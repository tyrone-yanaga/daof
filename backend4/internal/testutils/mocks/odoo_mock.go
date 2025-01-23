package mocks

import (
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"time"

	go_odoo "github.com/skilld-labs/go-odoo"
	"github.com/stretchr/testify/mock"
)

// MockOdooClient is a mock implementation of the Odoo client interface
type MockOdooClient struct {
	mock.Mock
}

var _ odoo.OdooClient = (*MockOdooClient)(nil)

// Create mocks the Create method
func (m *MockOdooClient) Create(model string, data []interface{}, options *go_odoo.Options) ([]int64, error) {
	args := m.Called(model, data, options)
	return args.Get(0).([]int64), args.Error(1)
}

// NewCriteria mocks the NewCriteria method
func (m *MockOdooClient) NewCriteria() *go_odoo.Criteria {
	args := m.Called()
	if args.Get(0) == nil {
		return go_odoo.NewCriteria()
	}
	return args.Get(0).(*go_odoo.Criteria)
}

// NewOptions mocks the NewOptions method
func (m *MockOdooClient) NewOptions() *go_odoo.Options {
	args := m.Called()
	if args.Get(0) == nil {
		return go_odoo.NewOptions()
	}
	return args.Get(0).(*go_odoo.Options)
}

// SearchRead mocks the SearchRead method
func (m *MockOdooClient) SearchRead(model string, criteria *go_odoo.Criteria, options *go_odoo.Options, result interface{}) error {
	args := m.Called(model, criteria, options, result)
	if products, ok := result.(*[]models.Product); ok {
		testTime := time.Date(2025, time.January, 23, 2, 24, 11, 379642000, time.Local)
		*products = []models.Product{
			{
				OdooID:      1,
				Name:        "Test Product 1",
				Description: "Test Description 1",
				BasePrice:   99.99,
				SKU:         "TEST001",
				Active:      true,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
			{
				OdooID:      2,
				Name:        "Test Product 2",
				Description: "Test Description 2",
				BasePrice:   149.99,
				SKU:         "TEST002",
				Active:      true,
				CreatedAt:   testTime,
				UpdatedAt:   testTime,
			},
		}
	}
	return args.Error(0)
}

// Read mocks the Read method
func (m *MockOdooClient) Read(model string, ids []int64, options *go_odoo.Options, result interface{}) error {
	args := m.Called(model, ids, options, result)
	return args.Error(0)
}
