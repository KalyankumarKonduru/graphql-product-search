import { Star, ShoppingCart } from 'lucide-react';
import type { Product } from '../graphql/types';

interface ProductCardProps {
  product: Product;
  onClick?: () => void;
}

export function ProductCard({ product, onClick }: ProductCardProps) {
  const formatPrice = (price: number) =>
    new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(price);

  return (
    <article
      onClick={onClick}
      className="group bg-white rounded-xl border border-gray-200 overflow-hidden hover:shadow-lg transition-all duration-300 cursor-pointer"
    >
      {/* Image */}
      <div className="relative aspect-square overflow-hidden bg-gray-100">
        <img
          src={product.imageUrl}
          alt={product.name}
          loading="lazy"
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
        />
        {product.stock <= 5 && product.stock > 0 && (
          <span className="absolute top-3 left-3 bg-amber-500 text-white text-xs font-medium px-2 py-1 rounded-full">
            Only {product.stock} left
          </span>
        )}
        {product.stock === 0 && (
          <span className="absolute top-3 left-3 bg-red-500 text-white text-xs font-medium px-2 py-1 rounded-full">
            Out of Stock
          </span>
        )}
      </div>

      {/* Content */}
      <div className="p-4">
        {/* Category */}
        <span className="text-xs font-medium text-primary-600 uppercase tracking-wide">
          {product.category.name}
        </span>

        {/* Name */}
        <h3 className="mt-1 text-lg font-semibold text-gray-900 line-clamp-1 group-hover:text-primary-600 transition-colors">
          {product.name}
        </h3>

        {/* Description */}
        <p className="mt-1 text-sm text-gray-500 line-clamp-2">
          {product.description}
        </p>

        {/* Rating */}
        <div className="mt-3 flex items-center gap-1">
          <div className="flex items-center">
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
          </div>
          <span className="text-sm text-gray-600 ml-1">
            {product.rating.toFixed(1)}
          </span>
          <span className="text-sm text-gray-400">
            ({product.reviewCount})
          </span>
        </div>

        {/* Price and Action */}
        <div className="mt-4 flex items-center justify-between">
          <span className="text-xl font-bold text-gray-900">
            {formatPrice(product.price)}
          </span>
          <button
            onClick={(e) => {
              e.stopPropagation();
              // Add to cart logic
            }}
            disabled={product.stock === 0}
            className="p-2 rounded-lg bg-primary-600 text-white hover:bg-primary-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
            aria-label={`Add ${product.name} to cart`}
          >
            <ShoppingCart className="w-5 h-5" />
          </button>
        </div>
      </div>
    </article>
  );
}

// Loading skeleton
export function ProductCardSkeleton() {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden animate-pulse">
      <div className="aspect-square bg-gray-200" />
      <div className="p-4 space-y-3">
        <div className="h-3 bg-gray-200 rounded w-1/4" />
        <div className="h-5 bg-gray-200 rounded w-3/4" />
        <div className="h-4 bg-gray-200 rounded w-full" />
        <div className="h-4 bg-gray-200 rounded w-2/3" />
        <div className="flex justify-between items-center mt-4">
          <div className="h-6 bg-gray-200 rounded w-1/4" />
          <div className="h-10 w-10 bg-gray-200 rounded-lg" />
        </div>
      </div>
    </div>
  );
}
