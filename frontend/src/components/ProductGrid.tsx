import { ProductCard, ProductCardSkeleton } from './ProductCard';
import { useInfiniteScroll } from '../hooks/useProducts';
import { GraphQLError } from './ErrorBoundary';
import { Loader2 } from 'lucide-react';
import type { Product, PageInfo } from '../graphql/types';

interface ProductGridProps {
  products: Product[];
  loading: boolean;
  error?: Error;
  pageInfo?: PageInfo;
  hasMore: boolean;
  onLoadMore: () => void;
  onProductClick?: (product: Product) => void;
  onRetry?: () => void;
}

export function ProductGrid({
  products,
  loading,
  error,
  pageInfo,
  hasMore,
  onLoadMore,
  onProductClick,
  onRetry,
}: ProductGridProps) {
  // Infinite scroll ref
  const loadMoreRef = useInfiniteScroll(onLoadMore, hasMore, loading);

  // Error state
  if (error && products.length === 0) {
    return (
      <div className="py-12">
        <GraphQLError error={error} onRetry={onRetry} />
      </div>
    );
  }

  // Initial loading state
  if (loading && products.length === 0) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {Array.from({ length: 8 }).map((_, i) => (
          <ProductCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  // Empty state
  if (!loading && products.length === 0) {
    return (
      <div className="py-16 text-center">
        <div className="w-24 h-24 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
          <svg
            className="w-12 h-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">
          No products found
        </h3>
        <p className="text-gray-500 max-w-sm mx-auto">
          Try adjusting your search or filters to find what you're looking for.
        </p>
      </div>
    );
  }

  return (
    <div>
      {/* Results count */}
      {pageInfo && (
        <p className="text-sm text-gray-500 mb-4">
          Showing {products.length} of {pageInfo.totalCount} products
        </p>
      )}

      {/* Product grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {products.map((product) => (
          <ProductCard
            key={product.id}
            product={product}
            onClick={() => onProductClick?.(product)}
          />
        ))}

        {/* Loading skeletons for infinite scroll */}
        {loading &&
          products.length > 0 &&
          Array.from({ length: 4 }).map((_, i) => (
            <ProductCardSkeleton key={`skeleton-${i}`} />
          ))}
      </div>

      {/* Load more trigger (for infinite scroll) */}
      <div ref={loadMoreRef} className="h-20 flex items-center justify-center">
        {loading && products.length > 0 && (
          <Loader2 className="w-6 h-6 text-primary-600 animate-spin" />
        )}
        {!loading && hasMore && (
          <button
            onClick={onLoadMore}
            className="px-6 py-2 text-primary-600 font-medium hover:text-primary-700"
          >
            Load more
          </button>
        )}
        {!hasMore && products.length > 0 && (
          <p className="text-sm text-gray-400">You've reached the end</p>
        )}
      </div>
    </div>
  );
}
