import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
  ApolloLink,
  from,
} from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { RetryLink } from '@apollo/client/link/retry';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000/graphql';

// Error handling link - logs and categorizes errors
const errorLink = onError(({ graphQLErrors, networkError, operation }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path, extensions }) => {
      const isAuthError = message.includes('unauthorized') || extensions?.code === 'UNAUTHENTICATED';

      if (isAuthError) {
        console.warn(`[Auth Error] Operation: ${operation.operationName}`);
        // Could dispatch to auth state management here
      } else {
        console.error(
          `[GraphQL Error] Message: ${message}, Location: ${JSON.stringify(locations)}, Path: ${path}`
        );
      }

      // Production error reporting
      if (import.meta.env.PROD) {
        reportError({ type: 'graphql', message, operation: operation.operationName, path });
      }
    });
  }

  if (networkError) {
    console.error(`[Network Error] ${networkError.message}`);

    if (import.meta.env.PROD) {
      reportError({ type: 'network', message: networkError.message });
    }
  }
});

// Retry link - automatically retries failed network requests
const retryLink = new RetryLink({
  delay: {
    initial: 300,
    max: 5000,
    jitter: true,
  },
  attempts: {
    max: 3,
    retryIf: (error) => !!error && !error.message?.includes('unauthorized'),
  },
});

// HTTP link with credentials
const httpLink = new HttpLink({
  uri: API_URL,
  credentials: 'same-origin',
  headers: {
    'X-Client-Version': import.meta.env.VITE_APP_VERSION || '1.0.0',
  },
});

// Request ID link - adds a unique ID to each request for tracing
const requestIdLink = new ApolloLink((operation, forward) => {
  operation.setContext(({ headers = {} }) => ({
    headers: {
      ...headers,
      'X-Request-ID': crypto.randomUUID?.() || `${Date.now()}-${Math.random().toString(36).slice(2)}`,
    },
  }));
  return forward(operation);
});

// Performance tracking link
const performanceLink = new ApolloLink((operation, forward) => {
  const startTime = performance.now();

  return forward(operation).map((response) => {
    const duration = performance.now() - startTime;

    if (import.meta.env.DEV) {
      const color = duration > 1000 ? '#ff4444' : duration > 300 ? '#ffaa00' : '#44aa44';
      console.log(
        `%c[GraphQL] ${operation.operationName}: ${duration.toFixed(1)}ms`,
        `color: ${color}; font-weight: bold`
      );
    }

    // Track slow queries in production
    if (import.meta.env.PROD && duration > 2000) {
      reportError({
        type: 'slow_query',
        operation: operation.operationName,
        duration,
      });
    }

    return response;
  });
});

// Cache configuration with optimized policies
const cache = new InMemoryCache({
  typePolicies: {
    Query: {
      fields: {
        products: {
          keyArgs: ['filter', 'sort'],
          merge(existing, incoming, { args }) {
            if (!args?.after) {
              return incoming;
            }
            return {
              ...incoming,
              edges: [...(existing?.edges || []), ...incoming.edges],
            };
          },
        },
      },
    },
    Product: {
      keyFields: ['id'],
    },
    Category: {
      keyFields: ['id'],
    },
    Review: {
      keyFields: ['id'],
    },
  },
});

// Compose the link chain
const link = from([
  requestIdLink,
  performanceLink,
  errorLink,
  retryLink,
  httpLink,
]);

export const client = new ApolloClient({
  link,
  cache,
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
      errorPolicy: 'all',
    },
    query: {
      errorPolicy: 'all',
    },
    mutate: {
      errorPolicy: 'all',
    },
  },
  connectToDevTools: import.meta.env.DEV,
});

// Error reporting utility (replace with your monitoring service)
function reportError(data: Record<string, unknown>) {
  // Example: Send to monitoring endpoint
  // navigator.sendBeacon('/api/errors', JSON.stringify(data));
  if (import.meta.env.DEV) {
    console.log('[Error Report]', data);
  }
}
