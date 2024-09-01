package api

import (
	"database/sql"
	"drfin/transactionstype"
	"log"
	"net/http"
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
	mux := http.NewServeMux()
	log.Println("Assign handlers...")
	transactionTypeHandler := transactionstype.NewHandler(apis.db)
	transactionTypeHandler.RegisterHandler(mux)
	log.Println("DRFIN Application Started...")
	err := http.ListenAndServe(apis.addr, mux)
	return err
}
