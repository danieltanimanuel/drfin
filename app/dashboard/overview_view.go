package dashboard

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"time"

	"drfin/app/accounts"
	"drfin/app/transactions"
)

// Account group IDs. They match the seeded account_groups rows and the
// group list already hard-coded in the create-account form and in
// accounts.Repository.CreateReceivableAccount (Loan = 7).
const (
	groupCash    = 1
	groupDebit   = 2
	groupCredit  = 3
	groupEMoney  = 4
	groupEWallet = 5
	groupDebt    = 6
	groupLoan    = 7
	groupSaving  = 8
)

const (
	kindAvailable = iota
	kindOwedToYou
	kindYouOwe
	kindOther
)

// groupDisplayOrder is the order of group sections on the dashboard.
var groupDisplayOrder = []int{
	groupCash, groupDebit, groupEMoney, groupEWallet, groupSaving,
	groupLoan, groupCredit, groupDebt,
}

// ----------------------
// View models
// ----------------------

type OverviewView struct {
	MonthLabel string
	Income     *big.Float
	Expense    *big.Float
	Net        *big.Float
	Blocks     []BalanceBlock
}

type BalanceBlock struct {
	Title  string
	Hint   string
	Total  *big.Float
	Groups []GroupSubtotalView
}

type GroupSubtotalView struct {
	Name  string
	Total *big.Float
}

type AccountSectionView struct {
	GroupID int
	Name    string
	Class   string
	Total   *big.Float
	Cards   []AccountCardView
}

type AccountCardView struct {
	ID              int
	Name            string
	GroupName       string
	Class           string
	Balance         *big.Float
	ShowCredit      bool
	CreditLimit     *big.Float
	AvailableCredit *big.Float
}

type CashflowView struct {
	Months []CashflowMonthView
}

type CashflowMonthView struct {
	Label      string // short, e.g. "Nov"
	FullLabel  string // e.g. "Nov 2026"
	Income     *big.Float
	Expense    *big.Float
	IncomePct  float64 // bar height, 0-100, relative to the largest value shown
	ExpensePct float64
}

type CategoryBreakdownView struct {
	CurrentMonthStart string
	CurrentMonthEnd   string
	Total             *big.Float
	Items             []CategoryRowView
}

type CategoryRowView struct {
	Name     string
	Amount   *big.Float
	PctLabel string  // e.g. "12.3"
	Bar      float64 // 0-100
}

type RecentTxnView struct {
	Txn          transactions.TransactionDTO
	Class        string // income, expense, transfer, adjustment, other
	Sign         string // "+", "-" or ""
	Amount       *big.Float
	ShowCategory bool
	ShowDest     bool
}

// monthTotal is the exact income/expense of one month.
type monthTotal struct {
	Start   time.Time
	Income  *big.Rat
	Expense *big.Rat
}

// ----------------------
// Money helpers
// ----------------------

// moneyRat converts a *big.Float holding a 2-decimal amount into an exact
// *big.Rat, so sums and differences have no binary rounding error.
func moneyRat(f *big.Float) *big.Rat {
	r := new(big.Rat)
	if f == nil {
		return r
	}
	if _, ok := r.SetString(f.Text('f', 2)); !ok {
		return new(big.Rat)
	}
	return r
}

// ratFloat converts an exact amount back to *big.Float for display only.
func ratFloat(r *big.Rat) *big.Float {
	return new(big.Float).SetRat(r)
}

func sumAggregates(aggs []transactions.CategoryAggregate) *big.Rat {
	total := new(big.Rat)
	for _, a := range aggs {
		total.Add(total, moneyRat(a.TotalAmount))
	}
	return total
}

// monthRange returns the first and last day of the month containing t.
func monthRange(t time.Time) (time.Time, time.Time) {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return first, first.AddDate(0, 1, -1)
}

// ----------------------
// Account grouping
// ----------------------

type accountGroup struct {
	ID       int
	Name     string
	Total    *big.Rat
	Accounts []accounts.AccountDTO
}

func groupRank(id int) int {
	for i, g := range groupDisplayOrder {
		if g == id {
			return i
		}
	}
	return len(groupDisplayOrder) + id
}

