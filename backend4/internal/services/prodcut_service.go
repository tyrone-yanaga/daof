package services

import (
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"strconv"
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

	criteria := s.odooClient.NewCriteria().Add("active", "=", true)

	options := s.odooClient.NewOptions().
		FetchFields(
			"name", "description", "list_price", "qty_available", "default_code",
		)

		// Implement product fetching from Odoo
	var products []map[string]interface{}
	err := s.odooClient.SearchRead("product.template", criteria, options, &products)
	if err != nil {
		return nil, err
	}

	var result []models.Product
	for _, p := range products {
		product := models.FromOdooProduct(p)
		result = append(result, product)
	}

	return result, nil
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
	err = s.odooClient.Read("product.template", []int64{int64(productID)}, options, &product)

	if err != nil {
		return nil, err
	}

	productReturn := models.FromOdooProduct(product)
	return &productReturn, nil
}
