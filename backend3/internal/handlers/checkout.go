// internal/handlers/checkout.go
package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"your-project/internal/config"
	"your-project/internal/models"
	"your-project/internal/services"

	"github.com/adyen/adyen-go-api-library/v4/adyen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CheckoutHandler struct {
	checkoutService services.CheckoutService
	cartService     services.CartService
	adyenClient     *adyen.Client
	config          *config.Config
}

type CheckoutSession struct {
	SessionData     string            `json:"sessionData"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Items           []models.CartItem `json:"items"`
	ShippingAddress models.Address    `json:"shippingAddress"`
	BillingAddress  models.Address    `json:"billingAddress"`
	CustomerEmail   string            `json:"customerEmail"`
	CreatedAt       time.Time         `json:"createdAt"`
}

func NewCheckoutHandler(cs services.CheckoutService, cartService services.CartService, config *config.Config) *CheckoutHandler {
	client := adyen.NewClient(&adyen.Config{
		ApiKey:      config.Adyen.APIKey,
		Environment: config.Adyen.Environment,
	})

	return &CheckoutHandler{
		checkoutService: cs,
		cartService:     cartService,
		adyenClient:     client,
		config:          config,
	}
}

// InitCheckout initializes a new checkout session
func (h *CheckoutHandler) InitCheckout(c *gin.Context) {
	var req struct {
		ShippingAddress models.Address `json:"shippingAddress" binding:"required"`
		BillingAddress  models.Address `json:"billingAddress" binding:"required"`
		CustomerEmail   string         `json:"customerEmail" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartID := getCartIDFromSession(c)
	if cartID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active cart found"})
		return
	}

	// Get cart and validate items
	cart, err := h.cartService.Get(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Validate stock availability
	if err := h.validateStock(c.Request.Context(), cart.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate totals
	amount, tax := h.calculateTotals(cart.Items, req.ShippingAddress.Country)

	// Create Adyen session
	reference := generateReference()
	session, err := h.createPaymentSession(c.Request.Context(), amount, tax, cart.Items, req.ShippingAddress, reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create payment session: %v", err)})
		return
	}

	// Store checkout session
	checkoutSession := CheckoutSession{
		SessionData:     session.SessionData,
		Amount:          amount,
		Currency:        "USD",
		Items:           cart.Items,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		CustomerEmail:   req.CustomerEmail,
		CreatedAt:       time.Now(),
	}

	if err := h.checkoutService.StoreCheckoutSession(c.Request.Context(), cartID, checkoutSession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessionId":   session.ID,
		"sessionData": session.SessionData,
		"amount":      amount,
		"currency":    "USD",
		"reference":   reference,
	})
}

