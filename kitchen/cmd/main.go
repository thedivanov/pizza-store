package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	"kitchen/db"
	"kitchen/handlers"
	"kitchen/rabbitmq"
)

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.WithError(err).Error("Parse config error")
		return
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
	}))

	db, err := db.NewDB(cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresAddr, cfg.PostgresPort, cfg.PostgresDBName)
	if err != nil {
		log.WithError(err).Error("Connect to database error")
		return
	}

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQAddress, cfg.RabbitMQPort, db)
	if err != nil {
		log.WithError(err).Error("Connect to rabbit error")
		return
	}

	err = rabbitMQ.ConsumeNewOrders(context.TODO())
	if err != nil {
		log.WithError(err).Error("Consume new orders error")
		return
	}

	handler, err := handlers.NewHandler(e, rabbitMQ, db)
	if err != nil {
		log.Fatal(err)
	}
	handler.AddURLs()

	log.Info("Server starting at port 9020")

	srv := &http.Server{
		Addr:         ":9020",
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
