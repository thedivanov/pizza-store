package models

type Order struct {
	ID       uint    `json:"-" gorm:"primary_key"`
	Items    []*Item `json:"items" validate:"required" gorm:"many2many:orders_items"`
	OrdersID uint    `json:"id" gorm:"column:order_id"`
}

func (o *Order) TableName() string {
	return "delivery_orders"
}

type Item struct {
	ID       uint    `json:"-" gorm:"primary_key"`
	Name     string  `json:"name" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
	Comment  string  `json:"comment"`
}

type Delivery struct {
	Address string `json:"delivery_address" gorm:"column:delivery_address"`
}

type ItemAdmin struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
	Comment  string  `json:"comment"`
}

func (ia *ItemAdmin) TableName() string {
	return "items"
}
