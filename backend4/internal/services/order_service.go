package services

import (
	"context"
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type OrderService struct {
	rabbitmqConn *amqp.Connection
	odooClient   *odoo.Client
}

func NewOrderService(rabbitmqConn *amqp.Connection, odooClient *odoo.Client) *OrderService {
	return &OrderService{
		rabbitmqConn: rabbitmqConn,
		odooClient:   odooClient,
	}
}

func (s *OrderService) CreateOrder(order *models.Order) (*models.Order, error) {
	// Create order in Odoo
	orderData := map[string]interface{}{
		"partner_id":   order.UserID,
		"state":        "draft",
		"amount_total": order.Total,
	}

	odooOrderID, err := s.odooClient.Create("sale.order", orderData)
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

	order.ID = uint(odooOrderID)
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
	if order.OdooID > 0 {
		odooOrder, err := s.getOdooOrder(order.OdooID)
		if err != nil {
			// Log the error but don't fail the request
			// We can still return the local order data
			fmt.Printf("failed to fetch order from Odoo: %v\n", err)
		} else {
			// Update local order status if it differs
			if odooOrder["state"] != order.Status {
				order.Status = odooOrder["state"].(string)
				s.db.Save(&order)
			}
		}
	}

	return &order, nil
}

// Helper method to fetch order from Odoo
func (s *OrderService) getOdooOrder(odooID int64) (map[string]interface{}, error) {
	criteria := &odoo.Criteria{
		Filters: [][]odoo.Filter{
			{
				{
					Field:    "id",
					Operator: "=",
					Value:    odooID,
				},
			},
		},
	}

	options := &odoo.Options{
		Fields: []string{
			"name",
			"state",
			"amount_total",
			"date_order",
			"partner_id",
		},
	}

	var result []map[string]interface{}
	err := s.odooClient.SearchRead("sale.order", criteria, options, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order from Odoo: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("order not found in Odoo")
	}

	return result[0], nil
}
