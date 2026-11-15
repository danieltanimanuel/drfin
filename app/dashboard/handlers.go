package dashboard

import (
	"context"
	"drfin/app/accounts"
	"drfin/app/categories"
	"drfin/app/transactions"
	"drfin/app/transactiontypes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	accountService     *accounts.Service
	transactionService *transactions.Service
	categoryService    *categories.Service
	transactionTypeService *transactiontypes.Service
	renderer           Renderer
}

type Renderer interface {
	Render(w http.ResponseWriter, templatePath string, templateName string, data interface{}) error
}

func NewHandler(
    accountService *accounts.Service,
    transactionService *transactions.Service,
    categoryService *categories.Service,
    transactionTypeService *transactiontypes.Service,
    r Renderer,
) *Handler {
    return &Handler{
        accountService:         accountService,
        transactionService:      transactionService,
        categoryService:         categoryService,
        transactionTypeService: transactionTypeService,
        renderer:                r,
    }
}

func (h *Handler) RegisterHandler(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})
	router.HandleFunc("/dashboard", h.Dashboard)
	router.HandleFunc("/dashboard/accounts", h.AccountsGrid)
	router.HandleFunc("/dashboard/pie", h.PieGrid)
	router.HandleFunc("/dashboard/categories", h.GetExpenseCategoryAggregates)
	router.HandleFunc("/dashboard/transactions", h.TransactionsGrid)
	router.HandleFunc(
		"/dashboard/transactions/new",
		h.CreateTransactionForm,
	).Methods("GET")
	router.HandleFunc("/dashboard/transactions/{id}/edit", h.EditTransactionForm)
	router.HandleFunc("/dashboard/transactions/{id}", h.UpdateTransactionCategory).Methods("PUT")
	router.HandleFunc(
		"/dashboard/transactions/{id}",
		h.DeleteTransaction,
	).Methods("DELETE")
}

// ----------------------
// Dashboard shell
// ----------------------
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Just render the shell (AJAX will load grids)
	err := h.renderer.Render(w, "dashboard/templates/", "dashboard.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ----------------------
// Accounts Grid
// ----------------------
func (h *Handler) AccountsGrid(w http.ResponseWriter, r *http.Request) {
	accountsInq, err := h.accountService.GetAllAccounts()
	if err != nil {
		http.Error(w, "failed to load accounts", http.StatusInternalServerError)
		return
	}

	err = h.renderer.Render(w, "dashboard/templates/", "accounts_grid.html", accountsInq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ----------------------
// Pie Grid (JSON only)
// ----------------------
func (h *Handler) PieGrid(w http.ResponseWriter, r *http.Request) {
	accountsInq, err := h.accountService.GetAllAccounts()
	if err != nil {
		http.Error(w, "failed to load pie data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountsInq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ----------------------
// Transactions Grid
// ----------------------
func (h *Handler) TransactionsGrid(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	transactionInq, err := h.transactionService.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, "failed to load transactions", http.StatusInternalServerError)
		return
	}

	err = h.renderer.Render(w, "dashboard/templates/", "transactions_grid.html", transactionInq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetTransactionAccountGroups handles GET requests to retrieve all transaction accountgroups
func (h *Handler) GetAllAcounts(w http.ResponseWriter, r *http.Request) {
	accountsInq, err := h.accountService.GetAllAccounts()
	if err != nil {
		http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(accountsInq)
	if err != nil {
		http.Error(w, "Failed to encode accounts", http.StatusInternalServerError)
		return
	}

	data := accounts.AccountPage{
		Accounts:     accountsInq,
		AccountsJSON: template.JS(jsonData), // mark safe for template
	}

	// Decide template based on HTMX header
	var templateName string
	if r.Header.Get("HX-Request") == "true" {
		templateName = "list_partial.html"
	} else {
		templateName = "list_full.html"
	}
	templatePath := "dashboard/templates/"

	log.Println(accountsInq)
	err = h.renderer.Render(w, templatePath, templateName, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Something error :(", http.StatusInternalServerError)
	}
}



func (h *Handler) EditTransactionForm(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	txn, err := h.transactionService.GetTransactionByIDDTO(ctx, id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		http.Error(w, "Failed to load accounts", http.StatusInternalServerError)
		return
	}

	transactionTypes, err := h.transactionTypeService.GetAllTransactionTypes()
	if err != nil {
		http.Error(w, "Failed to load transaction types", http.StatusInternalServerError)
		return
	}

	data := EditTransactionView{
		Transaction: txn,
		Categories:  categories,
		Accounts:    accounts,
    	TransactionTypes: transactionTypes,
	}

	err = h.renderer.Render(
		w,
		"dashboard/templates/",
		"transaction_edit_item.html",
		data,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) UpdateTransactionCategory(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Type        string `json:"type"`
		Category    string `json:"category"`
		Account     string `json:"account"`
		AccountDest string `json:"account_dest"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	category := payload.Category

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	dto := transactions.UpdateTransactionDTO{
		Category: &category,
		Account:  &payload.Account,
		AccountDest: &payload.AccountDest,
		Type: &payload.Type,
	}

	_, err = h.transactionService.UpdateTransaction(ctx, id, dto)
	if err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	// Re-render full grouped transactions grid
	h.TransactionsGrid(w, r)
}

func (h *Handler) GetExpenseCategoryAggregates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Example: This Month
	fromStr := r.URL.Query().Get("from")
    toStr   := r.URL.Query().Get("to")

	now := time.Now()

	// Default: current month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay  := firstDay.AddDate(0, 1, -1)

    var from time.Time
    var to   time.Time
    var err  error

    if fromStr == "" {
		from = firstDay
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "Invalid from date", http.StatusBadRequest)
			return
		}
	}

	if toStr == "" {
		to = lastDay
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "Invalid to date", http.StatusBadRequest)
			return
		}
	}

	expenseTotals, err := h.transactionService.GetExpenseCategoryAggregates(
		ctx,
		from,
		to,
	)
	if err != nil {
		http.Error(w, "Failed to load category summary", 500)
		return
	}

	data := struct {
		CurrentMonthStart string
		CurrentMonthEnd   string
		Items             interface{}
	}{
		CurrentMonthStart: from.Format("2006-01-02"),
		CurrentMonthEnd:   to.Format("2006-01-02"),
		Items:             expenseTotals,
	}

	err = h.renderer.Render(
		w,
		"dashboard/templates/",
		"categories_grid.html",
		data,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h *Handler) CreateTransactionForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		http.Error(w, "failed categories", 500)
		return
	}

	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		http.Error(w, "failed accounts", 500)
		return
	}

	data := CreateTransactionView{
		Categories: categories,
		Accounts:   accounts,
	}

	err = h.renderer.Render(
		w,
		"dashboard/templates/",
		"transaction_create.html",
		data,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.transactionService.DeleteTransaction(ctx, id)
	if err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}