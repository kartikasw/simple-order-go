package repository

import (
	"simple-order-go/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	CreateOrder(order entity.Order) (entity.Order, error)
	GetOrder(orderID int64) (entity.Order, error)
	GetAllOrders() (entity.Orders, error)
	UpdateOrder(order entity.Order) error
	DeleteOrder(orderID int64) error
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order entity.Order) (entity.Order, error) {
	err := r.db.Create(&order).Error
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) GetOrder(orderID int64) (order entity.Order, err error) {
	err = r.db.Model(&entity.Order{}).Preload("Items").Take(&order, "orders.id = ?", orderID).Error
	return
}

func (r *OrderRepository) GetAllOrders() (entity.Orders, error) {
	var orders []entity.Order
	err := r.db.Unscoped().Model(&entity.Order{}).Preload("Items").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateOrder(order entity.Order) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").Take(&entity.Order{}, "id = ?", order.ID).Error
		if err != nil {
			return err
		}

		err = tx.Omit("Items").Save(&order).Error
		if err != nil {
			return err
		}

		for _, item := range order.Items {
			var exist entity.Item
			err = tx.Where("description = ?", item.Description).
				Attrs(entity.Item{
					Name:        item.Name,
					Description: item.Description,
					Quantity:    item.Quantity,
					OrderID:     order.ID,
				}).
				FirstOrCreate(&exist).
				Error
			if err != nil {
				return err
			}

			if exist.ID != 0 {
				err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&exist).
					Where("id = ?", exist.ID).
					Updates(entity.Item{Name: item.Name, Quantity: item.Quantity}).
					Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (r *OrderRepository) DeleteOrder(orderID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Delete(&entity.Order{}, orderID).Error; err != nil {
			return err
		}
		return nil
	})
}
