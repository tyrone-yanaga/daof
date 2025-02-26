package services

import (
	"context"
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"ecommerce/pkg/queue"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type OrderService struct {
	rabbitmqConn *amqp.Connection
	odooClient   *odoo.Client
	db           *gorm.DB // Added missing db field
}

func NewOrderService(queueClient *queue.Client, odooClient *odoo.Client, db *gorm.DB) *OrderService {
	return &OrderService{
		rabbitmqConn: queueClient.ReturnConnection(),
		odooClient:   odooClient,
		db:           db,
	}
}

func (s *OrderService) CreateOrder(order *models.Order) (*models.Order, error) {
	// Create order in Odoo
	orderData := []interface{}{
		map[string]interface{}{
			"partner_id":   order.UserID,
			"state":        "draft",
			"amount_total": order.Total,
		},
	}

	options := s.odooClient.NewOptions()

	odooOrderIDs, err := s.odooClient.Create("sale.order", orderData, options)
	if err != nil {
		return nil, err
	}

	// Publish order to RabbitMQ for processing
	ch, err := s.rabbitmqConn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	orderBytes, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	err = ch.Publish(
		"orders",
		"new_order",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        orderBytes,
		},
	)
	if err != nil {
		return nil, err
	}

	// Assuming models.Order has ID field of type uint
	if len(odooOrderIDs) > 0 {
		order.ID = uint(odooOrderIDs[0]) // Store the Odoo ID separately
		order.ID = uint(1)               // You might want to generate this differently
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uint) (*models.Order, error) {
	// First try to get from local database
	var order models.Order
	result := s.db.Preload("Items").First(&order, orderID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to fetch order: %w", result.Error)
	}

	// If order has Odoo ID, fetch latest status from Odoo
	if order.ID > 0 {
		odooOrders, err := s.getOdooOrder(int64(order.ID))
		if err != nil {
			// Log the error but don't fail the request
			// We can still return the local order data
			fmt.Printf("failed to fetch order from Odoo: %v\n", err)
		} else {
			// Update local order status if it differs
			if odooOrders[0].Status != order.Status {
				order.Status = odooOrders[0].Status
				s.db.Save(&order)
			}

		}
	}

	return &order, nil
}

// Helper method to fetch order from Odoo
func (s *OrderService) getOdooOrder(odooID int64) ([]models.Order, error) {
	criteria := s.odooClient.NewCriteria().Add("id", "=", odooID)

	options := s.odooClient.NewOptions().
		FetchFields(
			"name",
			"state",
			"amount_total",
			"date_order",
			"partner_id",
		)

	var order []models.Order
	err := s.odooClient.SearchRead("sale.order", criteria, options, &order)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order from Odoo: %w", err)
	}

	if len(order) == 0 {
		return nil, fmt.Errorf("order not found in Odoo")
	}

	return order, nil
}
