"use client";

import { Suspense } from 'react';
import { CartProvider } from '@/components/CartContext';
import { Header } from '@/components/Header';
import ProductDetail from '@/components/ProductDetail';
import Loading from './loading';

const product = {
    id: 1,
    odoo_id: 123,
    name: "Wheel Talk - Calligraphy Long Sleeve (Black)",
    description: "Dual sleeve action for wind / weather protection...",
    base_price: 45.00,
    sku: "WT-CLS-BLK",
    active: true,
    created_at: "2021-10-01T00:00:00Z",
    updated_at: "2021-10-01T00:00:00Z",
    variants: [
        { id: 1, name: "Black", sku: "WT-CLS-BLK", active: true, created_at: "2021-10-01T00:00:00Z", updated_at: "2021-10-01T00:00:00Z", price: 45.00 },
        { id: 2, name: "White", sku: "WT-CLS-WHT", active: true, created_at: "2021-10-01T00:00:00Z", updated_at: "2021-10-01T00:00:00Z", price: 45.00 }
    ],
    attributes: [
        { id: 1, name: "Size", values: ["S", "M", "L", "XL"] },
        { id: 2, name: "Color", values: ["Black", "White"] }
    ],
    image_1920: "/images/products/wheel-talk-calligraphy-long-sleeve-black-1920.jpg",
    image_1024: "/images/products/wheel-talk-calligraphy-long-sleeve-black-1024.jpg",
    image_128: "/images/products/wheel-talk-calligraphy-long-sleeve-black-128.jpg"
};


export default function ProductPage() {
  return (
    <CartProvider>
      <div className="min-h-screen flex flex-col">
        <Header />
        <main className="flex-1">
          <Suspense fallback={<Loading />}>
            <ProductDetail product={product}/>
          </Suspense>
        </main>
      </div>
    </CartProvider>
  );
}