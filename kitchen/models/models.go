package models

type KitchenOrder struct {
	ID      uint `gorm:"primary_key"`
	Status  string
	Order   *Order `gorm:"foreignKey:OrderID;references:ID"`
	OrderID uint
}

type Order struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Items []*Item `json:"items" gorm:"many2many:orders_items"`
}

func (o *Order) TableName() string {
	return "orders"
}

type Item struct {
	ID      uint   `json:"-" gorm:"primary_key"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type OrdersResponse struct {
	Orders []*OrdersResponseOrder `json:"orders"`
	Total  int64                  `json:"total"`
	Offset int64                  `json:"offset"`
}

type OrdersResponseOrder struct {
	ID     uint    `json:"id"`
	Status string  `json:"status"`
	Items  []*Item `json:"items"`
}
