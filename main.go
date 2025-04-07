package main

import (
	"drfin/app/api"
	"drfin/app/database"
	"log"

	_ "github.com/lib/pq"
)

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
}
