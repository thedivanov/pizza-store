package db

import (
	"context"
	"kitchen/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetOrders(ctx context.Context, limit, offset int, needTotal bool) ([]*models.KitchenOrder, int64, error) {
	var count int64
	orders := []*models.KitchenOrder{}

	q := db.db.Table(`kitchen_orders`)
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if needTotal {
		q = q.Count(&count)
	}
	err := q.WithContext(ctx).Preload("Order.Items").Find(&orders).Error
	return orders, count, err
}

func (db *DB) CreateOrder(ctx context.Context, order models.Order) error {
	return db.db.Exec(`INSERT INTO kitchen_orders(order_id) VALUES (?)`, order.ID).Error
}

func (db *DB) SetCookingOrder(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "cooking", id)
}

func (db *DB) SetCookedOrder(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "cooked", id)
}

func (db *DB) SetHandoverOrder(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "handover", id)
}

func (db *DB) updateOrderStatus(ctx context.Context, newStatus string, id int) (*models.Order, error) {
	order := &models.KitchenOrder{}
	storeOrder := &models.Order{}

	err := db.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Model(order).Clauses(clause.Returning{}).Where("id = ?", id).Update("status", newStatus).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Table("orders").Where("id = ?", order.OrderID).Update("status", newStatus).Error
		if err != nil {
			return err
		}

		err = db.db.WithContext(ctx).Preload("Items").Where("id = ?", order.OrderID).Find(&storeOrder).Error
		if err != nil {
			return err
		}

		return nil
	})

	return storeOrder, err
}
