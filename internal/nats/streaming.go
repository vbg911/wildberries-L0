package nats

import (
	"Wildberries-L0/internal/common"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"log"
)

func Connect(ctx context.Context, c *cache.Cache, db *sql.DB) {
	sc, err := stan.Connect("test-cluster", "subscriber-123", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS Streaming server: %v", err)
	}
	defer sc.Close()

	// Subscribe to a channel
	subscription, err := sc.Subscribe("my_channel", func(msg *stan.Msg) {
		handleMessage(msg, c, db)
	})
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	defer subscription.Unsubscribe()

	log.Println("Subscription to channel established successfully")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("NATS connection canceled")
			return
		}
	}
}

func handleMessage(msg *stan.Msg, c *cache.Cache, db *sql.DB) {
	log.Printf("message recived")
	var order common.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("Error while unmarshal json: %v\n", err)
		return
	}

	log.Println(order)

	validate := validator.New()
	err = validate.Struct(order)
	if err != nil {
		log.Printf("Validate json error: %v\n", err)
		return
	}
	c.Set(order.OrderUid, order, cache.NoExpiration)
	log.Println("cache added")

	_, err = db.Exec("INSERT INTO orders (id, data) VALUES ($1, $2)", order.OrderUid, msg.Data)
	if err != nil {
		log.Println(err)
	}

	log.Println("db row added")
}
