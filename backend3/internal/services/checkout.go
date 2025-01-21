// internal/services/checkout.go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/adyen/adyen-go-api-library/v4/adyen"
    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "your-project/internal/config"
    "your-project/internal/models"
)

type CheckoutService struct {
    db          *gorm.DB
    redis       *redis.Client
    cartService *CartService
    adyenClient *adyen.Client
    odooClient  *OdooClient
    config      *config.Config
}

type CheckoutSession struct {
    CartID          uint            `json:"cartId"`
    SessionData     string          `json:"sessionData"`
    PaymentSession  *adyen.Session  `json:"paymentSession"`
    Amount          float64         `json:"amount"`
    Currency        string          `json:"currency"`
    ShippingAddress models.Address  `json:"shippingAddress"`
    BillingAddress  models.Address  `json:"billingAddress"`
    CustomerEmail   string          `json:"customerEmail"`
    Status          string          `json:"status"`
    ExpiresAt       time.Time       `json:"expiresAt"`
}

func NewCheckoutService(
    db *gorm.DB,
    redis *redis.Client,
    cartService *CartService,
    odooClient *OdooClient,
    config *config.Config,
) *CheckoutService {
    adyenClient := adyen.NewClient(&adyen.Config{
        ApiKey:      config.Adyen.APIKey,
        Environment: config.Adyen.Environment,
    })

    return &CheckoutService{
        db:          db,
        redis:       redis,
        cartService: cartService,
        adyenClient: adyenClient,
        odooClient:  odooClient,
        config:      config,
    }
}

// InitiateCheckout starts the checkout process
func (s *CheckoutService) InitiateCheckout(ctx context.Context, cartID uint, req models.CheckoutRequest) (*CheckoutSession, error) {
    // Validate cart and items
    cart, err := s.cartService.Get(ctx, cartID)
    if err != nil {
        return nil, fmt.Errorf("failed to get cart: %w", err)
    }

    if len(cart.Items) == 0 {
        return nil, fmt.Errorf("cart is empty")
    }

    // Validate stock availability
    if err := s.validateStock(ctx, cart.Items); err != nil {
        return nil, err
    }

    // Calculate totals
    subtotal, tax, total := s.calculateTotals(cart.Items, req.ShippingAddress.Country)

    // Create Adyen payment session
    reference := uuid.New().String()
    paymentSession, err := s.createPaymentSession(ctx, total, cart.Items, req.ShippingAddress, reference)
    if err != nil {
        return nil, fmt.Errorf("failed to create payment session: %w", err)
    }

    // Create checkout session
    session := &CheckoutSession{
        CartID:          cartID,
        SessionData:     paymentSession.SessionData,
        PaymentSession:  paymentSession,
        Amount:          total,
        Currency:        "USD",
        ShippingAddress: req.ShippingAddress,
        BillingAddress:  req.BillingAddress,
        CustomerEmail:   req.Email,
        Status:          "pending",
        ExpiresAt:       time.Now().Add(30 * time.Minute),
    }

    // Store session in Redis
    if err := s.storeSession(ctx, cartID, session); err != nil {
        return nil, fmt.Errorf("failed to store session: %w", err)
    }

    return session, nil
}

// CompleteCheckout finalizes the checkout process
func (s *CheckoutService) CompleteCheckout(ctx context.Context, cartID uint, paymentData string) (*models.Order, error) {
    // Get checkout session
    session, err := s.getSession(ctx, cartID)
    if err != nil {
        return nil, fmt.Errorf("failed to get session: %w", err)
    }

    // Verify payment with Adyen
    paymentResult, err := s.verifyPayment(ctx, paymentData, session)
    if err != nil {
        return nil, fmt.Errorf("payment verification failed: %w", err)
    }

    // Start database transaction
    var order *models.Order
    err = s.db.Transaction(func(tx *gorm.DB) error {
        // Create order
        order, err = s.createOrder(ctx, tx, session, paymentResult)
        if err != nil {
            return err
        }

        // Update inventory
        if err := s.updateInventory(ctx, tx, order.Items); err != nil {
            return err
        }

        // Clear cart
        if err := tx.Delete(&models.Cart{ID: cartID}).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("transaction failed: %w", err)
    }

    // Remove session from Redis
    s.redis.Del(ctx, fmt.Sprintf("checkout:%d", cartID))

    return order, nil
}

// HandleWebhook processes Adyen webhook notifications
func (s *CheckoutService) HandleWebhook(ctx context.Context, notification models.AdyenNotification) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        order, err := s.findOrderByPaymentReference(tx, notification.PspReference)
        if err != nil {
            return err
        }

        switch notification.EventCode {
        case "AUTHORISATION":
            if notification.Success {
                order.Status = "confirmed"
                order.PaymentStatus = "paid"
            } else {
                order.Status = "failed"
                order.PaymentStatus = "failed"
            }
        case "CANCELLATION":
            order.Status = "cancelled"
            order.PaymentStatus = "cancelled"
        case "REFUND":
            order.Status = "refunded"
            order.PaymentStatus = "refunded"
        }

        return tx.Save(order).Error
    })
}

