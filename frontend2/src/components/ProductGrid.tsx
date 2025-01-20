"use client";

import React, { useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';

interface Product {
  id: string;
  name: string;
  description: string;
  thumbnail: string;
  price: number;
}

const ProductGrid = () => {
  const [hoveredProduct, setHoveredProduct] = useState<string | null>(null);

  // Sample products data
  const products: Product[] = [
    {
      id: '1',
      name: 'Product 1',
      description: 'Short description of product 1',
      thumbnail: '/api/placeholder/300/300',
      price: 99.99
    },
    // Add more products...
  ];

  const handleProductClick = (productId: string) => {
    // Using standard window location for navigation
    window.location.href = `/products/${productId}`;
    
    // Alternatively, if you're using React Router, you could use:
    // navigate(`/products/${productId}`);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((product) => (
          <Card
            key={product.id}
            className="relative cursor-pointer transform transition-transform duration-200 hover:scale-105"
            onMouseEnter={() => setHoveredProduct(product.id)}
            onMouseLeave={() => setHoveredProduct(null)}
            onClick={() => handleProductClick(product.id)}
          >
            <div className="aspect-square overflow-hidden">
              <img
                src={product.thumbnail}
                alt={product.name}
                className="w-full h-full object-cover"
              />
            </div>
            
            {/* Hover overlay */}
            {hoveredProduct === product.id && (
              <div className="absolute inset-0 bg-black bg-opacity-50 flex flex-col justify-end p-4 text-white transition-opacity duration-200">
                <h3 className="text-lg font-semibold mb-2">{product.name}</h3>
                <p className="text-sm">{product.description}</p>
                <p className="text-lg font-bold mt-2">
                  ${product.price.toFixed(2)}
                </p>
              </div>
            )}
          </Card>
        ))}
      </div>
    </div>
  );
};

export default ProductGrid;