import { ProductCard, ProductCardSkeleton } from './ProductCard';
import { useInfiniteScroll } from '../hooks/useProducts';
import { Loader2, AlertCircle } from 'lucide-react';
import type { Product, PageInfo } from '../graphql/types';

interface ProductGridProps {
  products: Product[];
  loading: boolean;
  error?: Error;
  pageInfo?: PageInfo;
  hasMore: boolean;
  onLoadMore: () => void;
}

export function ProductGrid({
  products,
  loading,
  error,
  pageInfo,
  hasMore,
  onLoadMore,
}: ProductGridProps) {
  const loadMoreRef = useInfiniteScroll(onLoadMore, hasMore, loading);

  if (error) {
    return (
      <div className="py-12 text-center">
        <AlertCircle className="w-12 h-12 text-red-500 mx-auto mb-4" />
        <p className="text-red-600">Error loading products. Please try again.</p>
      </div>
    );
  }

  if (loading && products.length === 0) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {Array.from({ length: 8 }).map((_, i) => (
          <ProductCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="py-16 text-center">
        <p className="text-gray-500">No products found. Try a different search.</p>
      </div>
    );
  }

  return (
    <div>
      {pageInfo && (
        <p className="text-sm text-gray-500 mb-4">
          Showing {products.length} of {pageInfo.totalCount} products
        </p>
      )}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
        {loading && Array.from({ length: 4 }).map((_, i) => (
          <ProductCardSkeleton key={`skeleton-${i}`} />
        ))}
      </div>
      <div ref={loadMoreRef} className="h-20 flex items-center justify-center">
        {loading && products.length > 0 && (
          <Loader2 className="w-6 h-6 text-primary-600 animate-spin" />
        )}
        {!hasMore && products.length > 0 && (
          <p className="text-sm text-gray-400">You've reached the end</p>
        )}
      </div>
    </div>
  );
}
