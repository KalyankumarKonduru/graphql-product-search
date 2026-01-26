import { useState, useCallback } from 'react';
import { ApolloProvider } from '@apollo/client';
import { client } from './lib/apollo';
import { useProducts } from './hooks/useProducts';
import { SearchBar } from './components/SearchBar';
import { ProductGrid } from './components/ProductGrid';
import { Filters } from './components/Filters';
import { ErrorBoundary } from './components/ErrorBoundary';
import { WebVitalsDebugPanel } from './hooks/useWebVitals';
import { ShoppingBag } from 'lucide-react';
import type {
  ProductFilterInput,
  ProductSortInput,
  ProductSortField,
  SortOrder,
} from './graphql/types';

function ProductSearch() {
  // State for filters and sorting
  const [filter, setFilter] = useState<ProductFilterInput>({});
  const [sort, setSort] = useState<ProductSortInput>({
    field: 'CREATED_AT' as ProductSortField,
    order: 'DESC' as SortOrder,
  });

  // Fetch products with current filters
  const {
    products,
    pageInfo,
    loading,
    error,
    loadMore,
    refetch,
    hasMore,
  } = useProducts({
    first: 12,
    filter,
    sort,
  });

  // Handle search
  const handleSearch = useCallback((query: string) => {
    setFilter((prev) => ({
      ...prev,
      search: query || undefined,
    }));
  }, []);

  // Handle filter changes
  const handleFilterChange = useCallback((newFilter: ProductFilterInput) => {
    setFilter(newFilter);
  }, []);

  // Handle sort changes
  const handleSortChange = useCallback((newSort: ProductSortInput) => {
    setSort(newSort);
  }, []);

  // Handle product click
  const handleProductClick = useCallback((product: any) => {
    console.log('Product clicked:', product);
    // Could open modal or navigate to product page
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <div className="flex items-center gap-2">
              <div className="w-10 h-10 bg-primary-600 rounded-xl flex items-center justify-center">
                <ShoppingBag className="w-6 h-6 text-white" />
              </div>
              <span className="text-xl font-bold text-gray-900">GraphQL Store</span>
            </div>

            {/* Search - Desktop */}
            <div className="hidden md:block flex-1 max-w-xl mx-8">
              <SearchBar
                onSearch={handleSearch}
                placeholder="Search products..."
                initialValue={filter.search}
              />
            </div>

            {/* Actions */}
            <div className="flex items-center gap-4">
              <button className="text-sm font-medium text-gray-600 hover:text-gray-900">
                Sign In
              </button>
              <button className="relative p-2 text-gray-600 hover:text-gray-900">
                <ShoppingBag className="w-6 h-6" />
                <span className="absolute -top-1 -right-1 w-5 h-5 bg-primary-600 text-white text-xs rounded-full flex items-center justify-center">
                  0
                </span>
              </button>
            </div>
          </div>

          {/* Search - Mobile */}
          <div className="md:hidden pb-4">
            <SearchBar
              onSearch={handleSearch}
              placeholder="Search products..."
              initialValue={filter.search}
            />
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Page Title */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">
            {filter.search ? `Results for "${filter.search}"` : 'All Products'}
          </h1>
          <p className="text-gray-500 mt-1">
            Find the best products with GraphQL-powered search
          </p>
        </div>

        {/* Filters */}
        <div className="mb-6">
          <Filters
            filter={filter}
            sort={sort}
            onFilterChange={handleFilterChange}
            onSortChange={handleSortChange}
          />
        </div>

        {/* Product Grid */}
        <ProductGrid
          products={products}
          loading={loading}
          error={error}
          pageInfo={pageInfo}
          hasMore={hasMore}
          onLoadMore={loadMore}
          onProductClick={handleProductClick}
          onRetry={refetch}
        />
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <p className="text-gray-500 text-sm">
              Built with Go (gqlgen) + React + Apollo Client
            </p>
            <div className="flex items-center gap-6 text-sm text-gray-500">
              <a href="#" className="hover:text-gray-900">
                GitHub
              </a>
              <a href="#" className="hover:text-gray-900">
                Documentation
              </a>
              <a href="#" className="hover:text-gray-900">
                GraphQL Playground
              </a>
            </div>
          </div>
        </div>
      </footer>

      {/* Web Vitals Debug Panel (dev only) */}
      <WebVitalsDebugPanel />
    </div>
  );
}

function App() {
  return (
    <ApolloProvider client={client}>
      <ErrorBoundary>
        <ProductSearch />
      </ErrorBoundary>
    </ApolloProvider>
  );
}

export default App;
