package api

import (
	"database/sql"
	"log"
	"net/http"

	"drfin/accountgroups"
	"drfin/categories"
	"drfin/transactiontypes"

	"github.com/gorilla/mux"
)

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

func (apis *APIServer) Run() error {
	log.Println("Starting Application DRFIN...")

	// Initialize the router using gorilla/mux
	router := mux.NewRouter()

	log.Println("Assigning handlers...")
	// Initialize the handler for categories, accountgroups, types
	categoryHandler := categories.NewHandler(categories.NewService(categories.NewRepository(apis.db)))
	accountgroupsHandler := accountgroups.NewHandler(accountgroups.NewService(accountgroups.NewRepository(apis.db)))
	transactiontypesHandler := transactiontypes.NewHandler(transactiontypes.NewService(transactiontypes.NewRepository(apis.db)))

	// Register routes with specific methods (GET, POST, PUT, DELETE)
	categoryHandler.RegisterHandler(router)
	accountgroupsHandler.RegisterHandler(router)
	transactiontypesHandler.RegisterHandler(router)

	log.Println("DRFIN Application Started...")

	// Start the HTTP server with the specified address and router
	err := http.ListenAndServe(apis.addr, router)
	return err
}
