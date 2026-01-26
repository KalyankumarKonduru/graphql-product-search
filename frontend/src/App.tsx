import { useState, useCallback } from 'react';
import { ApolloProvider } from '@apollo/client';
import { client } from './lib/apollo';
import { useProducts } from './hooks/useProducts';
import { SearchBar } from './components/SearchBar';
import { ProductGrid } from './components/ProductGrid';
import { Filters } from './components/Filters';
import { ShoppingBag } from 'lucide-react';
import type { ProductFilterInput, ProductSortInput } from './graphql/types';

function ProductSearch() {
  const [filter, setFilter] = useState<ProductFilterInput>({});
  const [sort, setSort] = useState<ProductSortInput>({
    field: 'CREATED_AT',
    order: 'DESC',
  });

  const { products, pageInfo, loading, error, loadMore, hasMore } = useProducts({
    first: 12,
    filter,
    sort,
  });

  const handleSearch = useCallback((query: string) => {
    setFilter((prev) => ({ ...prev, search: query || undefined }));
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-2">
              <div className="w-10 h-10 bg-primary-600 rounded-xl flex items-center justify-center">
                <ShoppingBag className="w-6 h-6 text-white" />
              </div>
              <span className="text-xl font-bold text-gray-900">GraphQL Store</span>
            </div>
            <div className="hidden md:block flex-1 max-w-xl mx-8">
              <SearchBar onSearch={handleSearch} />
            </div>
          </div>
          <div className="md:hidden pb-4">
            <SearchBar onSearch={handleSearch} />
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">
            {filter.search ? `Results for "${filter.search}"` : 'All Products'}
          </h1>
        </div>

        <div className="mb-6">
          <Filters
            filter={filter}
            sort={sort}
            onFilterChange={setFilter}
            onSortChange={setSort}
          />
        </div>

        <ProductGrid
          products={products}
          loading={loading}
          error={error}
          pageInfo={pageInfo}
          hasMore={hasMore}
          onLoadMore={loadMore}
        />
      </main>

      <footer className="bg-white border-t border-gray-200 mt-12">
        <div className="max-w-7xl mx-auto px-4 py-8 text-center text-gray-500 text-sm">
          Built with Go (gqlgen) + React + Apollo Client
        </div>
      </footer>
    </div>
  );
}

function App() {
  return (
    <ApolloProvider client={client}>
      <ProductSearch />
    </ApolloProvider>
  );
}

export default App;
