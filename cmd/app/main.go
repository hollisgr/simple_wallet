package main

import (
	"cmd/app/main.go/internal/config"
	"cmd/app/main.go/internal/handler"
	"cmd/app/main.go/pkg/postgres"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	router := gin.Default()

	pgxPool, err := postgres.NewPool(context.Background(), 5, cfg.Postgresql.DSN)
	if err != nil {
		log.Fatalln("cant connect to db")
	}
	log.Println("db connection OK")

	err = pgxPool.Ping(context.Background())
	if err != nil {
		log.Fatalln("cant ping to db, err:", err)
	}
	log.Println("db ping OK")

	h := handler.New(router)
	h.Register()

	router.Run(cfg.Listen.Addr)
	defer pgxPool.Close()
}
