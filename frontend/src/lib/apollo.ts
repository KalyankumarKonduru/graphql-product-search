import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
  from,
} from '@apollo/client';
import { onError } from '@apollo/client/link/error';

// API URL - defaults to local development
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000/graphql';

// HTTP Link
const httpLink = new HttpLink({
  uri: API_URL,
});

// Error handling link
const errorLink = onError(({ graphQLErrors, networkError }: any) => {
  if (graphQLErrors) {
    graphQLErrors.forEach((err: { message: string }) => {
      console.error(`[GraphQL error]: ${err.message}`);
    });
  }

  if (networkError) {
    console.error(`[Network error]: ${networkError}`);
  }
});

// Cache configuration with type policies
const cache = new InMemoryCache({
  typePolicies: {
    Query: {
      fields: {
        // Cursor-based pagination for products
        products: {
          keyArgs: ['filter', 'sort'],
          merge(existing, incoming, { args }) {
            // If no cursor (first page), replace entirely
            if (!args?.after) {
              return incoming;
            }

            // Merge edges for infinite scroll
            const existingEdges = existing?.edges || [];
            const incomingEdges = incoming?.edges || [];

            return {
              ...incoming,
              edges: [...existingEdges, ...incomingEdges],
            };
          },
        },
        // Cursor-based pagination for reviews
        reviews: {
          keyArgs: ['productId'],
          merge(existing, incoming, { args }) {
            if (!args?.after) {
              return incoming;
            }

            const existingEdges = existing?.edges || [];
            const incomingEdges = incoming?.edges || [];

            return {
              ...incoming,
              edges: [...existingEdges, ...incomingEdges],
            };
          },
        },
      },
    },
    Product: {
      fields: {
        // Ensure category is properly cached
        category: {
          merge: true,
        },
      },
    },
  },
});

// Create Apollo Client
export const client = new ApolloClient({
  link: from([errorLink, httpLink]),
  cache,
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
      errorPolicy: 'all',
    },
    query: {
      fetchPolicy: 'cache-first',
      errorPolicy: 'all',
    },
    mutate: {
      errorPolicy: 'all',
    },
  },
});
