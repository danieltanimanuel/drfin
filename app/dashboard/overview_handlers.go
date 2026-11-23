package dashboard

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"time"
)

const (
	recentTransactionLimit = 8
	cashflowMonths         = 6
)

// monthTotals returns the exact Income and Expense totals for [from, to].
// It reuses the same aggregate queries as the category breakdown, so the
// expense total always equals the sum of the breakdown rows for that range.
// Transfers, adjustments, lend/borrow and cancelled rows are not counted.
func (h *Handler) monthTotals(ctx context.Context, from, to time.Time, ownerID *int) (*big.Rat, *big.Rat, error) {
	incomeAggs, err := h.transactionService.GetIncomeCategoryAggregates(ctx, from, to, ownerID)
	if err != nil {
		return nil, nil, err
	}
	expenseAggs, err := h.transactionService.GetExpenseCategoryAggregates(ctx, from, to, ownerID)
	if err != nil {
		return nil, nil, err
	}
	return sumAggregates(incomeAggs), sumAggregates(expenseAggs), nil
}

// OverviewGrid renders the month income/expense/net cards and the balance blocks.
func (h *Handler) OverviewGrid(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ownerID, ok := ownerFromRequest(w, r); if !ok { return }
	accountsInq, err := h.accountService.GetAllAccounts(&ownerID)
	if err != nil {
		http.Error(w, "failed to load accounts", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	from, to := monthRange(now)
	income, expense, err := h.monthTotals(ctx, from, to, &ownerID)
	if err != nil {
		log.Printf("Error loading month totals: %v", err)
		http.Error(w, "failed to load month totals", http.StatusInternalServerError)
		return
	}

	err = h.renderer.Render(w, "dashboard/templates/", "overview_grid.html",
		buildOverview(now, income, expense, accountsInq))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CashflowGrid renders income vs expenses for the last cashflowMonths months.
func (h *Handler) CashflowGrid(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ownerID, ok := ownerFromRequest(w, r); if !ok { return }
	currentMonth, _ := monthRange(time.Now())
	months := make([]monthTotal, 0, cashflowMonths)
	for i := cashflowMonths - 1; i >= 0; i-- {
		start := currentMonth.AddDate(0, -i, 0)
		_, end := monthRange(start)
		income, expense, err := h.monthTotals(ctx, start, end, &ownerID)
		if err != nil {
			log.Printf("Error loading cash flow for %s: %v", start.Format("2006-01"), err)
			http.Error(w, "failed to load cash flow", http.StatusInternalServerError)
			return
		}
		months = append(months, monthTotal{Start: start, Income: income, Expense: expense})
	}

	err := h.renderer.Render(w, "dashboard/templates/", "cashflow_grid.html", buildCashflow(months))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RecentGrid renders the latest recentTransactionLimit non-cancelled transactions.
func (h *Handler) RecentGrid(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ownerID, ok := ownerFromRequest(w, r); if !ok { return }
	days, err := h.transactionService.GetAllTransactions(ctx, &ownerID)
	if err != nil {
		http.Error(w, "failed to load transactions", http.StatusInternalServerError)
		return
	}

	err = h.renderer.Render(w, "dashboard/templates/", "recent_grid.html",
		buildRecent(days, recentTransactionLimit))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AccountCreateForm renders the New Account form for the modal.
func (h *Handler) AccountCreateForm(w http.ResponseWriter, r *http.Request) {
	err := h.renderer.Render(w, "dashboard/templates/", "account_create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
