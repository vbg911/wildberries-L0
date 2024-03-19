package server

import (
	"Wildberries-L0/internal/common"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
)

func StartHTTPServer(ctx context.Context, c *cache.Cache, db *sql.DB) {
	// Создаем маршруты для HTTP-сервера
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, c, db)
	})

	// Запускаем HTTP-сервер
	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting HTTP server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %s\n", err)
		}
	}()

	// Ожидаем отмены контекста и завершаем сервер
	<-ctx.Done()
	log.Println("HTTP server shutdown")
	server.Shutdown(context.Background())
}

func handleRequest(w http.ResponseWriter, r *http.Request, c *cache.Cache, db *sql.DB) {
	uuid := r.URL.Query().Get("uuid")
	log.Println("new request received with uuid =", uuid)
	if uuid != "" {
		data, ok := c.Get(uuid)
		if ok {
			// Если данные найдены в кеше, преобразуем их в JSON и отправляем клиенту
			log.Println("cache founded for uuid =", uuid)
			order := data.(common.Order)
			jsonData, err := json.Marshal(order)
			if err != nil {
				http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
			return
		}

		var dataDB string
		// Выполнить запрос для получения строки по id
		row := db.QueryRow("SELECT data FROM orders WHERE id = $1", uuid)
		err := row.Scan(&dataDB)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("response from DB for uuid =", uuid)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(dataDB))

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
