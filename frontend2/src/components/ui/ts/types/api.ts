// api.ts
const BASE_URL = 'http://localhost:8080';


export const api = {
    async getProducts(): Promise<Product[]> {
      const response = await fetch(`${BASE_URL}/api/products`);
      return response.json();
    },
  
    async getProduct(id: string): Promise<Product> {
      const response = await fetch(`${BASE_URL}/api/products/${id}`);
      return response.json();
    },

    //TODO promise type check - is []byte converted to string
    async getProductImage(id: string): Promise<string> {
      const response = await fetch(`${BASE_URL}/api/products/${id}/image`);
      return response.json();
    },
  
    async createCart(): Promise<{ id: string }> {
      const response = await fetch(`${BASE_URL}/api/carts`, {
        method: 'POST'
      });
      return response.json();
    },
  
    async getCart(id: string) {
      const response = await fetch(`${BASE_URL}/api/carts/${id}`);
      return response.json();
    },
  
    async addToCart(cartId: string, productId: string, quantity: number) {
      const response = await fetch(`${BASE_URL}/api/carts/${cartId}/items`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ productId, quantity })
      });
      return response.json();
    },
  
    async updateCartItem(cartId: string, itemId: string, quantity: number) {
      const response = await fetch(`${BASE_URL}/api/carts/${cartId}/items`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ itemId, quantity })
      });
      return response.json();
    },
  
    async initiateCheckout(cartId: string) {
      const response = await fetch(`${BASE_URL}/api/checkout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ cartId })
      });
      return response.json();
    }
  };
  