package main

import (
	"Wildberries-L0/internal/message"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test-cluster", "client-123", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Отправка сообщения в канал
	id := uuid.NewString()
	msg, err := message.Generate("model.json", id)
	err = sc.Publish("my_channel", msg)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	log.Println("Сообщение успешно отправлено в канал")
}
