package odoo

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/kolo/xmlrpc"
	"github.com/skilld-labs/go-odoo"
)

type Client struct {
	*odoo.Client
	config       Config
	commonClient *xmlrpc.Client
	objectClient *xmlrpc.Client
	uid          int
	transport    *http.Transport
}

type Config struct {
	URL           string
	Database      string
	Username      string
	Password      string
	MaxRetries    int
	RetryInterval time.Duration
}

type Product struct {
	ID          int64       `xmlrpc:"id"`
	Name        string      `xmlrpc:"name"`
	Description interface{} `xmlrpc:"description"` // Handle potential null/string
	ListPrice   float64     `xmlrpc:"list_price"`
	DefaultCode interface{} `xmlrpc:"default_code"` // Handle potential null/string
	Active      interface{} `xmlrpc:"active"`
}

type OdooProductTemplate struct {
	ID          int64       `xmlrpc:"id"`
	Name        string      `xmlrpc:"name"`
	Description interface{} `xmlrpc:"description"` // Handle potential null/string
	ListPrice   float64     `xmlrpc:"list_price"`
	DefaultCode interface{} `xmlrpc:"default_code"` // Handle potential null/string
	Active      interface{} `xmlrpc:"active"`       // Handle potential string/bool

	Image1920 string `xmlrpc:"image_1920"`
	Image1024 string `xmlrpc:"image_1024"`
	Image128  string `xmlrpc:"image_128"`
}
type ProductList struct {
	Items []Product
}

// OdooClient defines the interface for Odoo operations
type OdooClient interface {
	Create(model string, data []interface{}, options *odoo.Options) ([]int64, error)
	NewCriteria() *odoo.Criteria
	NewOptions() *odoo.Options
	SearchRead(model string, criteria *odoo.Criteria, options *odoo.Options, result interface{}) error
	SearchReads(model string, criteria interface{}, options []string, result interface{}) error
	Read(model string, ids []int64, options *odoo.Options, result interface{}) error
}

// Ensure Client implements OdooClient
var _ OdooClient = (*Client)(nil)

func NewClient(config Config) (*Client, error) {
	if config.MaxRetries == 0 {
		config.MaxRetries = 10
	}
	if config.RetryInterval == 0 {
		config.RetryInterval = 5 * time.Second
	}

	var client *Client
	var err error

	for i := 0; i < config.MaxRetries; i++ {
		log.Printf("Attempting to connect to Odoo (attempt %d/%d)", i+1, config.MaxRetries)

		// Try to create a test connection
		resp, err := http.Get(config.URL)
		if err != nil {
			log.Printf("Failed to connect to Odoo: %v", err)
			time.Sleep(config.RetryInterval)
			continue
		}
		resp.Body.Close()

		// If we can connect, try to initialize the client
		client, err = initClient(config)
		if err == nil {
			log.Printf("Successfully connected to Odoo")
			break
		}

		log.Printf("Failed to initialize Odoo client: %v", err)
		time.Sleep(config.RetryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create odoo client after %d attempts: %w",
			config.MaxRetries, err)
	}

	return client, nil
}

func initClient(config Config) (*Client, error) {
	// Create common client for authentication
	commonClient, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/common", config.URL), &http.Transport{
		ResponseHeaderTimeout: 60 * time.Second,
		DisableKeepAlives:     false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create common client: %w", err)
	}

	// Create object client for method calls
	objectClient, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", config.URL), &http.Transport{
		ResponseHeaderTimeout: 60 * time.Second,
		DisableKeepAlives:     false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create object client: %w", err)
	}

	baseConfig := &odoo.ClientConfig{
		Admin:    config.Username,
		Password: config.Password,
		Database: config.Database,
		URL:      config.URL,
	}

	// Initialize base odoo client first
	baseClient, err := odoo.NewClient(baseConfig)
	if err != nil {
		return nil, fmt.Errorf("new odoo client creation failed: %w", err)
	}

	client := &Client{
		Client:       baseClient, // Set the base client
		config:       config,
		commonClient: commonClient,
		objectClient: objectClient,
	}

	// Authenticate and get user ID
	var result int
	err = commonClient.Call("authenticate", []interface{}{
		config.Database,
		config.Username,
		config.Password,
		map[string]interface{}{},
	}, &result)

	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if result == 0 {
		return nil, fmt.Errorf("authentication failed: invalid credentials")
	}

	client.uid = result

	return client, nil
}

