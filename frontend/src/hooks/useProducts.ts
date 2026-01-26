import { useQuery, useMutation } from '@apollo/client';
import { useCallback, useState, useEffect, useRef } from 'react';
import {
  GET_PRODUCTS,
  GET_PRODUCT,
  GET_CATEGORIES,
  GET_SEARCH_SUGGESTIONS,
  GET_REVIEWS,
  CREATE_REVIEW,
} from '../graphql/operations';
import type {
  GetProductsResponse,
  GetProductResponse,
  GetCategoriesResponse,
  GetSearchSuggestionsResponse,
  GetReviewsResponse,
  ProductFilterInput,
  ProductSortInput,
  CreateReviewInput,
  Product,
  Review,
} from '../graphql/types';

// ============== PRODUCTS HOOK ==============

interface UseProductsOptions {
  first?: number;
  filter?: ProductFilterInput;
  sort?: ProductSortInput;
}

export function useProducts(options: UseProductsOptions = {}) {
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
      variables: {
        after: data.products.pageInfo.endCursor,
      },
    });
  }, [data, fetchMore]);

  return {
    products: data?.products.edges.map((edge: { node: Product }) => edge.node) || [],
    pageInfo: data?.products.pageInfo,
    loading,
    error,
    loadMore,
    refetch,
    hasMore: data?.products.pageInfo.hasNextPage || false,
  };
}

// ============== SINGLE PRODUCT HOOK ==============

export function useProduct(id: string) {
  const { data, loading, error } = useQuery<GetProductResponse>(GET_PRODUCT, {
    variables: { id },
    skip: !id,
  });

  return {
    product: data?.product || null,
    loading,
    error,
  };
}

// ============== CATEGORIES HOOK ==============

export function useCategories() {
  const { data, loading, error } = useQuery<GetCategoriesResponse>(GET_CATEGORIES);

  return {
    categories: data?.categories || [],
    loading,
    error,
  };
}

// ============== SEARCH SUGGESTIONS HOOK ==============

export function useSearchSuggestions(query: string, limit = 5) {
  const { data, loading } = useQuery<GetSearchSuggestionsResponse>(
    GET_SEARCH_SUGGESTIONS,
    {
      variables: { query, limit },
      skip: query.length < 2,
      fetchPolicy: 'cache-first',
    }
  );

  return {
    suggestions: data?.searchSuggestions || [],
    loading,
  };
}

// ============== DEBOUNCED SEARCH HOOK ==============

export function useDebouncedSearch(delay = 300) {
  const [searchTerm, setSearchTerm] = useState('');
  const [debouncedTerm, setDebouncedTerm] = useState('');
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = setTimeout(() => {
      setDebouncedTerm(searchTerm);
    }, delay);

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [searchTerm, delay]);

  return {
    searchTerm,
    setSearchTerm,
    debouncedTerm,
    isDebouncing: searchTerm !== debouncedTerm,
  };
}

// ============== REVIEWS HOOK ==============

export function useReviews(productId: string, first = 5) {
  const { data, loading, error, fetchMore } = useQuery<GetReviewsResponse>(
    GET_REVIEWS,
    {
      variables: { productId, first },
      skip: !productId,
    }
  );

  const loadMore = useCallback(() => {
    if (!data?.reviews.pageInfo.hasNextPage) return;

    return fetchMore({
      variables: {
        after: data.reviews.pageInfo.endCursor,
      },
    });
  }, [data, fetchMore]);

  return {
    reviews: data?.reviews.edges.map((edge: { node: Review }) => edge.node) || [],
    pageInfo: data?.reviews.pageInfo,
    loading,
    error,
    loadMore,
    hasMore: data?.reviews.pageInfo.hasNextPage || false,
  };
}

// ============== CREATE REVIEW MUTATION ==============

export function useCreateReview() {
  const [mutate, { loading, error }] = useMutation(CREATE_REVIEW);

  const createReview = useCallback(
    async (input: CreateReviewInput) => {
      const result = await mutate({
        variables: { input },
        // Optimistic UI - show review immediately
        optimisticResponse: {
          createReview: {
            __typename: 'Review',
            id: `temp-${Date.now()}`,
            productId: input.productId,
            userId: 'current-user',
            userName: 'You',
            rating: input.rating,
            comment: input.comment,
            createdAt: new Date().toISOString(),
          },
        },
      });

      return result.data?.createReview;
    },
    [mutate]
  );

  return {
    createReview,
    loading,
    error,
  };
}

// ============== INFINITE SCROLL HOOK ==============

export function useInfiniteScroll(
  callback: () => void,
  hasMore: boolean,
  loading: boolean
) {
  const observerRef = useRef<IntersectionObserver | null>(null);
  const targetRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (loading) return;

    if (observerRef.current) {
      observerRef.current.disconnect();
    }

    observerRef.current = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading) {
          callback();
        }
      },
      { threshold: 0.1 }
    );

    if (targetRef.current) {
      observerRef.current.observe(targetRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [callback, hasMore, loading]);

  return targetRef;
}
