import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000/graphql';

const httpLink = new HttpLink({
  uri: API_URL,
});

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
  },
});

export const client = new ApolloClient({
  link: httpLink,
  cache,
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
    },
  },
});
