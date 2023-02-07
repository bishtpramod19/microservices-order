package ports

import "github.com/bishtpramod19/microservices/order/internal/application/core/domain"

type ApiPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
