package services_test

import (
	"ecommerce/internal/services"
	"ecommerce/internal/testutils/mocks"
	"testing"

	"github.com/skilld-labs/go-odoo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProducts(t *testing.T) {
	mockClient := new(mocks.MockOdooClient)
	service := services.NewProductService(mockClient)

	t.Run("successful fetch", func(t *testing.T) {
		mockCriteria := &odoo.Criteria{}
		mockOptions := &odoo.Options{}

		mockClient.On("NewCriteria").Return(mockCriteria)
		mockClient.On("NewOptions").Return(mockOptions)
		mockClient.On("SearchRead",
			"product.product",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("*[]models.Product")).Return(nil)

		products, err := service.GetProducts()

		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, int64(1), products[0].OdooID)
		assert.Equal(t, int64(2), products[1].OdooID)
		mockClient.AssertExpectations(t)
	})
}
