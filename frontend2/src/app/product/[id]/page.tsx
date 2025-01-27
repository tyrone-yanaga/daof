// src/app/product/[id]/page.tsx
'use client';

import React from 'react';
import ProductDetail from '@/components/ProductDetail';
import { useParams } from 'next/navigation';

export default function ProductPage() {
  const params = useParams();
  return <ProductDetail id={params.id as string} />;
}