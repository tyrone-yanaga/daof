// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  // Get the pathname from the request URL
  const pathname = request.nextUrl.pathname;

  // Allow client-side routes to be handled by Next.js
  if (pathname.startsWith('/product/')) {
    return NextResponse.next();
  }

  // Allow API requests to be handled by the rewrite rules
  if (pathname.startsWith('/api/')) {
    return NextResponse.next();
  }

  // For all other routes, proceed normally
  return NextResponse.next();
}

// Configure which paths should be handled by this middleware
export const config = {
  matcher: [
    // Match all routes that start with /product/
    '/product/:path*',
    // Match all API routes
    '/api/:path*',
  ],
};