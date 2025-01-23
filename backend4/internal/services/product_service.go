package services

import (
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"fmt"
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
	var odooProducts []odoo.OdooProductTemplate
	criteria := s.odooClient.NewCriteria().Add("active", "=", true)
	options := s.odooClient.NewOptions().FetchFields("name", "description", "list_price", "default_code", "active")

	err := s.odooClient.SearchRead("product.template", criteria, options, &odooProducts)
	fmt.Println("-----error---- ", err)

	if err != nil {
		return nil, err
	}

	var products []models.Product
	for _, op := range odooProducts {
		products = append(products, models.Product{
			OdooID:    op.ID,
			Name:      op.Name,
			BasePrice: op.ListPrice,
			Active:    true,
		})
	}
	return products, nil
}

func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	// Implement single product fetching from Odoo
	// productID, err := strconv.Atoi(id)
	// if err != nil {
	// 	return nil, err
	// }
	// options := s.odooClient.NewOptions().
	// 	FetchFields(
	// 		"name", "description", "list_price", "qty_available", "default_code",
	// 	)
	var product map[string]interface{}
	fmt.Println("ProductService GetProduct: ", id) // Debug print
	// err = s.odooClient.Read("product.template", []int64{int64(productID)}, options, &product)

	// if err != nil {
	// 	return nil, err
	// }

	productReturn := models.FromOdooProduct(product)
	return &productReturn, nil
}
