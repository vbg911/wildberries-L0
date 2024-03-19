package main

import (
	"Wildberries-L0/internal/common"
	"Wildberries-L0/internal/database"
	"Wildberries-L0/internal/nats"
	"Wildberries-L0/internal/server"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

func main() {
	// Создаем контекст для управления выполнением
	ctx := context.Background()
	db := database.ConnectDB()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Println(err)
	}

	c := cache.New(cache.NoExpiration, 10*time.Minute)
	orders, err := database.GetAllOrders(db)
	if err != nil {
		log.Println(err)
	}

	// Заполнить кеш данными о заказах
	for id, data := range orders {
		var order common.Order
		err := json.Unmarshal([]byte(data), &order)
		if err != nil {
			log.Printf("Error while unmarshal json: %v\n", err)
			return
		}
		c.Set(id, order, cache.NoExpiration)
	}

	// Запускаем goroutine для подключения к NATS
	go func(c *cache.Cache, db *sql.DB) {
		nats.Connect(ctx, c, db)
	}(c, db)

	// Запускаем goroutine для запуска HTTP-сервера
	go func(c *cache.Cache, db *sql.DB) {
		server.StartHTTPServer(ctx, c, db)
	}(c, db)

	// Бесконечный цикл, чтобы основная горутина не завершилась
	select {}
}