// Helper functions

func (s *CheckoutService) createPaymentSession(ctx context.Context, amount float64, items []models.CartItem, address models.Address, reference string) (*adyen.Session, error) {
    req := &adyen.CheckoutSessionRequest{
        Amount: adyen.Amount{
            Currency: "USD",
            Value:    uint64(amount * 100), // Convert to cents
        },
        CountryCode:      address.Country,
        MerchantAccount:  s.config.Adyen.MerchantAccount,
        Reference:        reference,
        ReturnUrl:        s.config.Server.FrontendURL + "/checkout/complete",
        ExpiresAt:        time.Now().Add(30 * time.Minute),
        AllowedPaymentMethods: []string{"scheme", "ideal", "paypal"},
        LineItems:        s.convertToLineItems(items),
    }

    return s.adyenClient.Checkout.Sessions(ctx, req)
}

func (s *CheckoutService) verifyPayment(ctx context.Context, paymentData string, session *CheckoutSession) (*adyen.PaymentVerificationResponse, error) {
    req := &adyen.PaymentVerificationRequest{
        Amount: adyen.Amount{
            Currency: session.Currency,
            Value:    uint64(session.Amount * 100),
        },
        PaymentData: paymentData,
        Reference:   session.PaymentSession.Reference,
    }

    return s.adyenClient.Checkout.PaymentsVerification(ctx, req)
}

func (s *CheckoutService) createOrder(ctx context.Context, tx *gorm.DB, session *CheckoutSession, payment *adyen.PaymentVerificationResponse) (*models.Order, error) {
    order := &models.Order{
        OrderNumber:      fmt.Sprintf("ORD-%s", uuid.New().String()[:8]),
        CartID:          session.CartID,
        CustomerEmail:   session.CustomerEmail,
        ShippingAddress: session.ShippingAddress,
        BillingAddress:  session.BillingAddress,
        Total:          session.Amount,
        Currency:       session.Currency,
        Status:         "processing",
        PaymentStatus:  "pending",
        PaymentID:      payment.PspReference,
    }

    if err := tx.Create(order).Error; err != nil {
        return nil, err
    }

    return order, nil
}

func (s *CheckoutService) validateStock(ctx context.Context, items []models.CartItem) error {
    for _, item := range items {
        available, err := s.odooClient.CheckStockAvailability(ctx, item.ProductID, item.Quantity)
        if err != nil {
            return err
        }
        if !available {
            return fmt.Errorf("insufficient stock for product %d", item.ProductID)
        }
    }
    return nil
}

func (s *CheckoutService) updateInventory(ctx context.Context, tx *gorm.DB, items []models.OrderItem) error {
    for _, item := range items {
        if err := s.odooClient.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
            return err
        }
    }
    return nil
}

func (s *CheckoutService) storeSession(ctx context.Context, cartID uint, session *CheckoutSession) error {
    sessionJSON, err := json.Marshal(session)
    if err != nil {
        return err
    }

    key := fmt.Sprintf("checkout:%d", cartID)
    return s.redis.Set(ctx, key, sessionJSON, 30*time.Minute).Err()
}

func (s *CheckoutService) getSession(ctx context.Context, cartID uint) (*CheckoutSession, error) {
    key := fmt.Sprintf("checkout:%d", cartID)
    data, err := s.redis.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    var session CheckoutSession
    if err := json.Unmarshal([]byte(data), &session); err != nil {
        return nil, err
    }

    return &session, nil
}

func (s *CheckoutService) calculateTotals(items []models.CartItem, country string) (subtotal, tax, total float64) {
    for _, item := range items {
        subtotal += item.Price * float64(item.Quantity)
    }

    taxRate := s.getTaxRate(country)
    tax = subtotal * taxRate
    total = subtotal + tax

    return
}

func (s *CheckoutService) getTaxRate(country string) float64 {
    // Simplified tax logic - in production, use a tax service
    switch country {
    case "US":
        return 0.0825 // 8.25%
    case "GB":
        return 0.20 // 20% VAT
    default:
        return 0.0
    }
}

func (s *CheckoutService) convertToLineItems(items []models.CartItem) []adyen.LineItem {
    lineItems := make([]adyen.LineItem, len(items))
    for i, item := range items {
        lineItems[i] = adyen.LineItem{
            Quantity:           uint32(item.Quantity),
            AmountIncludingTax: uint64(item.Price * 100),
            Description:        item.Name,
            ID:                fmt.Sprintf("%d", item.ProductID),
        }
    }
    return lineItems
}