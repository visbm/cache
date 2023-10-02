package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"cache/internal/config"
	"cache/internal/handlers"
	"cache/internal/httpserver"
	"cache/internal/repository"
	"cache/internal/repository/pgsqldb"
	"cache/internal/repository/redisdb"
	"cache/internal/service"
	"cache/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := logger.GetLogger()
	logger.Info("Starting application")

	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("config loaded ", config.Env)

	cache, err := redisdb.NewRedisDatabase(config.Redis, logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("cache created")

	db, err := pgsqldb.NewPgSqlDB(config.PgSQL, logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("db created")

	repository := repository.NewMainRepository(cache, db, logger)
	logger.Info("repository created ")

	service := service.NewArticleService(repository, logger)
	articleHandler := handlers.NewArticleHandler(service, logger)
	handler := handlers.NewHandler(articleHandler)

	server := httpserver.NewServer(config.HttpServer, handler.InitRoutes(), logger)

	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		err = server.Shutdown(context.Background())
		if err != nil {
			logger.Errorf("Error occurred on server shutting down: %s", err.Error())
		}

		err = db.Close()
		if err != nil {
			logger.Errorf("Error occurred on db connection close: %s", err.Error())
		}

		err = cache.Close()
		if err != nil {
			logger.Errorf("Error occurred on db connection close: %s", err.Error())
		}

		logger.Info("shutting down")
		os.Exit(0)
	}()

	if err := server.Start(); err != http.ErrServerClosed {
		logger.Panicf("Error while starting server:%s", err)
	}
	<-idleConnsClosed

}
