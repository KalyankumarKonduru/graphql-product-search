# GraphQL Product Search

A full-stack GraphQL application with **Go (gqlgen)** backend and **React (Apollo Client)** frontend. Features real-time search, cursor-based pagination, and optimistic UI updates.

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go)
![React](https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6?style=flat-square&logo=typescript)
![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=flat-square&logo=graphql)
![Apollo](https://img.shields.io/badge/Apollo_Client-311C87?style=flat-square&logo=apollographql)

## 🎯 Features

### Backend (Go + gqlgen)
- GraphQL API with type-safe resolvers
- Cursor-based pagination (Relay-style)
- Full-text search with filtering
- CORS-enabled for frontend

### Frontend (React + Apollo Client)
- Debounced search with autocomplete
- Infinite scroll pagination
- Normalized caching
- Optimistic UI updates
- Error boundaries
- Web Vitals monitoring
- MSW for testing

## 🏗️ Architecture

```
┌─────────────────────────────────────┐
│         React Frontend              │
│  ┌───────────────────────────────┐  │
│  │      Apollo Client            │  │
│  │  • Normalized Cache           │  │
│  │  • Optimistic Updates         │  │
│  │  • Error Handling             │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
                │
                │ GraphQL (HTTP)
                ▼
┌─────────────────────────────────────┐
│         Go Backend                  │
│  ┌───────────────────────────────┐  │
│  │         gqlgen                │  │
│  │  • Schema-first              │  │
│  │  • Type-safe resolvers       │  │
│  │  • Code generation           │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- Node.js 18+
- npm or yarn

### Backend

```bash
cd backend

# Download dependencies
go mod download

# Generate GraphQL code
go run github.com/99designs/gqlgen generate

# Run server
go run server.go
```

Server runs at `http://localhost:4000` with GraphQL Playground.

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev
```

Frontend runs at `http://localhost:5173`.

## 📊 GraphQL Schema

### Queries

```graphql
# Search products with pagination
query SearchProducts($search: String!, $first: Int!, $after: String) {
  products(
    first: $first
    after: $after
    filter: { search: $search }
    sort: { field: RATING, order: DESC }
  ) {
    edges {
      cursor
      node {
        id
        name
        price
        rating
        category { name }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
      totalCount
    }
  }
}

# Autocomplete suggestions
query Suggestions($query: String!) {
  searchSuggestions(query: $query, limit: 5)
}
```

### Mutations

```graphql
mutation CreateReview($input: CreateReviewInput!) {
  createReview(input: $input) {
    id
    rating
    comment
  }
}
```

## 🧪 Testing

### Frontend Tests (with MSW)

```bash
cd frontend
npm test
npm run test:coverage
```

Tests use Mock Service Worker to mock GraphQL:

```typescript
graphql.query('GetProducts', ({ variables }) => {
  return HttpResponse.json({
    data: { products: filteredProducts },
  });
});
```

## 📁 Project Structure

```
graphql-product-search/
├── backend/
│   ├── graph/
│   │   ├── schema.graphqls   # GraphQL schema
│   │   ├── resolver.go       # Query/Mutation resolvers
│   │   └── model/            # Data models
│   ├── internal/
│   │   └── database/         # In-memory database
│   ├── server.go             # Entry point
│   └── go.mod
│
└── frontend/
    ├── src/
    │   ├── components/       # React components
    │   ├── graphql/          # Queries, mutations, types
    │   ├── hooks/            # Custom hooks
    │   ├── lib/              # Apollo client config
    │   └── __mocks__/        # MSW handlers
    └── package.json
```

## 🔧 Key Implementation Details

### Apollo Client Cache

```typescript
const cache = new InMemoryCache({
  typePolicies: {
    Query: {
      fields: {
        products: {
          keyArgs: ['filter', 'sort'],
          merge(existing, incoming, { args }) {
            if (!args?.after) return incoming;
            return {
              ...incoming,
              edges: [...existing.edges, ...incoming.edges],
            };
          },
        },
      },
    },
  },
});
```

### gqlgen Resolver

```go
func (r *queryResolver) Products(
  ctx context.Context,
  first *int,
  after *string,
  filter *model.ProductFilterInput,
  sort *model.ProductSortInput,
) (*model.ProductConnection, error) {
  return database.DB.GetProducts(filter, sort, first, after, nil, nil), nil
}
```

## 📈 Performance

- **Debounced Search**: 300ms delay to reduce API calls
- **Cursor Pagination**: Efficient database queries
- **Normalized Cache**: Prevents duplicate data
- **Web Vitals**: LCP, FID, CLS monitoring

## 👤 Author

**Kalyankumar Konduru**

## 📄 License

MIT
