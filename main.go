package main

import (
	"drfin/api"
	"drfin/database"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func DRFINHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RequestLogger2Middleware(r)
		next.ServeHTTP(w, r)
	}
}

func RequestLogger2Middleware(r *http.Request) {
	log.Printf("method: %s, path: %s ", r.Method, r.URL)
}

func main() {

	db, err := database.InitiateDatabaseConnection()
	if err != nil {
		log.Panic(err)
		return
	}
	apiServer := api.NewAPIServer(db, "localhost:8080")
	err = apiServer.Run()
	if err != nil {
		log.Panic(err)
		return
	}

	// transactionstype.InitHandlers(mux)
}
