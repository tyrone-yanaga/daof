// internal/odoo/client.go
package odoo

import (
	"context"
	"net/http"
	"time"
)

type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		baseURL:  cfg.BaseURL,
		username: cfg.Username,
		password: cfg.Password,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SyncInventory(ctx context.Context, productID string) (int, error) {
	// Implement Odoo XML-RPC call to get inventory levels
	// Return current inventory level
	return 0, nil
}
