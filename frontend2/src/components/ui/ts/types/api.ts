// api.ts
const api = {
    async getProducts() {
      const response = await fetch('/api/products');
      return response.json();
    },
  
    async getProduct(id: string) {
      const response = await fetch(`/api/products/${id}`);
      return response.json();
    },
  
    async createCart() {
      const response = await fetch('/api/carts', {
        method: 'POST'
      });
      return response.json();
    },
  
    async getCart(id: string) {
      const response = await fetch(`/api/carts/${id}`);
      return response.json();
    },
  
    async addToCart(cartId: string, productId: string, quantity: number) {
      const response = await fetch(`/api/carts/${cartId}/items`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ productId, quantity })
      });
      return response.json();
    },
  
    async updateCartItem(cartId: string, itemId: string, quantity: number) {
      const response = await fetch(`/api/carts/${cartId}/items`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ itemId, quantity })
      });
      return response.json();
    },
  
    async initiateCheckout(cartId: string) {
      const response = await fetch('/api/checkout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ cartId })
      });
      return response.json();
    }
  };
  