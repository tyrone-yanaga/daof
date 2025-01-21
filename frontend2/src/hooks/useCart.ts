// src/hooks/useCart.ts 
import { create } from 'zustand';

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

interface CartState {
  items: CartItem[];
  isLoading: boolean;
  error: string | null;
  addItem: (productId: string, quantity: number) => Promise<void>;
  updateItem: (itemId: string, quantity: number) => Promise<void>;
  removeItem: (itemId: string) => Promise<void>;
  clearCart: () => void;
}

export const useCart = create<CartState>((set, get) => ({
  items: [],
  isLoading: false,
  error: null,

  addItem: async (productId: string, quantity: number) => {
    set({ isLoading: true });
    try {
      // Fetch product details
      const response = await fetch(`/api/products/${productId}`);
      if (!response.ok) throw new Error('Failed to fetch product');
      const product: Product = await response.json();

      // Update state
      const currentItems = get().items;
      const existingItem = currentItems.find(item => item.productId === productId);

      if (existingItem) {
        // Update existing item
        const updatedItems = currentItems.map(item =>
          item.productId === productId
            ? { ...item, quantity: item.quantity + quantity }
            : item
        );
        set({ items: updatedItems, isLoading: false });
      } else {
        // Add new item
        const newItem: CartItem = {
          id: Date.now().toString(),
          productId,
          quantity,
          product
        };
        set({ items: [...currentItems, newItem], isLoading: false });
      }
    } catch (error) {
      set({ error: 'Failed to add item to cart', isLoading: false });
    }
  },

  updateItem: async (itemId: string, quantity: number) => {
    if (quantity === 0) {
      get().removeItem(itemId);
      return;
    }

    set(state => ({
      items: state.items.map(item =>
        item.id === itemId ? { ...item, quantity } : item
      )
    }));
  },

  removeItem: async (itemId: string) => {
    set(state => ({
      items: state.items.filter(item => item.id !== itemId)
    }));
  },

  clearCart: () => {
    set({ items: [], error: null });
  }
}));

// Helper function to calculate cart totals
export const useCartTotals = () => {
  const items = useCart(state => state.items);
  
  const subtotal = items.reduce(
    (sum, item) => sum + item.product.price * item.quantity, 
    0
  );
  const tax = subtotal * 0.1; // 10% tax
  const total = subtotal + tax;

  return { subtotal, tax, total };
};