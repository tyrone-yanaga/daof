
// Store.tsx
import React, { useEffect, useState } from 'react';
import { useCart } from './CartContext';
import { api } from './ui/ts/types/api'; 

const Store = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [sortOrder, setSortOrder] = useState<'high-to-low' | 'low-to-high'>('high-to-low');
  const { addItem } = useCart();

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const data = await api.getProducts();
        setProducts(data);
      } catch (error) {
        console.error('Failed to fetch products:', error);
      }
    };

    fetchProducts();
  }, []);

  const sortedProducts = [...products].sort((a, b) => {
    return sortOrder === 'high-to-low' ? b.price - a.price : a.price - b.price;
  });

  const getImageUrl = (productId: string) => {
    return `http://localhost:8080/api/products/${productId}/image`;
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">STORE</h1>
        <select
          className="p-2 border rounded"
          value={sortOrder}
          onChange={(e) => setSortOrder(e.target.value as 'high-to-low' | 'low-to-high')}
        >
          <option value="high-to-low">Sort by price: high to low</option>
          <option value="low-to-high">Sort by price: low to high</option>
        </select>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {sortedProducts.map((product) => (
          <div key={product.id} className="group relative">
            <div className="aspect-square overflow-hidden">
            <img
              src={getImageUrl(product.odoo_id)}
              alt={product.name}
              className="w-full h-full object-cover transform transition-transform group-hover:scale-105"
            />
            </div>
            <div className="mt-4 flex justify-between items-center">
              <h2 className="text-lg font-medium">{product.name}</h2>
              <p className="text-lg">${product.base_price.toFixed(2)}</p>
            </div>
            <button
              onClick={() => addItem(product.id, 1)}
              className="mt-2 w-full bg-black text-white py-2 hover:bg-gray-800 transition-colors"
            >
              Add to Cart
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export { Header, Store, CartProvider };