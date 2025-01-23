// File: internal/services/product_service_test.go

package services_test

import (
	"testing"

	"ecommerce/internal/services"
	"ecommerce/internal/testutils/fixtures"
	"ecommerce/internal/testutils/mocks"

	go_odoo "github.com/skilld-labs/go-odoo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProducts(t *testing.T) {
	// Arrange
	mockOdoo := new(mocks.MockOdooClient)
	service := services.NewProductService(mockOdoo)
	testProducts := fixtures.GetTestProducts()

	// Setup mock expectations
	mockOdoo.On("NewCriteria").Return(go_odoo.NewCriteria())
	mockOdoo.On("NewOptions").Return(go_odoo.NewOptions())
	mockOdoo.On("SearchRead", "product.product", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			result := args.Get(3).(*[]map[string]interface{})
			*result = testProducts
		}).
		Return(nil)

	// Act
	products, err := service.GetProducts()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, products, len(testProducts))
	assert.Equal(t, "Test Product 1", products[0].Name)
	mockOdoo.AssertExpectations(t)
}

// ... rest of your test cases ...
