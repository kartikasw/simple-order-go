package service

import (
	"simple-order-go/internal/entity"
	"simple-order-go/internal/repository"
)

type OrderService struct {
	orderRepo repository.IOrderRepository
}

type IOrderService interface {
	CreateOrder(order entity.OrderViewModel) (entity.OrderViewModel, error)
	GetOrder(orderID int64) (entity.OrderViewModel, error)
	GetAllOrders() ([]entity.OrderViewModel, error)
	UpdateOrder(order entity.OrderViewModel) error
	DeleteOrder(orderID int64) error
}

func NewOrderService(orderRepo repository.IOrderRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo}

}

func (s *OrderService) CreateOrder(order entity.OrderViewModel) (entity.OrderViewModel, error) {
	result, err := s.orderRepo.CreateOrder(order.ToEntity())
	if err != nil {
		return entity.OrderViewModel{}, err
	}

	return result.ToViewModel(), nil
}

func (s *OrderService) GetOrder(orderID int64) (entity.OrderViewModel, error) {
	result, err := s.orderRepo.GetOrder(orderID)
	if err != nil {
		return entity.OrderViewModel{}, err
	}

	return result.ToViewModel(), nil
}

func (s *OrderService) GetAllOrders() ([]entity.OrderViewModel, error) {
	result, err := s.orderRepo.GetAllOrders()
	if err != nil {
		return []entity.OrderViewModel{}, err
	}

	return result.ToViewModel(), nil
}

func (s *OrderService) UpdateOrder(order entity.OrderViewModel) error {
	err := s.orderRepo.UpdateOrder(order.ToEntity())
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) DeleteOrder(orderID int64) error {
	err := s.orderRepo.DeleteOrder(orderID)
	if err != nil {
		return err
	}

	return nil
}
