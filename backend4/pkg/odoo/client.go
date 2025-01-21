package odoo

import (
	"fmt"

	"github.com/skilld-labs/go-odoo"
)

type Client struct {
	*odoo.Client
}

type Config struct {
	URL      string
	Database string
	Username string
	Password string
}

type Product struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ListPrice    float64 `json:"list_price"`
	QtyAvailable float64 `json:"qty_available"`
	DefaultCode  string  `json:"default_code"`
}

type ProductList struct {
	Items []Product
}

func NewClient(config Config) (*Client, error) {
	c, err := odoo.NewClient(&odoo.ClientConfig{
		Admin:    config.Username,
		Password: config.Password,
		Database: config.Database,
		URL:      config.URL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create odoo client: %w", err)
	}
	return &Client{c}, nil
}

func (c *Client) GetProducts(offset, limit int) ([]Product, error) {
	criteria := odoo.NewCriteria()
	criteria.Add("active", "=", true)

	options := odoo.NewOptions().
		Offset(offset).
		Limit(limit).
		FetchFields("name", "description", "list_price", "qty_available", "default_code")

	var products ProductList
	err := c.SearchRead("product.template", criteria, options, &products.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products.Items, nil
}

func (c *Client) GetProduct(id int64) (*Product, error) {
	criteria := odoo.NewCriteria()
	criteria.Add("id", "=", id)
	criteria.Add("active", "=", true)

	options := odoo.NewOptions().
		Limit(1).
		FetchFields("name", "description", "list_price", "qty_available", "default_code")

	var products ProductList
	err := c.SearchRead("product.template", criteria, options, &products.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if len(products.Items) == 0 {
		return nil, fmt.Errorf("product not found")
	}
	return &products.Items[0], nil
}

func (c *Client) CreateOrder(orderData map[string]interface{}) (int64, error) {
	ids, err := c.Create("sale.order", []interface{}{orderData}, odoo.NewOptions())
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	if len(ids) == 0 {
		return 0, fmt.Errorf("no order ID returned")
	}
	return ids[0], nil
}

func (c *Client) SearchRead(model string, criteria *odoo.Criteria, options *odoo.Options, result interface{}) error {
	return c.Client.SearchRead(model, criteria, options, result)
}

func (c *Client) Read(model string, data []int64, options *odoo.Options, result interface{}) error {
	return c.Client.Read(model, data, options, result)
}

func (c Client) NewCriteria() *odoo.Criteria {
	return odoo.NewCriteria()
}

func (c *Client) NewOptions() *odoo.Options {
	return odoo.NewOptions()
}
