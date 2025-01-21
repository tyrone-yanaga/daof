package services

import (
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
)

type ProductService struct {
	odooClient *odoo.Client
}

func NewProductService(odooClient *odoo.Client) *ProductService {
	return &ProductService{
		odooClient: odooClient,
	}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	// Implement product fetching from Odoo
	products, err := s.odooClient.SearchRead("product.template", []int{}, []string{
		"name", "description", "list_price", "qty_available", "default_code",
	})
	if err != nil {
		return nil, err
	}

	var result []models.Product
	for _, p := range products {
		product := models.Product{
			Name:        p["name"].(string),
			Description: p["description"].(string),
			Price:       p["list_price"].(float64),
			Stock:       int(p["qty_available"].(float64)),
			SKU:         p["default_code"].(string),
		}
		result = append(result, product)
	}

	return result, nil
}

func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	// Implement single product fetching from Odoo
	product, err := s.odooClient.Read("product.template", []int{id}, []string{
		"name", "description", "list_price", "qty_available", "default_code",
	})
	if err != nil {
		return nil, err
	}

	return &models.Product{
		Name:        product["name"].(string),
		Description: product["description"].(string),
		Price:       product["list_price"].(float64),
		Stock:       int(product["qty_available"].(float64)),
		SKU:         product["default_code"].(string),
	}, nil
}
