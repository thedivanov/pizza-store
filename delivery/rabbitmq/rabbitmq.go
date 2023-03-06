package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"delivery/models"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	rabbitMQConn *amqp.Connection
	db           DB
}

type DB interface {
	CreateOrder(ctx context.Context, order models.Order) error
}

func NewRabbitMQ(amqpUser, amqpPassword, amqpAddr, amqpPort string, db DB) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", amqpUser, amqpPassword, amqpAddr, amqpPort))
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		rabbitMQConn: conn,
		db:           db,
	}, nil
}

func (r *RabbitMQ) SendDeliveredMessage(order models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.sendToQueue(orderJSON)
}

func (r *RabbitMQ) sendToQueue(val []byte) error {
	ch, err := r.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"delivered",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        val,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) ConsumeNewOrders(ctx context.Context) error {
	rChan, err := r.rabbitMQConn.Channel()
	if err != nil {
		return err
	}

	msgChan, err := rChan.Consume(
		"handover",
		"delivery",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgChan {
			order := models.Order{}
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.WithError(err).Error("Unmarshal RabbitMQ message error")
				continue
			}
			log.Println(order.OrdersID)
			err = r.db.CreateOrder(ctx, order)
			if err != nil {
				log.WithError(err).Error("Unmarshal RabbitMQ message error")
				continue
			}
		}
	}()

	return nil
}
