package pgsqldb

import (
	"cache/internal/config"
	"cache/pkg/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PgSqlRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPgSqlDB(cfg config.PgSQL, logger logger.Logger) (*PgSqlRepository, error) {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		logger.Panicf("Database open error:%s", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Errorf("DB ping error:%s", err)
		return nil, err
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS articles(
		id TEXT PRIMARY KEY,		
		url TEXT NOT NULL,
		title TEXT NOT NULL);
	`)
	if err != nil {
		logger.Errorf("DB Prepare error:%s", err)
		return nil, fmt.Errorf("%s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		logger.Errorf("DB Exec error:%s", err)
		return nil, fmt.Errorf("%s", err)
	}

	return &PgSqlRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (m PgSqlRepository) Close() error {
	return m.db.Close()
}
