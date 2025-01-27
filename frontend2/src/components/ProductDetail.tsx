import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { BASE_URL } from './ui/ts/types/api';

// Type definitions
type ProductVariant = {
  // Add variant properties as needed
};

type ProductAttribute = {
  // Add attribute properties as needed
};
interface ProductDetailProps {
  id: string;
}

type Product = {
  id: number;
  odoo_id: number;
  name: string;
  description: string;
  base_price: number;
  sku: string;
  active: boolean;
  created_at: string;
  updated_at: string;
  variants: ProductVariant[];
  attributes: ProductAttribute[];
  image_1920: string;
  image_1024: string;
  image_128: string;
};

const ProductDetail: React.FC<ProductDetailProps> = ({id}) => {
  const router = useRouter();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSize, setSelectedSize] = useState<string>('M');
  const [quantity, setQuantity] = useState<number>(1);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        setLoading(true);
        console.log('---------Fetching product:', id);

        const response = await fetch(`${BASE_URL}/api/products/${id}`);
        if (!response.ok) {
          throw new Error('Product not found');
        }
        const data = await response.json();
        console.log('---------Fetched product:', data);
        setProduct(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch product');
        console.error('Error fetching product:', err);
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchProduct();
    }
  }, [id]);

  if (!mounted) {
    return null;
  }

  const handleBackToStore = () => {
    router.push('/');  // Using Next.js router for navigation
  };

  // Loading state
  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center h-64">
          <div className="text-xl">Loading...</div>
        </div>
      </div>
    );
  }

  // Error state
  if (error || !product) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col items-center space-y-4">
          <div className="text-xl text-red-600">
            {error || 'Product not found'}
          </div>
          <button
            onClick={() => handleBackToStore}
            className="px-4 py-2 bg-black text-white hover:bg-gray-800 transition-colors"
          >
            Back to Store
          </button>
        </div>
      </div>
    );
  }

  // Handle adding to cart
  const handleAddToCart = () => {
    console.log('Adding to cart:', {
      product,
      size: selectedSize,
      quantity
    });
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <button
        onClick={() => handleBackToStore()}
        className="mb-6 text-gray-600 hover:text-gray-800 flex items-center"
      >
        ← Back to Store
      </button>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Product Image Section */}
        <div className="relative overflow-hidden aspect-square">
          <div className="w-full h-full group">
            <img
              src={product.image_1920 || product.image_1024 || product.image_128 || `/api/products/${product.odoo_id}/image`}
              alt={product.name}
              className="w-full h-full object-cover transition-transform duration-500 ease-in-out group-hover:scale-125"
            />
          </div>
        </div>

        {/* Product Details Section */}
        <div className="flex flex-col space-y-6">
          <h1 className="text-3xl font-bold">{product.name}</h1>
          <div className="text-2xl">${product.base_price.toFixed(2)}</div>

          <div className="space-y-4">
            <p className="text-gray-600 italic">{product.description}</p>

            <div className="space-y-2">
              <div>• 100% combed and ring-spun cotton</div>
              <div>• Two sleeves for as many as two arms</div>
              <div>• Lightweight / breathable</div>
              <div>• Available for a limited time</div>
              <div>• 3 Color options (White, Black, Olive)</div>
              <div>• Printed by the order</div>
            </div>
          </div>

          {/* Size Selector */}
          <div className="space-y-2">
            <label className="block text-sm font-medium">Size</label>
            <div className="flex items-center space-x-2">
              <select
                value={selectedSize}
                onChange={(e) => setSelectedSize(e.target.value)}
                className="p-2 border rounded bg-gray-50 w-full md:w-48"
              >
                <option value="S">S</option>
                <option value="M">M</option>
                <option value="L">L</option>
                <option value="XL">XL</option>
              </select>
              <button 
                onClick={() => setSelectedSize('M')}
                className="text-sm text-yellow-600 hover:text-yellow-700"
              >
                Clear
              </button>
            </div>
          </div>

          {/* Quantity and Add to Cart */}
          <div className="flex space-x-4">
            <div className="w-24">
              <input
                type="number"
                min="1"
                value={quantity}
                onChange={(e) => setQuantity(Math.max(1, Number(e.target.value)))}
                className="w-full p-2 border rounded text-center"
              />
            </div>
            <button
              onClick={handleAddToCart}
              className="flex-1 bg-black text-white py-2 px-4 hover:bg-gray-800 transition-colors"
            >
              ADD TO CART
            </button>
          </div>

          {/* Product Details */}
          <div className="pt-6 space-y-2 text-sm text-gray-500">
            <div>SKU: {product.sku}</div>
            <div>Category: Shirts</div>
            <div>Tag: Shirts</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductDetail;