package handlers

type Handlers struct {
	Product  ProductHandler
	Order    OrderHandler
	Checkout CheckoutHandler
	Cart     CartHandler
}
