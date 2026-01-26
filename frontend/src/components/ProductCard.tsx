import { Star, ShoppingCart } from 'lucide-react';
import type { Product } from '../graphql/types';

export function ProductCard({ product }: { product: Product }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow">
      <div className="aspect-square bg-gray-100 overflow-hidden">
        <img
          src={product.imageUrl}
          alt={product.name}
          loading="lazy"
          className="w-full h-full object-cover hover:scale-105 transition-transform"
        />
      </div>
      <div className="p-4">
        <span className="text-xs font-medium text-primary-600 uppercase">
          {product.category.name}
        </span>
        <h3 className="mt-1 text-lg font-semibold text-gray-900 line-clamp-1">
          {product.name}
        </h3>
        <p className="mt-1 text-sm text-gray-500 line-clamp-2">
          {product.description}
        </p>
        <div className="mt-3 flex items-center gap-1">
          {[1, 2, 3, 4, 5].map((star) => (
            <Star
              key={star}
              className={`w-4 h-4 ${
                star <= Math.round(product.rating)
                  ? 'text-amber-400 fill-amber-400'
                  : 'text-gray-300'
              }`}
            />
          ))}
          <span className="text-sm text-gray-600 ml-1">{product.rating.toFixed(1)}</span>
          <span className="text-sm text-gray-400">({product.reviewCount})</span>
        </div>
        <div className="mt-4 flex items-center justify-between">
          <span className="text-xl font-bold text-gray-900">
            ${product.price.toFixed(2)}
          </span>
          <button className="p-2 rounded-lg bg-primary-600 text-white hover:bg-primary-700">
            <ShoppingCart className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}

export function ProductCardSkeleton() {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden animate-pulse">
      <div className="aspect-square bg-gray-200" />
      <div className="p-4 space-y-3">
        <div className="h-3 bg-gray-200 rounded w-1/4" />
        <div className="h-5 bg-gray-200 rounded w-3/4" />
        <div className="h-4 bg-gray-200 rounded w-full" />
        <div className="h-6 bg-gray-200 rounded w-1/3 mt-4" />
      </div>
    </div>
  );
}
