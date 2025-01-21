package adyen

import (
	"context"
	"fmt"

	"github.com/adyen/adyen-go-api-library/v5/src/adyen"
	"github.com/adyen/adyen-go-api-library/v5/src/checkout"
	"github.com/adyen/adyen-go-api-library/v5/src/common"
)

type Client struct {
	checkout *checkout.Checkout
	config   *Config
}

type Config struct {
	ApiKey      string
	Environment string
	MerchantID  string
	ClientKey   string
	ReturnURL   string
}

func NewClient(cfg *Config) (*Client, error) {
	env := common.Environment(cfg.Environment)
	client := adyen.NewClient(&common.Config{
		ApiKey:      cfg.ApiKey,
		Environment: env,
	})
	return &Client{
		checkout: client.Checkout,
		config:   cfg,
	}, nil
}

type PaymentRequest struct {
	Amount      float64
	Currency    string
	Reference   string
	Description string
	ReturnURL   string
}

type PaymentResponse struct {
	SessionData string                 `json:"session_data"`
	ClientKey   string                 `json:"client_key"`
	Config      map[string]interface{} `json:"config"`
}

func (c *Client) CreatePaymentSession(req *PaymentRequest) (*PaymentResponse, error) {
	amount := checkout.Amount{
		Currency: req.Currency,
		Value:    int64(req.Amount * 100), // Convert to cents
	}

	returnURL := req.ReturnURL
	if returnURL == "" {
		returnURL = c.config.ReturnURL
	}

	request := &checkout.PaymentSetupRequest{
		MerchantAccount: c.config.MerchantID,
		Amount:          amount,
		Reference:       req.Reference,
		ReturnUrl:       returnURL,
		CountryCode:     "US", // You might want to make this configurable
		Channel:         "Web",
	}

	ctx := context.Background()
	session, httpResp, err := c.checkout.PaymentSession(request, ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create payment session: %w\nhttp response: %v",
			err, httpResp,
		)
	}

	return &PaymentResponse{
		SessionData: session.PaymentSession,
		ClientKey:   c.config.ClientKey,
		Config: map[string]interface{}{
			"environment": c.config.Environment,
			"clientKey":   c.config.ClientKey,
			"locale":      "en-US",
		},
	}, nil
}

func (c *Client) GetPaymentDetails(paymentID string) (*checkout.PaymentDetailsResponse, error) {
	request := &checkout.DetailsRequest{
		PaymentData: paymentID,
	}

	ctx := context.Background()
	details, httpResp, err := c.checkout.PaymentsDetails(request, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment details: %w\nhttp response: %v",
			err, httpResp,
		)
	}

	return &details, nil
}

func (c *Client) HandleWebhook(notification map[string]interface{}) error {
	// Add HMAC validation if using webhooks
	// Process the notification based on eventCode
	eventCode := notification["eventCode"].(string)
	success := notification["success"].(string)

	switch eventCode {
	case "AUTHORISATION":
		if success == "true" {
			// Payment was successful
			// Update order status
		} else {
			// Payment failed
			// Handle failure
		}
	case "CANCELLATION":
		// Handle cancellation
	case "REFUND":
		// Handle refund
	}

	return nil
}
