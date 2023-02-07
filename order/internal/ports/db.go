package ports

import "github.com/bishtpramod19/microservices/order/internal/application/core/domain"

type DBPort interface {
	Save(*domain.Order) error
	Get(id string) (domain.Order, error)
}
