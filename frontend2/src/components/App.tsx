// components/App.tsx
import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Store from './Store';
import ProductDetail from './ProductDetail';
import { CartProvider } from './CartContext';

const App: React.FC = () => {
  return (
    <CartProvider>
      <BrowserRouter>
        <div className="min-h-screen">
          <Routes>
            <Route path="/" element={<Store />} />
            <Route path="/product/:id" element={<ProductDetail />} />
          </Routes>
        </div>
      </BrowserRouter>
    </CartProvider>
  );
};

export default App;