package db

import (
	"context"

	"delivery/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) CreateOrder(ctx context.Context, order models.Order) error {
	return db.db.Exec(`INSERT INTO delivery_orders(order_id) VALUES (?)`, order.OrdersID).Error
}

func (db *DB) SetStatusDelivered(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "delivered", id)
}

func (db *DB) updateOrderStatus(ctx context.Context, newStatus string, id int) (*models.Order, error) {
	order := &models.Order{}
	storeOrder := &models.Order{}

	err := db.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Model(order).Clauses(clause.Returning{}).Where("id = ?", id).Update("status", newStatus).Error
		if err != nil {
			return err
		}

		err = db.db.WithContext(ctx).Preload("Items").Where("id = ?", order.OrdersID).Find(&storeOrder).Error
		if err != nil {
			return err
		}

		return nil
	})

	return storeOrder, err
}
