"use client";

import { Suspense } from 'react';
import { CartProvider } from '@/components/CartContext';
import { Header } from '@/components/Header';
import { Store } from '@/components/Store';
import Loading from './loading';

export default function StorePage() {
  return (
    <CartProvider>
      <div className="min-h-screen flex flex-col">
        <Header />
        <main className="flex-1">
          <Suspense fallback={<Loading />}>
            <Store />
          </Suspense>
        </main>
      </div>
    </CartProvider>
  );
}