// Close closes the client connections
func (c *Client) Close() error {
	if c.transport != nil {
		c.transport.CloseIdleConnections()
	}
	c.commonClient = nil
	c.objectClient = nil
	c.transport = nil
	return nil
}

func (c *Client) GetProducts(offset, limit int) ([]Product, error) {
	criteria := odoo.NewCriteria()
	criteria.Add("active", "=", true)

	options := odoo.NewOptions().
		Offset(offset).
		Limit(limit).
		FetchFields("name", "description", "list_price", "qty_available", "default_code")

	var products ProductList
	err := c.SearchRead("product.product", criteria, options, &products.Items)
	fmt.Printf("Odoo CLIENT GetProducts call: %+v\n", products) // Debug print

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
	err := c.SearchRead("product.product", criteria, options, &products.Items)
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
	var records []OdooProductTemplate // Changed from []odoo.ProductTemplate
	fields := []string{"id", "name", "description", "list_price", "default_code", "active"}

	err := c.objectClient.Call("execute_kw", []interface{}{
		c.config.Database,
		c.uid,
		c.config.Password,
		model,
		"search_read",
		[]interface{}{[]interface{}{}},
		map[string]interface{}{
			"fields": fields,
		},
	}, &records)

	if err != nil {
		return fmt.Errorf("search_read failed: %w", err)
	}
	fmt.Printf("Odoo CLIENT SearchRead records: %+v\n", records) // Debug print

	// Convert the result to the expected type
	resultValue := reflect.ValueOf(result).Elem()
	resultValue.Set(reflect.ValueOf(records))
	fmt.Printf("Odoo CLIENT SearchRead result: %+v\n", result) // Debug print

	return nil
}
func (c *Client) SearchReads(model string, productID interface{}, options []string, result interface{}) error {
	var records []map[string]interface{}

	if productID == nil {
		return fmt.Errorf("no product ID found in criteria")
	}

	// Create the domain with the specific product ID
	domain := []interface{}{
		[]interface{}{"id", "=", productID},
	}

	fmt.Printf("Debug - Searching for product ID: %v\n", productID)
	err := c.objectClient.Call("execute_kw", []interface{}{
		c.config.Database,
		c.uid,
		c.config.Password,
		model,
		"search_read",
		[]interface{}{domain},
		map[string]interface{}{
			"fields": options,
		},
	}, &records)

	if err != nil {
		fmt.Printf("Search read error: %v\n", err)
		return fmt.Errorf("search_reads failed: %w", err)
	}

	fmt.Printf("Debug - Raw records: %+v\n", records)

	var templates []OdooProductTemplate
	for idx, record := range records {
		fmt.Printf("Debug - Processing record %d: %+v\n", idx, record)

		id := toInt64(record["id"])
		fmt.Printf("Debug - Converted ID: %v\n", id)

		template := OdooProductTemplate{
			ID:          id,
			Image1920:   toString(record["image_1920"]),
			Image1024:   toString(record["image_1024"]),
			Image128:    toString(record["image_128"]),
			DefaultCode: toString(record["default_code"]),
			Active:      toBoolInterface(record["active"]),
		}
		templates = append(templates, template)
	}

	// Set the result
	resultValue := reflect.ValueOf(result).Elem()
	resultValue.Set(reflect.ValueOf(templates))

	return nil
}

// Updated helper functions with better nil handling
func toInt64(v interface{}) int64 {
	if v == nil {
		fmt.Printf("Debug - Received nil value for int64 conversion\n")
		return 0
	}

	fmt.Printf("Debug - Converting value type %T: %v\n", v, v)

	switch v := v.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	fmt.Printf("Debug - Unsupported type for int64 conversion: %T\n", v)
	return 0
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func toBoolInterface(v interface{}) interface{} {
	if v == nil {
		return false
	}
	switch v := v.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1" || v == "yes"
	case int:
		return v != 0
	case float64:
		return v != 0
	default:
		return false
	}
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

// Options represents query options for Odoo API calls
type Options struct {
	*odoo.Options
}

// Criteria represents search criteria for Odoo API calls
type Criteria struct {
	*odoo.Criteria
}

// NewOptions creates a new Options instance
func NewOptions() *Options {
	return &Options{
		Options: odoo.NewOptions(),
	}
}

// NewCriteria creates a new Criteria instance
func NewCriteria() *Criteria {
	return &Criteria{
		Criteria: odoo.NewCriteria(),
	}
}
