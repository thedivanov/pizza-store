package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	"shop/db"
	"shop/handlers"
	"shop/rabbitmq"
)

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.WithError(err).Error("ParseConfig err")
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
		log.WithError(err).Error("Postgres connection error")
		return
	}

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQAddress, cfg.RabbitMQPort, db)
	if err != nil {
		log.WithError(err).Error("Connect to RabbitMQ error")
		return
	}
	rabbitMQ.CreateQueue("order")
	rabbitMQ.CreateQueue("handover")
	rabbitMQ.CreateQueue("delivered")

	err = rabbitMQ.ConsumeNewOrders(context.TODO())
	if err != nil {
		log.WithError(err).Error("Consume delivered orders error")
		return
	}

	handler, err := handlers.NewHandler(e, rabbitMQ, db)
	if err != nil {
		log.WithError(err).Error("Init handler error")
		return
	}
	handler.AddURLs()

	srv := &http.Server{
		Addr:         ":9010",
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Info("Server starting at port 9010")

	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Error("Server error")
		return
	}
}
