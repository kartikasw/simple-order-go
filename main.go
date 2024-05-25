package main

import (
	"fmt"
	"log"
	"simple-order-go/api"
	"simple-order-go/internal/handler"
	"simple-order-go/internal/repository"
	"simple-order-go/internal/service"
	config "simple-order-go/pkg/config"
	database "simple-order-go/pkg/db"
)

func main() {
	cfg := config.LoadConfig("app.yaml")

	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Init DB error: %v", err)
	}

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	handler := handler.NewOrderHandler(orderService)

	server := api.NewServer(*handler)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
