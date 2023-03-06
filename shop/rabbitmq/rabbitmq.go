package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	RabbitMQConn *amqp.Connection
	db           DB
}

type DB interface {
	SetCompletedOrder(ctx context.Context, id int) (*models.Order, error)
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

func (r *RabbitMQ) CreateQueue(name string) error {
	rChan, err := r.RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	_, err = rChan.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) CreateOreder(order models.Order) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.sendToQueue(orderJSON)
}

func (r *RabbitMQ) sendToQueue(val []byte) error {
	ch, err := r.RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order",
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
	log.Debug("Send to queue:", val)

	return nil
}

func (r *RabbitMQ) ConsumeNewOrders(ctx context.Context) error {
	rChan, err := r.RabbitMQConn.Channel()
	if err != nil {
		return err
	}

	msgChan, err := rChan.Consume(
		"delivered",
		"shop",
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
			order := models.DeliveredOrder{}
			json.Unmarshal(d.Body, &order)
			_, err := r.db.SetCompletedOrder(context.TODO(), int(order.OrderID))
			if err != nil {
				log.WithError(err).Error("Set deliver status error")
				continue
			}
		}
	}()

	return nil
}