// CompleteCheckout handles the payment completion
func (h *CheckoutHandler) CompleteCheckout(c *gin.Context) {
	var req struct {
		PaymentData string `json:"paymentData" binding:"required"`
		SessionID   string `json:"sessionId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartID := getCartIDFromSession(c)
	if cartID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active cart found"})
		return
	}

	// Retrieve checkout session
	session, err := h.checkoutService.GetCheckoutSession(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid checkout session"})
		return
	}

	// Verify payment with Adyen
	payment, err := h.verifyPayment(c.Request.Context(), req.PaymentData, session)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Payment verification failed: %v", err)})
		return
	}

	// Create order
	order := models.Order{
		ID:              uuid.New().String(),
		CartID:          cartID,
		Items:           session.Items,
		TotalAmount:     session.Amount,
		PaymentID:       payment.PspReference,
		CustomerEmail:   session.CustomerEmail,
		ShippingAddress: session.ShippingAddress,
		BillingAddress:  session.BillingAddress,
		Status:          "processing",
		CreatedAt:       time.Now(),
	}

	order, err = h.checkoutService.CreateOrder(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Update inventory
	go h.processOrderInventory(context.Background(), order)

	// Clear cart
	c.SetCookie("cart_id", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"orderId": order.ID,
		"status":  "success",
		"message": "Order placed successfully",
	})
}

// WebhookHandler processes Adyen webhooks
func (h *CheckoutHandler) WebhookHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Verify HMAC signature
	if !h.verifyHMAC(c.Request.Header.Get("X-Adyen-Signature"), body) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid HMAC signature"})
		return
	}

	var notification models.AdyenNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification format"})
		return
	}

	// Process each notification item
	for _, item := range notification.NotificationItems {
		if err := h.processNotificationItem(c.Request.Context(), item); err != nil {
			log.Printf("Error processing notification item: %v", err)
			continue
		}
	}

	// Adyen expects this specific response format
	c.JSON(http.StatusOK, gin.H{"notificationResponse": "[accepted]"})
}

// Helper functions

func (h *CheckoutHandler) createPaymentSession(ctx context.Context, amount, tax float64, items []models.CartItem, address models.Address, reference string) (*adyen.CheckoutSession, error) {
	lineItems := h.convertToLineItems(items)

	req := &adyen.CheckoutSessionRequest{
		Amount: adyen.Amount{
			Currency: "USD",
			Value:    uint64(amount * 100), // Convert to cents
		},
		CountryCode:      address.Country,
		MerchantAccount:  h.config.Adyen.MerchantAccount,
		Reference:        reference,
		ReturnUrl:        h.config.Server.FrontendURL + "/checkout/complete",
		LineItems:        lineItems,
		ShopperReference: reference,
		Channel:          "Web",
		ShopperLocale:    "en-US",
	}

	return h.adyenClient.Checkout.Sessions(ctx, req)
}

func (h *CheckoutHandler) verifyPayment(ctx context.Context, paymentData string, session CheckoutSession) (*adyen.PaymentVerificationResponse, error) {
	req := &adyen.PaymentVerificationRequest{
		Amount: adyen.Amount{
			Currency: session.Currency,
			Value:    uint64(session.Amount * 100),
		},
		PaymentData: paymentData,
		Reference:   generateReference(),
	}

	return h.adyenClient.Checkout.PaymentsVerification(ctx, req)
}

func (h *CheckoutHandler) verifyHMAC(signature string, payload []byte) bool {
	key, err := hex.DecodeString(h.config.Adyen.WebhookHMACKey)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(payload)
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (h *CheckoutHandler) processNotificationItem(ctx context.Context, item models.NotificationItem) error {
	order, err := h.checkoutService.GetOrderByPaymentID(ctx, item.PspReference)
	if err != nil {
		return err
	}

	switch item.EventCode {
	case "AUTHORISATION":
		if item.Success {
			return h.checkoutService.UpdateOrderStatus(ctx, order.ID, "confirmed")
		} else {
			return h.checkoutService.UpdateOrderStatus(ctx, order.ID, "failed")
		}
	case "CANCELLATION":
		return h.checkoutService.UpdateOrderStatus(ctx, order.ID, "cancelled")
	case "REFUND":
		return h.checkoutService.UpdateOrderStatus(ctx, order.ID, "refunded")
	}

	return nil
}

func (h *CheckoutHandler) validateStock(ctx context.Context, items []models.CartItem) error {
	for _, item := range items {
		available, err := h.cartService.CheckStockAvailability(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
		if !available {
			return fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}
	}
	return nil
}

func (h *CheckoutHandler) processOrderInventory(ctx context.Context, order models.Order) {
	for _, item := range order.Items {
		if err := h.checkoutService.UpdateInventory(ctx, item.ProductID, -item.Quantity); err != nil {
			log.Printf("Error updating inventory for product %s: %v", item.ProductID, err)
			// Consider implementing a retry mechanism or alerting system
		}
	}
}

func (h *CheckoutHandler) convertToLineItems(items []models.CartItem) []adyen.LineItem {
	lineItems := make([]adyen.LineItem, len(items))
	for i, item := range items {
		lineItems[i] = adyen.LineItem{
			Quantity:           uint32(item.Quantity),
			AmountIncludingTax: uint64(item.Price * 100),
			Description:        item.Name,
			ID:                 item.ProductID,
		}
	}
	return lineItems
}

func (h *CheckoutHandler) calculateTotals(items []models.CartItem, country string) (float64, float64) {
	var subtotal float64
	for _, item := range items {
		subtotal += item.Price * float64(item.Quantity)
	}

	// Calculate tax based on country
	taxRate := h.getTaxRate(country)
	tax := subtotal * taxRate

	return subtotal + tax, tax
}

func (h *CheckoutHandler) getTaxRate(country string) float64 {
	// Implement tax rate logic based on country
	// This is a simplified example
	switch country {
	case "US":
		return 0.0825 // 8.25%
	case "GB":
		return 0.20 // 20% VAT
	default:
		return 0.0
	}
}

func generateReference() string {
	return fmt.Sprintf("ORDER_%s", uuid.New().String())
}
