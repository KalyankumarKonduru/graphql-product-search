import { gql } from '@apollo/client';

export const GET_PRODUCTS = gql`
  query GetProducts(
    $first: Int
    $after: String
    $filter: ProductFilterInput
    $sort: ProductSortInput
  ) {
    products(first: $first, after: $after, filter: $filter, sort: $sort) {
      edges {
        cursor
        node {
          id
          name
          description
          price
          imageUrl
          stock
          rating
          reviewCount
          category {
            id
            name
          }
        }
      }
      pageInfo {
        hasNextPage
        endCursor
        totalCount
      }
    }
  }
`;

export const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
      slug
      productCount
    }
  }
`;

export const GET_SEARCH_SUGGESTIONS = gql`
  query GetSearchSuggestions($query: String!, $limit: Int) {
    searchSuggestions(query: $query, limit: $limit)
  }
`;
