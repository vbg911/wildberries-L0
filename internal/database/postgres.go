package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

const connStr = "user=postgres password=12345 dbname=wildberries host=localhost port=5432 sslmode=disable"

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to the database")
	return db
}

func GetAllOrders(db *sql.DB) (map[string]string, error) {
	orders := make(map[string]string)

	// Выполнить запрос для получения всех записей из таблицы orders
	rows, err := db.Query("SELECT id, data FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Проход по каждой строке результата запроса и добавление данных в карту
	for rows.Next() {
		var id string
		var data string
		if err := rows.Scan(&id, &data); err != nil {
			return nil, err
		}
		orders[id] = data
	}

	// Проверка наличия ошибок после прохода по всем строкам
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
