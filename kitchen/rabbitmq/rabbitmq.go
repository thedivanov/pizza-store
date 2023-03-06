package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"kitchen/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	RabbitMQConn *amqp.Connection
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
		RabbitMQConn: conn,
		db:           db,
	}, nil
}

func (r *RabbitMQ) CreateHandoverOrder(order models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.sendToQueue(orderJSON, "handover")
}

func (r *RabbitMQ) sendToQueue(val []byte, queueName string) error {
	ch, err := r.RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
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
	log.Debug(" [x] Sent %s\n", val)

	return nil
}

func (r *RabbitMQ) ConsumeNewOrders(ctx context.Context) error {
	rChan, err := r.RabbitMQConn.Channel()
	if err != nil {
		return err
	}

	msgChan, err := rChan.Consume(
		"order",
		"kitchen",
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
			err = r.db.CreateOrder(ctx, order)
			if err != nil {
				log.WithError(err).Error("Create order error")
				continue
			}
		}
	}()

	return nil
}
