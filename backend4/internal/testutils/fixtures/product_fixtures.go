package fixtures

// ProductFixtures provides test data for product tests
var ProductFixtures = struct {
	SingleProduct  map[string]interface{}
	ProductList    []map[string]interface{}
	InvalidProduct map[string]interface{}
}{
	SingleProduct: map[string]interface{}{
		"id":            int64(1),
		"name":          "Test Product",
		"description":   "Test Description",
		"list_price":    99.99,
		"qty_available": 10.0,
		"default_code":  "TEST001",
		"active":        true,
	},
	ProductList: []map[string]interface{}{
		{
			"id":            int64(1),
			"name":          "Test Product 1",
			"description":   "Test Description 1",
			"list_price":    99.99,
			"qty_available": 10.0,
			"default_code":  "TEST001",
			"active":        true,
		},
		{
			"id":            int64(2),
			"name":          "Test Product 2",
			"description":   "Test Description 2",
			"list_price":    149.99,
			"qty_available": 5.0,
			"default_code":  "TEST002",
			"active":        true,
		},
		{
			"id":            int64(3),
			"name":          "Test Product 3",
			"description":   "Test Description 3",
			"list_price":    199.99,
			"qty_available": 0.0,
			"default_code":  "TEST003",
			"active":        true,
		},
	},
	InvalidProduct: map[string]interface{}{
		"id":            "invalid", // Invalid type for ID
		"name":          123,       // Invalid type for name
		"list_price":    "invalid", // Invalid type for price
		"qty_available": "none",    // Invalid type for quantity
	},
}

// GetTestProducts returns a copy of the product list to prevent test pollution
func GetTestProducts() []map[string]interface{} {
	products := make([]map[string]interface{}, len(ProductFixtures.ProductList))
	for i, product := range ProductFixtures.ProductList {
		products[i] = make(map[string]interface{})
		for k, v := range product {
			products[i][k] = v
		}
	}
	return products
}

// GetTestProduct returns a copy of the single product to prevent test pollution
func GetTestProduct() map[string]interface{} {
	product := make(map[string]interface{})
	for k, v := range ProductFixtures.SingleProduct {
		product[k] = v
	}
	return product
}
