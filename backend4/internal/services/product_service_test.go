package services_test

import (
	"ecommerce/internal/services"
	"ecommerce/internal/testutils/mocks"
	"errors"
	"testing"

	go_odoo "github.com/skilld-labs/go-odoo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProducts(t *testing.T) {
	t.Run("successful fetch", func(t *testing.T) {
		mockClient := new(mocks.MockOdooClient)
		service := services.NewProductService(mockClient)

		criteria := go_odoo.NewCriteria()
		options := go_odoo.NewOptions()

		mockClient.On("NewCriteria").Return(criteria)
		mockClient.On("NewOptions").Return(options)
		mockClient.On("SearchRead",
			"product.template",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("*[]odoo.OdooProductTemplate")).Return(nil)

		products, err := service.GetProducts()
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, int64(1), products[0].OdooID)
		assert.Equal(t, int64(2), products[1].OdooID)
		mockClient.AssertExpectations(t)
	})

	t.Run("fetch error", func(t *testing.T) {
		mockClient := new(mocks.MockOdooClient)
		service := services.NewProductService(mockClient)

		criteria := go_odoo.NewCriteria()
		options := go_odoo.NewOptions()
		expectedError := errors.New("connection error")

		mockClient.On("NewCriteria").Return(criteria)
		mockClient.On("NewOptions").Return(options)
		mockClient.On("SearchRead",
			"product.template",
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("*[]odoo.OdooProductTemplate")).Return(expectedError)

		products, err := service.GetProducts()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, products)
		mockClient.AssertExpectations(t)
	})
}
