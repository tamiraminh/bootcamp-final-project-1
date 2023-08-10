package order

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
)

type OrderService interface {
	CreateOrder(order Order) (err error)
	CreateOrderItem(order OrderItem) (err error)
}

type OrderServiceImpl struct {
	OrderRepository OrderRepository
	Config          *configs.Config
}

func ProvideOrderServiceImpl(orderRepository OrderRepository, config *configs.Config) *OrderServiceImpl {
	s := new(OrderServiceImpl)
	s.OrderRepository = orderRepository
	s.Config = config

	return s
}

func (s *OrderServiceImpl) CreateOrder(order Order) (err error)  {
	err = s.OrderRepository.CreateOrder(order)
	if err != nil {
		return failure.BadRequest(err)
	}
	
	return 
}

func (s *OrderServiceImpl) CreateOrderItem(orderItem OrderItem) (err error)  {
	err = s.OrderRepository.CreateOrderItem(orderItem)
	if err != nil {
		return failure.BadRequest(err)
	}
	
	return 
}