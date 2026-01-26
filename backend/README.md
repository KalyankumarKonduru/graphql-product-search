# GraphQL Product Search - Go Backend

A GraphQL API built with **Go** and **gqlgen** for product search with cursor-based pagination.

## Tech Stack

- **Go 1.22+**
- **gqlgen** - Type-safe GraphQL server
- **Chi** - HTTP router
- **CORS** middleware

## Features

- 📦 Product CRUD operations
- 🔍 Full-text search with filters
- 📄 Cursor-based pagination (Relay-style)
- ⭐ Reviews with ratings
- 🏷️ Category management
- 💡 Search suggestions (autocomplete)

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Make (optional)

### Installation

```bash
# Clone and navigate
cd backend

# Download dependencies
go mod download

# Generate GraphQL code (first time only)
go run github.com/99designs/gqlgen generate

# Run the server
go run server.go
```

### Endpoints

| Endpoint | Description |
|----------|-------------|
| `http://localhost:4000/` | GraphQL Playground |
| `http://localhost:4000/graphql` | GraphQL API |
| `http://localhost:4000/health` | Health check |

## GraphQL Schema Highlights

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
      }
    }
    pageInfo {
      hasNextPage
      endCursor
      totalCount
    }
  }
}

# Get search suggestions
query Suggestions($query: String!) {
  searchSuggestions(query: $query, limit: 5)
}
```

### Mutations

```graphql
# Create a review
mutation CreateReview($input: CreateReviewInput!) {
  createReview(input: $input) {
    id
    rating
    comment
  }
}
```

## Project Structure

```
backend/
├── graph/
│   ├── schema.graphqls    # GraphQL schema
│   ├── model/
│   │   └── models.go      # Data models
│   ├── resolver.go        # Query/Mutation resolvers
│   └── generated.go       # Auto-generated (gqlgen)
├── internal/
│   └── database/
│       └── database.go    # In-memory DB with seed data
├── server.go              # Entry point
├── gqlgen.yml             # gqlgen configuration
└── go.mod                 # Dependencies
```

## Development

```bash
# Regenerate after schema changes
go run github.com/99designs/gqlgen generate

# Run with hot reload (install air first)
air

# Run tests
go test ./...
```

## Deployment

```bash
# Build binary
go build -o server .

# Run
./server
```

For Docker:
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:latest
COPY --from=builder /app/server /server
EXPOSE 4000
CMD ["/server"]
```
