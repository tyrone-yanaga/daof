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
	return args.Error(0)
}

// Read mocks the Read method
func (m *MockOdooClient) Read(model string, ids []int64, options *go_odoo.Options, result interface{}) error {
	args := m.Called(model, ids, options, result)
	return args.Error(0)
}
