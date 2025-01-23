package mocks

import (
	"ecommerce/pkg/odoo"

	go_odoo "github.com/skilld-labs/go-odoo"
	"github.com/stretchr/testify/mock"
)

// MockOdooClient is a mock implementation of the Odoo client interface
type MockOdooClient struct {
	mock.Mock
}

// Make sure MockCriteria implements the same interface as go_odoo.Criteria
type MockCriteria struct {
	mock.Mock
	*go_odoo.Criteria
}

// Make sure MockOptions implements the same interface as go_odoo.Options
type MockOptions struct {
	mock.Mock
	*go_odoo.Options
}

var _ odoo.OdooClient = (*MockOdooClient)(nil)

// Create mocks the Create method
func (m *MockOdooClient) Create(model string, data []interface{}, options *go_odoo.Options) ([]int64, error) {
	args := m.Called(model, data, options)
	return args.Get(0).([]int64), args.Error(1)
}

// Update MockOdooClient methods to work with the embedded types
func (m *MockOdooClient) NewCriteria() *go_odoo.Criteria {
	args := m.Called()
	if args.Get(0) == nil {
		return go_odoo.NewCriteria()
	}
	return args.Get(0).(*go_odoo.Criteria)
}

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
	err := args.Error(0)
	if err != nil {
		return err
	}
	if products, ok := result.(*[]odoo.OdooProductTemplate); ok {
		*products = []odoo.OdooProductTemplate{
			{
				ID:          int64(1),
				Name:        "Test Product 1",
				Description: "Test Description 1",
				ListPrice:   99.99,
				DefaultCode: "TEST001",
				Active:      true,
			},
			{
				ID:          int64(2),
				Name:        "Test Product 2",
				Description: "Test Description 2",
				ListPrice:   149.99,
				DefaultCode: "TEST002",
				Active:      true,
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

func (m *MockCriteria) Add(field string, operator string, value interface{}) *MockCriteria {
	args := m.Called(field, operator, value)
	return args.Get(0).(*MockCriteria)
}

func (m *MockOptions) FetchFields(fields ...string) *MockOptions {
	// Convert the variadic string parameters to []interface{}
	interfaceSlice := make([]interface{}, len(fields))
	for i, field := range fields {
		interfaceSlice[i] = field
	}

	args := m.Called(interfaceSlice...)
	return args.Get(0).(*MockOptions)
}
