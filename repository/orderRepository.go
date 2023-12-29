package repository

import (
	"awesomeProject/models"
	"gorm.io/gorm"
)

// OrderRepository представляет собой репозиторий для работы с заказами в базе данных.
type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetOrder(orderUID string) (*models.Order, error) {
	order := &models.Order{}
	err := r.db.Where("order_uid = ?", orderUID).Preload("Delivery").Preload("Payment").Preload("Items").First(order).Error
	return order, err
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) DeleteOrder(orderUID string) error {
	return r.db.Where("order_uid = ?", orderUID).Delete(models.Order{}).Error
}
