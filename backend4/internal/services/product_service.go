package services

import (
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type ProductService struct {
	odooClient odoo.OdooClient
	imageCache *cache.Cache
}

func NewProductService(odooClient odoo.OdooClient) *ProductService {
	return &ProductService{
		odooClient: odooClient,
		imageCache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	var odooProducts []odoo.OdooProductTemplate
	criteria := s.odooClient.NewCriteria().Add("active", "=", true)
	options := s.odooClient.NewOptions().FetchFields(
		"name", "description", "list_price", "default_code", "active",
		"image", "image_medium", "image_small",
	)

	err := s.odooClient.SearchRead("product.template", criteria, options, &odooProducts)

	if err != nil {
		return nil, err
	}
	println("ProductService GetProducts: ", odooProducts) // Debug print
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

func (s *ProductService) GetProductImage(productID string) ([]byte, error) {

	// if cachedImage, exists := s.imageCache.Get(productID); exists {
	// 	return cachedImage.([]byte), nil
	// }
	if productID == "" {
		return nil, fmt.Errorf("product ID cannot be empty")
	}
	var products []odoo.OdooProductTemplate

	// Set up options to fetch only image fields
	fields := []string{"image_1920", "image_1024", "image_128"}

	// Create criteria to fetch specific product
	criteria := s.odooClient.NewCriteria()
	idInt, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}
	fmt.Printf("Debug - Fetching image for product ID: %d\n", idInt)
	criteria.Add("id", "=", idInt)
	// Fetch the product with image
	err = s.odooClient.SearchReads(
		"product.template",
		idInt,
		fields,
		&products,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product image: %w", err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("product not found: %s", productID)
	}

	fmt.Printf("Debug - Found product: %+v\n", products[0])
	var imageData string
	if products[0].Image1920 != "" {
		imageData = products[0].Image1920
		fmt.Printf("Debug - Using Image1920 for product %s\n", productID)
	} else if products[0].Image1024 != "" {
		imageData = products[0].Image1024
		fmt.Printf("Debug - Using Image1024 for product %s\n", productID)
	} else if products[0].Image128 != "" {
		imageData = products[0].Image128
		fmt.Printf("Debug - Using Image128 for product %s\n", productID)
	}

	if imageData == "" {
		return nil, fmt.Errorf("no image data found for product: %s", productID)
	}
	fmt.Printf("PS - Base64 image data length: %d\n", len(products[0].Image1920))
	fmt.Printf("PS - Base64 image data length: %d\n", len(products[0].Image1024))
	fmt.Printf("PS - Base64 image data length: %d\n", len(products[0].Image128))

	// Odoo typically stores images as base64 strings
	// Decode the base64 string back to bytes
	imageBytes, err := base64.StdEncoding.DecodeString(products[0].Image1920)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image data: %w", err)
	}

	// s.imageCache.Set(productID, imageBytes, cache.DefaultExpiration)

	return imageBytes, nil
}
