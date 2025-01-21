// types.ts
interface Product {
    id: string;
    name: string;
    price: number;
    image: string;
    description?: string;
  }
  
  interface CartItem {
    id: string;
    productId: string;
    quantity: number;
    product: Product;
  }
  
  interface Cart {
    id: string;
    items: CartItem[];
  }
  