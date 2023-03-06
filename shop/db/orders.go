package db

import (
	"context"
	"shop/models"

	"gorm.io/gorm/clause"
)

func (db *DB) GetOrders(ctx context.Context, offset, limit int, needTotal bool) ([]*models.Order, int64, error) {
	var count int64
	orders := []*models.Order{}

	q := db.db.Table("orders")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if needTotal {
		q = q.Count(&count)
	}

	err := q.WithContext(ctx).Preload("Items").Find(&orders).Error
	return orders, count, err
}

func (db *DB) CreateOreder(ctx context.Context, order *models.Order) error {
	return db.db.WithContext(ctx).Create(order).Error
}

func (db *DB) SetConfirmOrder(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "confirmed", id)
}

func (db *DB) SetCancelOrder(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "canceled", id)
}

func (db *DB) SetCompletedOrder(ctx context.Context, id int) (*models.Order, error) {
	return db.updateOrderStatus(ctx, "completed", id)
}

func (db *DB) updateOrderStatus(ctx context.Context, newStatus string, id int) (*models.Order, error) {
	order := &models.Order{}

	err := db.db.WithContext(ctx).Model(order).Clauses(clause.Returning{}).Table("orders").Where("id = ?", id).Update("status", newStatus).Error
	if err != nil {
		return nil, err
	}

	err = db.db.WithContext(ctx).Preload("Items").Find(&order).Where("id = ?").Error

	return order, err
}
