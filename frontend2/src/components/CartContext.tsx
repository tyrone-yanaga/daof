
// CartContext.tsx
import React, { createContext, useContext, useState, useEffect } from 'react';

interface CartContextType {
  cart: Cart | null;
  isLoading: boolean;
  addItem: (productId: string, quantity: number) => Promise<void>;
  updateItem: (itemId: string, quantity: number) => Promise<void>;
  initiateCheckout: () => Promise<void>;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [cart, setCart] = useState<Cart | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initCart = async () => {
      try {
        // Try to get cart ID from localStorage
        let cartId = localStorage.getItem('cartId');
        
        if (!cartId) {
          // Create new cart if none exists
          const newCart = await api.createCart();
          cartId = newCart.id;
          localStorage.setItem('cartId', cartId);
        }

        const cartData = await api.getCart(cartId);
        setCart(cartData);
      } catch (error) {
        console.error('Failed to initialize cart:', error);
      } finally {
        setIsLoading(false);
      }
    };

    initCart();
  }, []);

  const addItem = async (productId: string, quantity: number) => {
    if (!cart) return;
    
    try {
      const updatedCart = await api.addToCart(cart.id, productId, quantity);
      setCart(updatedCart);
    } catch (error) {
      console.error('Failed to add item to cart:', error);
    }
  };

  const updateItem = async (itemId: string, quantity: number) => {
    if (!cart) return;

    try {
      const updatedCart = await api.updateCartItem(cart.id, itemId, quantity);
      setCart(updatedCart);
    } catch (error) {
      console.error('Failed to update cart item:', error);
    }
  };

  const initiateCheckout = async () => {
    if (!cart) return;

    try {
      const checkoutSession = await api.initiateCheckout(cart.id);
      // Redirect to checkout page or handle checkout flow
      window.location.href = checkoutSession.url;
    } catch (error) {
      console.error('Failed to initiate checkout:', error);
    }
  };

  return (
    <CartContext.Provider value={{ cart, isLoading, addItem, updateItem, initiateCheckout }}>
      {children}
    </CartContext.Provider>
  );
};

export const useCart = () => {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};