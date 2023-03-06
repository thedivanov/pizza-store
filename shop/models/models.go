package models

type Order struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Items    []*Item   `json:"items" validate:"required" gorm:"many2many:orders_items"`
	Delivery *Delivery `json:"delivery" validate:"required" gorm:"embedded"`
	Status   string    `json:"status" gorm:"default:new"`
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

type GetOrdersResponse struct {
	Orders []*Order `json:"orders"`
}

type GetOrdersAdminResponse struct {
	Orders []*Order `json:"orders"`
	Total  int64    `json:"total"`
	Offset int64    `json:"offset"`
}

type DeliveredOrder struct {
	ID       uint      `gorm:"primary_key"`
	Items    []*Item   `json:"items" validate:"required" gorm:"many2many:orders_items"`
	Delivery *Delivery `json:"delivery" validate:"required" gorm:"embedded"`
	OrderID  uint      `json:"id" gorm:"foreignKey"`
}
