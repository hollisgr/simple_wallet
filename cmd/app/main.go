package main

import (
	"cmd/app/main.go/internal/app"
	"cmd/app/main.go/internal/config"
	"cmd/app/main.go/internal/db"
	"cmd/app/main.go/internal/service"
)

func main() {

	cfg := config.GetConfig()

	pool := app.ConnectToDB(cfg)
	defer pool.Close()

	storage := db.New(pool)
	ws := service.New(storage)

	srv := app.SetupServer(cfg, ws)

	app.StartServer(srv)

	app.HandleQuit(srv)
}
