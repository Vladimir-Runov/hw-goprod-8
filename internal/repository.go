package internal

import (
  "database/sql"
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

// RepositoryWriter определяет контракт для работы с базой данных.
type RepositoryWriter interface {
  CreateOrder(order Order) error
}

// SQLiteDatabase реализует интерфейс RepositoryWriter для SQLite.
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

// PostgresqlDatabase реализация для PostgreSQL.
type PostgresqlDatabase struct {
  db *sql.DB
}

func (p *PostgresqlDatabase) CreateOrder(order Order) error {
  panic("Not yet implemented")
}

// Создает и инициализирует базу данных.
func InitDatabase(useMssSQL bool) {
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
}