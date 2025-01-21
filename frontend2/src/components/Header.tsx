
// Header.tsx
import React, { useState } from 'react';
import { useCart } from './CartContext';
import { ShoppingCart } from 'lucide-react';

export const Header = () => {
  const [isCartOpen, setIsCartOpen] = useState(false);
  const { cart } = useCart();

  const itemCount = cart?.items.reduce((sum, item) => sum + item.quantity, 0) ?? 0;

  return (
    <header className="bg-black text-white">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <nav className="flex-1 flex justify-center space-x-8">
            <a href="/videos" className="hover:text-gray-300">VIDEOS</a>
            <a href="/team" className="hover:text-gray-300">TEAM</a>
            <a href="/store" className="hover:text-gray-300">STORE</a>
          </nav>
          
          <button
            onClick={() => setIsCartOpen(true)}
            className="flex items-center space-x-2 hover:text-gray-300"
          >
            <span>CART</span>
            <ShoppingCart className="h-5 w-5" />
            {itemCount > 0 && (
              <span className="bg-white text-black rounded-full w-5 h-5 flex items-center justify-center text-sm">
                {itemCount}
              </span>
            )}
          </button>
        </div>
      </div>

      {/* Cart Sidebar */}
      {isCartOpen && (
        <Cart onClose={() => setIsCartOpen(false)} />
      )}
    </header>
  );
};