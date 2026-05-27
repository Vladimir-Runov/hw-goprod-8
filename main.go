// go get github.com/mattn/go-sqlite3
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Order struct {
	ID       int
	Customer string
	Products string
	Total    float64
	Status   string
}

type OrderSystem struct {
	db *sql.DB
}

func NewOrderSystem(db *sql.DB) *OrderSystem {
	return &OrderSystem{db: db}
}

func (s *OrderSystem) CreateOrder(customer string, products []string, total float64) error {
	// Создание заказа в БД
	_, err := s.db.Exec(
		"INSERT INTO orders (customer, products, total, status) VALUES (?, ?, ?, ?)",
		customer, fmt.Sprintf("%v", products), total, "pending",
	)
	if err != nil {
		return err
	}

	// Отправка уведомления
	s.sendEmailNotification(customer)

	return nil
}

func (s *OrderSystem) sendEmailNotification(customer string) {
	fmt.Printf("Уведомление отправлено клиенту %s\n", customer)
}

func main() {
	db, err := sql.Open("sqlite3", "sqlite_db/orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer TEXT NOT NULL,
        products TEXT NOT NULL,
        total REAL NOT NULL,
        status TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	system := NewOrderSystem(db)

	err = system.CreateOrder("Иван", []string{"apple", "banana"}, 10.5)
	if err != nil {
		log.Fatal(err)
	}
}
