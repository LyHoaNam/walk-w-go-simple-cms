package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// DB contains database connection and goqu dialect
type DB struct {
	SQL     *sql.DB
	Dialect goqu.DialectWrapper
}

// Config contains database configuration
type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Connect creates a connection to MySQL database
func Connect(cfg Config) (*DB, error) {
	// Create DSN (Data Source Name) for MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	// Open database connection
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test kết nối
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create goqu dialect wrapper
	dialect := goqu.Dialect("mysql")

	return &DB{
		SQL:     sqlDB,
		Dialect: dialect,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.SQL.Close()
}
