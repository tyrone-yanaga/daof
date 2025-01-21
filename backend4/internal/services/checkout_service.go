package services

import (
	"context"
	"ecommerce/internal/models"
	"ecommerce/pkg/odoo"
	"ecommerce/pkg/queue"
	"ecommerce/pkg/redis"
	"fmt"
	"time"

	"github.com/adyen/adyen-go-api-library/v5/src/checkout"
	"github.com/google/uuid"
)

type CheckoutService struct {
	cartService  *CartService
	orderService *OrderService
	redisClient  *redis.Client
	odooClient   *odoo.Client
	queueClient  *queue.Client
	adyenClient  *checkout.Client
	baseURL      string
}

func NewCheckoutService(
	cartService *CartService,
	orderService *OrderService,
	redisClient *redis.Client,
	odooClient *odoo.Client,
	queueClient *queue.Client,
	adyenClient *checkout.Client,
	baseURL string,
) *CheckoutService {
	return &CheckoutService{
		cartService:  cartService,
		orderService: orderService,
		redisClient:  redisClient,
		odooClient:   odooClient,
		queueClient:  queueClient,
		adyenClient:  adyenClient,
		baseURL:      baseURL,
	}
}

func (s *CheckoutService) InitiateCheckout(ctx context.Context, req *models.CheckoutRequest) (*models.CheckoutSession, error) {
	// Get cart
	cart, err := s.cartService.GetCart(ctx, req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Validate cart is not empty
	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Create checkout session
	session := &models.CheckoutSession{
		ID:           uuid.New().String(),
		CartID:       cart.ID,
		UserID:       cart.UserID,
		Status:       "pending",
		Total:        cart.Total,
		Currency:     req.Currency,
		ShippingInfo: req.ShippingInfo,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(30 * time.Minute),
	}

	// Save checkout session
	err = s.redisClient.Set(ctx, fmt.Sprintf("checkout:%s", session.ID), session, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to save checkout session: %w", err)
	}

	return session, nil
}

func (s *CheckoutService) CreatePaymentSession(ctx context.Context, checkoutID string) (*checkout.PaymentLinksResponse, error) {
	// Get checkout session
	var session models.CheckoutSession
	err := s.redisClient.Get(ctx, fmt.Sprintf("checkout:%s", checkoutID), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkout session: %w", err)
	}

	// Create Adyen payment session
	amount := checkout.Amount{
		Currency: session.Currency,
		Value:    int64(session.Total * 100), // Convert to cents
	}

	req := &checkout.PaymentLinkRequest{
		Reference:   checkoutID,
		Amount:      amount,
		Description: fmt.Sprintf("Order %s", checkoutID),
		ReturnUrl:   fmt.Sprintf("%s/api/checkout/%s/complete", s.baseURL, checkoutID),
	}

	resp, err := s.adyenClient.PaymentLinks.PaymentLinks(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment session: %w", err)
	}

	// Update checkout session with payment ID
	session.PaymentID = resp.PaymentLinkID
	session.Status = "processing"

	err = s.redisClient.Set(ctx, fmt.Sprintf("checkout:%s", session.ID), session, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to update checkout session: %w", err)
	}

	return resp, nil
}

func (s *CheckoutService) CompleteCheckout(ctx context.Context, checkoutID string, paymentData map[string]interface{}) error {
	// Get checkout session
	var session models.CheckoutSession
	err := s.redisClient.Get(ctx, fmt.Sprintf("checkout:%s", checkoutID), &session)
	if err != nil {
		return fmt.Errorf("failed to get checkout session: %w", err)
	}

	// Get cart
	cart, err := s.cartService.GetCart(ctx, session.CartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Create order
	order := &models.Order{
		UserID:       session.UserID,
		Status:       "pending",
		Total:        session.Total,
		PaymentID:    session.PaymentID,
		ShippingInfo: session.ShippingInfo,
		Items:        make([]models.OrderItem, len(cart.Items)),
	}

	// Convert cart items to order items
	for i, item := range cart.Items {
		order.Items[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// Create order in database and Odoo
	order, err = s.orderService.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Publish order created event
	err = s.queueClient.Publish("orders", queue.Message{
		Type:    "order.created",
		Payload: order,
	})
	if err != nil {
		// Log error but don't fail the checkout
		fmt.Printf("failed to publish order created event: %v\n", err)
	}

	// Clean up cart and checkout session
	s.redisClient.Delete(ctx, fmt.Sprintf("cart:%s", session.CartID))
	s.redisClient.Delete(ctx, fmt.Sprintf("checkout:%s", session.ID))

	return nil
}
