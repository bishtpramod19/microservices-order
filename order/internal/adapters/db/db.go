package db

import (
	"fmt"

	"github.com/bishtpramod19/microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerId int64
	Status     string
	Orderitems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderId     uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openerr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openerr != nil {
		return nil, fmt.Errorf("db connection error : %v", openerr)
	}
	err := db.AutoMigrate(&Order{}, OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("db automigrate error : %v", err)
	}

	return &Adapter{db: db}, nil

}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order
	res := a.db.First(&orderEntity, id)
	var orderItems []domain.OrderItem
	for _, orderItem := range orderEntity.Orderitems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerId,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}

	return order, res.Error

}

func (a Adapter) Save(order *domain.Order) error {
	var orderItems []OrderItem
	for _, orderitem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderitem.ProductCode,
			UnitPrice:   orderitem.UnitPrice,
			Quantity:    orderitem.Quantity,
		})
	}

	orderModel := Order{
		CustomerId: order.CustomerID,
		Status:     order.Status,
		Orderitems: orderItems,
	}

	res := a.db.Create(&orderModel)
	if res.Error == nil {
		order.ID = int64(orderModel.ID)
	}

	return res.Error
}
