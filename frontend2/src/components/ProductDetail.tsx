"use client";

import React, { useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { ChevronUp, ChevronDown, ShoppingCart } from 'lucide-react';

interface ProductVariation {
  id: string;
  name: string;
  options: string[];
}

interface ProductImage {
  id: string;
  url: string;
  alt: string;
}

interface Product {
  id: string;
  name: string;
  price: number;
  description: string;
  variations: ProductVariation[];
  images: ProductImage[];
  details: string;
  sizeChart: string;
}

const ProductDetail = () => {
  const [selectedImage, setSelectedImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const [selectedVariations, setSelectedVariations] = useState<Record<string, string>>({});

  // Sample product data
  const product: Product = {
    id: '1',
    name: 'Sample Product',
    price: 99.99,
    description: 'Detailed product description goes here...',
    variations: [
      { id: 'size', name: 'Size', options: ['S', 'M', 'L', 'XL'] },
      { id: 'color', name: 'Color', options: ['Red', 'Blue', 'Black'] }
    ],
    images: [
      { id: '1', url: '/api/placeholder/600/600', alt: 'Main product image' },
      { id: '2', url: '/api/placeholder/150/150', alt: 'Product detail 1' },
      { id: '3', url: '/api/placeholder/150/150', alt: 'Product detail 2' }
    ],
    details: 'Additional product details...',
    sizeChart: 'Size chart information...'
  };

  const handleQuantityChange = (increment: boolean) => {
    setQuantity(prev => increment ? prev + 1 : Math.max(1, prev - 1));
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Left column - Images */}
        <div className="space-y-4">
          <div className="relative aspect-square">
            <img
              src={product.images[selectedImage].url}
              alt={product.images[selectedImage].alt}
              className="w-full h-full object-cover rounded-lg"
            />
          </div>
          <div className="grid grid-cols-4 gap-2">
            {product.images.map((image, index) => (
              <button
                key={image.id}
                onClick={() => setSelectedImage(index)}
                className={`aspect-square rounded-md overflow-hidden border-2 ${
                  selectedImage === index ? 'border-blue-500' : 'border-transparent'
                }`}
              >
                <img
                  src={image.url}
                  alt={image.alt}
                  className="w-full h-full object-cover"
                />
              </button>
            ))}
          </div>
        </div>

        {/* Right column - Product details */}
        <div className="space-y-6">
          <h1 className="text-3xl font-bold">{product.name}</h1>
          <p className="text-gray-600">{product.description}</p>
          
          {/* Variations */}
          <div className="space-y-4">
            {product.variations.map(variation => (
              <div key={variation.id}>
                <label className="block text-sm font-medium mb-2">
                  {variation.name}
                </label>
                <Select
                  onValueChange={(value) => 
                    setSelectedVariations(prev => ({...prev, [variation.id]: value}))
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder={`Select ${variation.name}`} />
                  </SelectTrigger>
                  <SelectContent>
                    {variation.options.map(option => (
                      <SelectItem key={option} value={option}>
                        {option}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            ))}
          </div>

          {/* Price and quantity */}
          <div className="space-y-4">
            <p className="text-2xl font-bold">${product.price.toFixed(2)}</p>
            
            <div className="flex items-center space-x-4">
              <span className="font-medium">Quantity:</span>
              <div className="flex items-center border rounded-md">
                <Button 
                  variant="ghost"
                  size="sm"
                  onClick={() => handleQuantityChange(false)}
                >
                  <ChevronDown className="h-4 w-4" />
                </Button>
                <span className="px-4">{quantity}</span>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleQuantityChange(true)}
                >
                  <ChevronUp className="h-4 w-4" />
                </Button>
              </div>
            </div>

            <Button className="w-full flex items-center justify-center gap-2">
              <ShoppingCart className="h-5 w-5" />
              Add to Cart
            </Button>
          </div>
        </div>
      </div>

      {/* Tabs section */}
      <div className="mt-12">
        <Tabs defaultValue="details">
          <TabsList className="w-full justify-start">
            <TabsTrigger value="details">Details</TabsTrigger>
            <TabsTrigger value="size">Size Chart</TabsTrigger>
            <TabsTrigger value="reviews">Reviews</TabsTrigger>
          </TabsList>
          
          <TabsContent value="details">
            <Card>
              <CardContent className="pt-6">
                {product.details}
              </CardContent>
            </Card>
          </TabsContent>
          
          <TabsContent value="size">
            <Card>
              <CardContent className="pt-6">
                {product.sizeChart}
              </CardContent>
            </Card>
          </TabsContent>
          
          <TabsContent value="reviews">
            <Card>
              <CardContent className="pt-6">
                {/* Reviews component would go here */}
                <p>Product reviews...</p>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default ProductDetail;