func groupKind(id int) int {
	switch id {
	case groupCash, groupDebit, groupEMoney, groupEWallet, groupSaving:
		return kindAvailable
	case groupLoan:
		return kindOwedToYou
	case groupCredit, groupDebt:
		return kindYouOwe
	default:
		return kindOther
	}
}

func groupClass(id int) string {
	switch id {
	case groupCash:
		return "cash"
	case groupDebit:
		return "debit"
	case groupCredit:
		return "credit"
	case groupEMoney:
		return "emoney"
	case groupEWallet:
		return "ewallet"
	case groupDebt:
		return "debt"
	case groupLoan:
		return "loan"
	case groupSaving:
		return "saving"
	default:
		return "other"
	}
}

func groupAccounts(accs []accounts.AccountDTO) []*accountGroup {
	byID := make(map[int]*accountGroup)
	var list []*accountGroup
	for _, a := range accs {
		g, ok := byID[a.AccountGroup.ID]
		if !ok {
			name := a.AccountGroup.Name
			if name == "" {
				name = "Other"
			}
			g = &accountGroup{ID: a.AccountGroup.ID, Name: name, Total: new(big.Rat)}
			byID[g.ID] = g
			list = append(list, g)
		}
		g.Total.Add(g.Total, moneyRat(a.Balance))
		g.Accounts = append(g.Accounts, a)
	}
	sort.SliceStable(list, func(i, j int) bool {
		return groupRank(list[i].ID) < groupRank(list[j].ID)
	})
	return list
}

// ----------------------
// Builders
// ----------------------

// buildOverview never adds balances of different kinds together: available
// funds, money owed to you (Loan) and money you owe (Credit, Debt) are
// separate blocks. "You owe" is shown as a positive amount owed.
func buildOverview(now time.Time, income, expense *big.Rat, accs []accounts.AccountDTO) OverviewView {
	type blockAcc struct {
		block BalanceBlock
		total *big.Rat
	}
	blocks := map[int]*blockAcc{
		kindAvailable: {block: BalanceBlock{Title: "Available funds", Hint: "Cash, debit, e-money, e-wallet and savings"}, total: new(big.Rat)},
		kindOwedToYou: {block: BalanceBlock{Title: "Owed to you", Hint: "Loan accounts"}, total: new(big.Rat)},
		kindYouOwe:    {block: BalanceBlock{Title: "You owe", Hint: "Credit and debt accounts, shown as amount owed"}, total: new(big.Rat)},
		kindOther:     {block: BalanceBlock{Title: "Other accounts", Hint: "Not included in the blocks above"}, total: new(big.Rat)},
	}

	for _, g := range groupAccounts(accs) {
		kind := groupKind(g.ID)
		b := blocks[kind]
		v := new(big.Rat).Set(g.Total)
		if kind == kindYouOwe {
			v.Neg(v)
		}
		b.total.Add(b.total, v)
		b.block.Groups = append(b.block.Groups, GroupSubtotalView{Name: g.Name, Total: ratFloat(v)})
	}

	view := OverviewView{
		MonthLabel: now.Format("January 2006"),
		Income:     ratFloat(income),
		Expense:    ratFloat(expense),
		Net:        ratFloat(new(big.Rat).Sub(income, expense)),
	}
	for _, kind := range []int{kindAvailable, kindOwedToYou, kindYouOwe, kindOther} {
		b := blocks[kind]
		if len(b.block.Groups) == 0 && kind != kindAvailable {
			continue
		}
		b.block.Total = ratFloat(b.total)
		view.Blocks = append(view.Blocks, b.block)
	}
	return view
}

func buildAccountSections(accs []accounts.AccountDTO) []AccountSectionView {
	groups := groupAccounts(accs)
	sections := make([]AccountSectionView, 0, len(groups))
	for _, g := range groups {
		sec := AccountSectionView{
			GroupID: g.ID,
			Name:    g.Name,
			Class:   groupClass(g.ID),
			Total:   ratFloat(g.Total),
		}
		for _, a := range g.Accounts {
			card := AccountCardView{
				ID:        a.ID,
				Name:      a.Name,
				GroupName: g.Name,
				Class:     sec.Class,
				Balance:   a.Balance,
			}
			// Credit balances are stored negative (amount owed), so the
			// available credit is limit + balance.
			if g.ID == groupCredit && a.CreditLimit != nil {
				limit := moneyRat(a.CreditLimit)
				if limit.Sign() > 0 {
					card.ShowCredit = true
					card.CreditLimit = a.CreditLimit
					card.AvailableCredit = ratFloat(new(big.Rat).Add(limit, moneyRat(a.Balance)))
				}
			}
			sec.Cards = append(sec.Cards, card)
		}
		sections = append(sections, sec)
	}
	return sections
}

func barPct(v, max *big.Rat) float64 {
	if max.Sign() == 0 || v.Sign() <= 0 {
		return 0
	}
	f, _ := new(big.Rat).Quo(v, max).Float64()
	p := math.Round(f*1000) / 10
	if p < 2 {
		p = 2 // keep a non-zero bar visible
	}
	return p
}

func buildCashflow(months []monthTotal) CashflowView {
	max := new(big.Rat)
	for _, m := range months {
		if m.Income.Cmp(max) > 0 {
			max.Set(m.Income)
		}
		if m.Expense.Cmp(max) > 0 {
			max.Set(m.Expense)
		}
	}
	view := CashflowView{Months: make([]CashflowMonthView, 0, len(months))}
	for _, m := range months {
		view.Months = append(view.Months, CashflowMonthView{
			Label:      m.Start.Format("Jan"),
			FullLabel:  m.Start.Format("Jan 2006"),
			Income:     ratFloat(m.Income),
			Expense:    ratFloat(m.Expense),
			IncomePct:  barPct(m.Income, max),
			ExpensePct: barPct(m.Expense, max),
		})
	}
	return view
}

func buildCategoryBreakdown(from, to time.Time, aggs []transactions.CategoryAggregate) CategoryBreakdownView {
	amounts := make([]*big.Rat, len(aggs))
	total := new(big.Rat)
	for i, a := range aggs {
		amounts[i] = moneyRat(a.TotalAmount)
		total.Add(total, amounts[i])
	}

	rows := make([]CategoryRowView, 0, len(aggs))
	for i, a := range aggs {
		pct := 0.0
		if total.Sign() > 0 {
			f, _ := new(big.Rat).Quo(amounts[i], total).Float64()
			pct = f * 100
		}
		rows = append(rows, CategoryRowView{
			Name:     a.CategoryName,
			Amount:   a.TotalAmount,
			PctLabel: fmt.Sprintf("%.1f", pct),
			Bar:      math.Round(pct*10) / 10,
		})
	}

	return CategoryBreakdownView{
		CurrentMonthStart: from.Format("2006-01-02"),
		CurrentMonthEnd:   to.Format("2006-01-02"),
		Total:             ratFloat(total),
		Items:             rows,
	}
}

// buildRecent flattens the day groups (already newest day first, newest id
// first within a day) and returns the first `limit` non-cancelled rows.
func buildRecent(days []transactions.DailyTransactionsDTO, limit int) []RecentTxnView {
	rows := make([]RecentTxnView, 0, limit)
	for _, day := range days {
		for _, tx := range day.Transactions {
			if tx.Status == transactions.StatusCancelled {
				continue
			}
			rows = append(rows, recentRow(tx))
			if len(rows) == limit {
				return rows
			}
		}
	}
	return rows
}

func recentRow(tx transactions.TransactionDTO) RecentTxnView {
	row := RecentTxnView{
		Txn:          tx,
		Class:        "other",
		Amount:       tx.Amount,
		ShowCategory: true,
	}
	switch tx.Type.Name {
	case "Income":
		row.Class = "income"
		row.Sign = "+"
	case "Expense":
		row.Class = "expense"
		row.Sign = "-"
	case "Transfer":
		row.Class = "transfer"
		row.ShowCategory = false
		// Only transfers have a meaningful destination; other types may
		// carry a stored destination that the balance logic ignores.
		row.ShowDest = tx.AccountDest != nil
	case "Adjustment":
		row.Class = "adjustment"
		row.ShowCategory = false
		// The stored amount is a signed delta; show it as sign + magnitude.
		if tx.Amount != nil {
			switch tx.Amount.Sign() {
			case -1:
				row.Sign = "-"
				row.Amount = new(big.Float).Abs(tx.Amount)
			case 1:
				row.Sign = "+"
			}
		}
	}
	return row
}
