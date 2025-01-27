// next.config.ts
import path from 'path';
import { NextConfig } from 'next';

const nextConfig: NextConfig = {
  // Keep your existing configuration
  reactStrictMode: true,
  
  // Add webpack configuration for path aliases
  webpack: (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      '@': path.join(__dirname, 'src')
    };
    return config;
  },

  // Add rewrite rules to handle API requests
  async rewrites() {
    return [
      {
        // Forward all API requests to your backend server
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      }
    ];
  },

  // Handle trailing slashes consistently
  trailingSlash: false,

  // Enable strict mode for better development experience
  typescript: {
    // Enable type checking in production builds
    ignoreBuildErrors: false,
  },

  // Additional configuration for handling images and other static assets
  images: {
    domains: ['localhost'],
    // This allows images to be served from your backend
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '8080',
        pathname: '/api/products/**',
      },
    ],
  }
};

export default nextConfig;