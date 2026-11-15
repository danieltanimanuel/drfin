package api

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"

	"drfin/app/accountgroups"
	"drfin/app/accounts"
	"drfin/app/categories"
	"drfin/app/dashboard"
	"drfin/app/gmail"
	"drfin/app/lastemail"
	"drfin/app/processedgmail"
	"drfin/app/shared/renderer"
	"drfin/app/transactions"
	"drfin/app/transactiontypes"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var tmplFS embed.FS

type APIServer struct {
	db   *sql.DB
	addr string
}

func NewAPIServer(db *sql.DB, addr string) *APIServer {
	return &APIServer{
		db:   db,
		addr: addr,
	}
}

func DRFINHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the X-Request-ID header if it exists
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// If not provided, generate a new UUID
			requestID = uuid.New().String()
		}

		// Log the request ID
		log.Printf("Request ID: %s", requestID)

		// Add the request ID to the request's context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "requestID", requestID)
		r = r.WithContext(ctx)

		// Continue with the next handler
		next.ServeHTTP(w, r)
	})
}

func (apis *APIServer) Run() error {
	log.Println("Starting Application DRFIN...")

	router := mux.NewRouter()
	router.Use(DRFINHandler)

	rend := renderer.NewRenderer("./app")

	log.Println("Assigning handlers...")
	// Initialize categories
	categoryRepository := categories.NewRepository(apis.db)
	categoryService := categories.NewService(categoryRepository)
	categoryHandler := categories.NewHandler(categoryService)

	// Initialize accountsGroups
	accountgroupsRepository := accountgroups.NewRepository(apis.db)
	accountgroupsService := accountgroups.NewService(accountgroupsRepository)
	accountgroupsHandler := accountgroups.NewHandler(accountgroupsService)

	// Initialize Transaction Types
	transactiontypesRepository := transactiontypes.NewRepository(apis.db)
	transactiontypesService := transactiontypes.NewService(transactiontypesRepository)
	transactiontypesHandler := transactiontypes.NewHandler(transactiontypesService)

	// Initialize Accounts
	accountsRepository := accounts.NewRepository(apis.db)
	accountsService := accounts.NewService(accountsRepository)
	accountsHandler := accounts.NewHandler(accountsService, rend)

	// Initialize Transactions
	transactionsRepository := transactions.NewRepository(apis.db)
	transactionsService := transactions.NewService(transactionsRepository, accountsRepository, transactiontypesRepository, categoryRepository)
	transactionsHandler := transactions.NewHandler(transactionsService)

	// Initialize last email
	lastEmailRepository := lastemail.NewRepository(apis.db)

	// Initialize processed gmail
	processedGmailRepository := processedgmail.NewRepository(apis.db)

	// Initialize Gmail
	gmailService := gmail.NewService(transactionsService, accountsService, lastEmailRepository, processedGmailRepository)
	gmailHandler := gmail.NewHandler(gmailService)

	dashboardHandler := dashboard.NewHandler(accountsService, transactionsService, categoryService, transactiontypesService, rend)

	// Register routes with specific methods (GET, POST, PUT, DELETE)
	categoryHandler.RegisterHandler(router)
	accountgroupsHandler.RegisterHandler(router)
	accountsHandler.RegisterHandler(router)
	transactiontypesHandler.RegisterHandler(router)
	transactionsHandler.RegisterHandler(router)
	gmailHandler.RegisterHandler(router)
	dashboardHandler.RegisterHandler(router)

	// Serve static files via router
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("app/static"))),
	)
	log.Println("DRFIN Application Started...")

	// Start the HTTP server with the specified address and router
	err := http.ListenAndServe(apis.addr, router)
	return err
}
