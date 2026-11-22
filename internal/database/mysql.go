package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/bintangmada/bookstore-api-go/internal/config"
	"github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg config.Config) (*sql.DB, error) {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.DBUser
	mysqlCfg.Passwd = cfg.DBPassword
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = net.JoinHostPort(cfg.DBHost, cfg.DBPort)
	mysqlCfg.DBName = cfg.DBName
	mysqlCfg.ParseTime = true
	mysqlCfg.Params = map[string]string{
		"charset": "utf8mb4",
	}

	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	// Pengaturan connection pool (mirip HikariCP di Java)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql database %q: %w", cfg.DBName, err)
	}

	return db, nil
}
