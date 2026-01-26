import { useQuery } from '@apollo/client';
import { useCallback, useState, useEffect, useRef } from 'react';
import { GET_PRODUCTS, GET_CATEGORIES, GET_SEARCH_SUGGESTIONS } from '../graphql/operations';
import type {
  GetProductsResponse,
  GetCategoriesResponse,
  GetSearchSuggestionsResponse,
  ProductFilterInput,
  ProductSortInput,
  Product,
} from '../graphql/types';

export function useProducts(options: {
  first?: number;
  filter?: ProductFilterInput;
  sort?: ProductSortInput;
} = {}) {
  const { first = 12, filter, sort } = options;

  const { data, loading, error, fetchMore, refetch } = useQuery<GetProductsResponse>(
    GET_PRODUCTS,
    {
      variables: { first, filter, sort },
      notifyOnNetworkStatusChange: true,
    }
  );

  const loadMore = useCallback(() => {
    if (!data?.products.pageInfo.hasNextPage) return;
    return fetchMore({
      variables: { after: data.products.pageInfo.endCursor },
    });
  }, [data, fetchMore]);

  return {
    products: data?.products.edges.map((e: { node: Product }) => e.node) || [],
    pageInfo: data?.products.pageInfo,
    loading,
    error,
    loadMore,
    refetch,
    hasMore: data?.products.pageInfo.hasNextPage || false,
  };
}

export function useCategories() {
  const { data, loading, error } = useQuery<GetCategoriesResponse>(GET_CATEGORIES);
  return {
    categories: data?.categories || [],
    loading,
    error,
  };
}

export function useSearchSuggestions(query: string, limit = 5) {
  const { data, loading } = useQuery<GetSearchSuggestionsResponse>(
    GET_SEARCH_SUGGESTIONS,
    {
      variables: { query, limit },
      skip: query.length < 2,
    }
  );
  return {
    suggestions: data?.searchSuggestions || [],
    loading,
  };
}

export function useDebouncedSearch(delay = 300) {
  const [searchTerm, setSearchTerm] = useState('');
  const [debouncedTerm, setDebouncedTerm] = useState('');
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    timeoutRef.current = setTimeout(() => setDebouncedTerm(searchTerm), delay);
    return () => {
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
  }, [searchTerm, delay]);

  return { searchTerm, setSearchTerm, debouncedTerm };
}

export function useInfiniteScroll(callback: () => void, hasMore: boolean, loading: boolean) {
  const observerRef = useRef<IntersectionObserver | null>(null);
  const targetRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (loading) return;
    if (observerRef.current) observerRef.current.disconnect();

    observerRef.current = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting && hasMore && !loading) {
        callback();
      }
    });

    if (targetRef.current) observerRef.current.observe(targetRef.current);
    return () => observerRef.current?.disconnect();
  }, [callback, hasMore, loading]);

  return targetRef;
}
