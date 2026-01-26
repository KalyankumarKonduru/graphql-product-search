# GraphQL Product Search

Full-stack GraphQL application with **Go (gqlgen)** backend and **React (Apollo Client)** frontend.

## Quick Start

### 1. Start Backend

```bash
cd backend

# Download dependencies
go mod tidy

# Generate GraphQL code
go run github.com/99designs/gqlgen generate

# IMPORTANT: After generate, open graph/schema.resolvers.go
# and replace the panic() lines with the implementations from README.md

# Run the server
go run server.go
```

Backend runs at: **http://localhost:4000**

### 2. Start Frontend

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev
```

Frontend runs at: **http://localhost:5173**

## Test the API

Open http://localhost:4000 and run:

```graphql
query {
  products(first: 5) {
    edges {
      node {
        id
        name
        price
        rating
        category {
          name
        }
      }
    }
    pageInfo {
      hasNextPage
      totalCount
    }
  }
}
```

## Features

- **GraphQL API** with cursor-based pagination
- **Full-text search** with debouncing
- **Category filtering** and sorting
- **Infinite scroll** with Apollo Client cache merging
- **Responsive UI** with Tailwind CSS

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go, gqlgen, Chi router |
| Frontend | React, TypeScript, Apollo Client |
| Styling | Tailwind CSS |
| Icons | Lucide React |
