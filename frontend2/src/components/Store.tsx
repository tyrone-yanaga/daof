import React, { useEffect, useState } from 'react';
import { useCart } from './CartContext';
import { api } from './ui/ts/types/api';
import { useRouter } from 'next/navigation';

type Product = {
  id: number;
  odoo_id: string;
  name: string;
  description: string;
  base_price: number;
  price?: number;
  sku: string;
  active: boolean;
  created_at: string;
  updated_at: string;
};

const Store: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [sortOrder, setSortOrder] = useState<'high-to-low' | 'low-to-high'>('high-to-low');
  const { addItem } = useCart();
  const router = useRouter();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);
  
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const data = await api.getProducts();
        console.log('Fetched products:', data); // Debug log
        setProducts(data);
      } catch (error) {
        console.error('Failed to fetch products:', error);
      }
    };

    fetchProducts();
  }, []);

  if (!mounted) {
    return null;
  }

  const handleProductClick = async (productId: number) => {
    try {
      router.push(`/product/${productId}`);
    } catch (error) {
      console.error('Failed to fetch related product:', error);
    }
  };

  const getImageUrl = (productId: string) => {
    return `http://localhost:8080/api/products/${productId}/image`;
  };

  const sortedProducts = [...products].sort((a, b) => {
    const priceA = a.price ?? a.base_price;
    const priceB = b.price ?? b.base_price;
    return sortOrder === 'high-to-low' ? priceB - priceA : priceA - priceB;
  });

  const handleAddToCart = (e: React.MouseEvent, productId: number) => {
    e.stopPropagation();
    addItem(productId.toString(), 1);
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
        {sortedProducts.map((product) => {
          console.log('Rendering product:', product.odoo_id); // Debug log per product
          return (
            <div
              key={product.odoo_id}
              className="group relative cursor-pointer"
              onClick={() => {
                console.log('Clicked product:', product.odoo_id); // Debug log on click
                handleProductClick((product.odoo_id as unknown) as number);
              }}
            >
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
                onClick={(e) => handleAddToCart(e, (product.odoo_id as unknown) as number)}
                className="mt-2 w-full bg-black text-white py-2 hover:bg-gray-800 transition-colors"
              >
                Add to Cart
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Store;