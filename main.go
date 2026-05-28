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

// + базовый интерфейсы, которые будут определять общие контракты для всех типов баз данных и отправителей уведомлений.
type RepositoryWriter interface {
	CreateOrder(order Order) error
}

type Notifier interface {
	Send(customer string) error
}

// OrderService, зависит от интерфейсов RepositoryWriter и Notifier
type OrderService struct {
	repo     RepositoryWriter
	notifier Notifier
}

// EmailSender и SMSSender реализуют интерфейс Notifier, что позволяет легко переключаться между разными способами уведомлений в зависимости от условий.
type EmailSender struct{}

func (e *EmailSender) Send(customer string) error {
	fmt.Printf("E-mail: Уведомление отправлено клиенту %s\n", customer)
	return nil
}

type SMSSender struct{}

func (e *SMSSender) Send(customer string) error {
	fmt.Printf("SMS: Уведомление отправлено клиенту %s\n", customer)
	return nil
}

// SQLiteDatabase и PostgresqlDatabase реализуют интерфейс RepositoryWriter

// SQLiteDatabase для SQLite или любой другой СУБД, которая будет реализовывать интерфейс RepositoryWriter
type SQLiteDatabase struct{}

func (s *SQLiteDatabase) CreateOrder(order Order) error {
	db1, err := sql.Open("sqlite3", "sqlite_db/orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db1.Close()

	_, err = db1.Exec("INSERT INTO orders (customer, products, total, status) VALUES (?, ?, ?, ?)", order.Customer, order.Products, order.Total, order.Status)

	if err != nil {
		return err
	}

	return nil
}

// PostgresqlDatabase реализация для PostgreSQL
type PostgresqlDatabase struct {
	db *sql.DB
}

func (p *PostgresqlDatabase) CreateOrder(order Order) error {
	panic("Not yet implemented")
}

func CreateOrderService(useMssSQL bool, useSMS bool) *OrderService {
	var notifier Notifier
	if useSMS {
		notifier = &SMSSender{}
	} else {
		notifier = &EmailSender{}
	}

	if useMssSQL {
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
		return &OrderService{repo: &SQLiteDatabase{}, notifier: notifier}
	}

	return &OrderService{repo: &PostgresqlDatabase{}, notifier: notifier}
}

func (s *OrderService) CreateOrder(customerName string, products []string, total float64) (Order, error) {
	// Создание заказа в БД
	order := Order{
		Customer: customerName,
		Products: fmt.Sprintf("%v", products),
		Total:    total,
		Status:   "pending",
	}
	err := s.repo.CreateOrder(order)
	if err != nil {
		return order, err
	}

	return order, nil
}

func main() {

	// Создание сервиса с использованием SQLite и E-mail уведомлений
	serv := CreateOrderService(true, false)

	order, err := serv.CreateOrder("Иван", []string{"apple", "banana"}, 10.5)
	if err != nil {
		log.Fatal(err)
	}
	// Отправка уведомления
	serv.notifier.Send(order.Customer)

	order, err = serv.CreateOrder("Марья", []string{"prune", "carrot"}, 110.5)
	if err != nil {
		log.Fatal(err)
	}
	// Отправка уведомления
	serv.notifier.Send(order.Customer)

}
