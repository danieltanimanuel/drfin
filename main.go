package main

import (
	"drfin/app/api"
	"drfin/app/db"
	"io"
	"log"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	multi := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewJSONHandler(multi, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	db, err := db.InitiateDatabaseConnection()
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
