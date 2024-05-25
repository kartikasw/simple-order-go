package entity

import (
	"time"
)

type Orders []Order

type Order struct {
	ID           int64     `gorm:"primary_key;column:id;autoIncrement"`
	CustomerName string    `gorm:"column:customer_name"`
	OrderedAt    time.Time `gorm:"column:ordered_at"`
	Items        []Item    `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:updated_at;autoCreateTime"`
}

type OrderViewModel struct {
	ID           int64           `json:"id"`
	CustomerName string          `json:"customer_name"`
	OrderedAt    time.Time       `json:"ordered_at"`
	Items        []ItemViewModel `json:"items"`
}

func (e Order) ToViewModel() OrderViewModel {
	return OrderViewModel{
		ID:           e.ID,
		CustomerName: e.CustomerName,
		Items:        itemListToViewModel(e.Items),
	}
}

func (e Orders) ToViewModel() []OrderViewModel {
	orders := make([]OrderViewModel, len(e))

	for i, order := range e {
		orders[i] = OrderViewModel{
			ID:           order.ID,
			CustomerName: order.CustomerName,
			OrderedAt:    order.OrderedAt,
			Items:        itemListToViewModel(order.Items),
		}
	}

	return orders
}

func (vm OrderViewModel) ToEntity() Order {
	return Order{
		ID:           vm.ID,
		CustomerName: vm.CustomerName,
		OrderedAt:    vm.OrderedAt,
		Items:        itemViewModelListToEntity(int64(vm.ID), vm.Items),
	}
}
