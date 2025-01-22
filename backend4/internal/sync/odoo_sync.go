package sync

import (
	"ecommerce/internal/models"
	"fmt"
	"time"

	"ecommerce/pkg/odoo"

	"gorm.io/gorm"
)

type OdooSync struct {
	db         *gorm.DB
	odooClient *odoo.Client
}

func NewOdooSync(db *gorm.DB, odooClient *odoo.Client) *OdooSync {
	return &OdooSync{
		db:         db,
		odooClient: odooClient,
	}
}

// SyncProducts synchronizes products between Odoo and local database
func (s *OdooSync) SyncProducts() error {
	// Create criteria and options for Odoo API
	criteria := s.odooClient.NewCriteria().Add("active", "=", true)
	options := s.odooClient.NewOptions().
		FetchFields("id", "name", "description", "list_price", "qty_available", "default_code")

	var products []interface{}
	if err := s.odooClient.SearchRead("product.template", criteria, options, &products); err != nil {
		return fmt.Errorf("failed to fetch products from Odoo: %w", err)
	}

	// Begin transaction
	tx := s.db.Begin()
	for _, p := range products {
		productMap := p.(map[string]interface{})
		product := models.FromOdooProduct(productMap)

		// Upsert product
		result := tx.Where("odoo_id = ?", product.OdooID).
			Assign(product).
			FirstOrCreate(&product)

		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert product: %w", result.Error)
		}
	}
	return tx.Commit().Error
}

// SyncOrders synchronizes orders from local database to Odoo
func (s *OdooSync) SyncOrders() error {
	// Get unsynchronized orders
	var orders []models.Order
	if err := s.db.Where("odoo_id IS NULL").Find(&orders).Error; err != nil {
		return fmt.Errorf("failed to fetch unsynchronized orders: %w", err)
	}

	for _, order := range orders {
		// Create order in Odoo
		orderData := []interface{}{
			map[string]interface{}{
				"partner_id":   order.UserID, // Assuming UserID maps to Odoo partner_id
				"date_order":   order.CreatedAt.Format(time.RFC3339),
				"state":        "draft",
				"amount_total": order.Total,
			},
		}

		options := s.odooClient.NewOptions()
		odooOrderIDs, err := s.odooClient.Create("sale.order", orderData, options)
		if err != nil {
			return fmt.Errorf("failed to create order in Odoo: %w", err)
		}

		odooOrderID := odooOrderIDs[0]

		// Create order lines
		for _, item := range order.Items {
			lineData := []interface{}{
				map[string]interface{}{
					"order_id":        odooOrderID,
					"product_id":      item.ProductID,
					"product_uom_qty": item.Quantity,
					"price_unit":      item.Price,
				},
			}

			_, err := s.odooClient.Create("sale.order.line", lineData, options)
			if err != nil {
				return fmt.Errorf("failed to create order line in Odoo: %w", err)
			}
		}

		// Update local order with Odoo ID
		if err := s.db.Model(&order).Update("odoo_id", odooOrderID).Error; err != nil {
			return fmt.Errorf("failed to update order with Odoo ID: %w", err)
		}
	}

	return nil
}
