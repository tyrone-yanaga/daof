package services

import (
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"strconv"
)

type ProductService struct {
	odooClient odoo.OdooClient
}

func NewProductService(odooClient odoo.OdooClient) *ProductService {
	return &ProductService{
		odooClient: odooClient,
	}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {

	criteria := s.odooClient.NewCriteria().Add("active", "=", true)

	options := s.odooClient.NewOptions().
		FetchFields(
			"name", "description", "list_price", "qty_available", "default_code",
		)

		// Implement product fetching from Odoo
	var products []models.Product
	err := s.odooClient.SearchRead("product.product", criteria, options, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	// Implement single product fetching from Odoo
	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	options := s.odooClient.NewOptions().
		FetchFields(
			"name", "description", "list_price", "qty_available", "default_code",
		)
	var product map[string]interface{}
	err = s.odooClient.Read("product.product", []int64{int64(productID)}, options, &product)

	if err != nil {
		return nil, err
	}

	productReturn := models.FromOdooProduct(product)
	return &productReturn, nil
}
