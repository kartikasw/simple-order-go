package entity

import (
	"time"
)

type Items []Item

type ItemViewModels []ItemViewModel

type Item struct {
	ID          int64     `gorm:"primary_key;column:id;autoIncrement"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description;unique;<-:create"`
	Quantity    int32     `gorm:"column:quantity"`
	OrderID     int64     `gorm:"index;column:order_id"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt   time.Time `gorm:"column:updated_at;autoCreateTime"`
}

type ItemViewModel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int32  `json:"quantity"`
}

func (e Item) toViewModel() ItemViewModel {
	return ItemViewModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Quantity:    e.Quantity,
	}
}

func itemListToViewModel(e Items) []ItemViewModel {
	items := make([]ItemViewModel, len(e))
	for i, itemEntity := range e {
		item := itemEntity.toViewModel()
		items[i] = item
	}

	return items
}

func (vm ItemViewModel) toEntity(orderID int64) Item {
	return Item{
		ID:          vm.ID,
		Name:        vm.Name,
		Description: vm.Description,
		Quantity:    vm.Quantity,
		OrderID:     orderID,
	}
}

func itemViewModelListToEntity(orderID int64, vm ItemViewModels) []Item {
	items := make([]Item, len(vm))
	for i, itemViewModel := range vm {
		item := itemViewModel.toEntity(orderID)
		items[i] = item
	}

	return items
}
