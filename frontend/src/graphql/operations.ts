import { gql } from '@apollo/client';

// ============== FRAGMENTS ==============
// Reusable fragments for consistent data fetching

export const PRODUCT_FRAGMENT = gql`
  fragment ProductFields on Product {
    id
    name
    description
    price
    imageUrl
    stock
    rating
    reviewCount
    createdAt
    category {
      id
      name
      slug
    }
  }
`;

export const PAGE_INFO_FRAGMENT = gql`
  fragment PageInfoFields on PageInfo {
    hasNextPage
    hasPreviousPage
    startCursor
    endCursor
    totalCount
  }
`;

export const REVIEW_FRAGMENT = gql`
  fragment ReviewFields on Review {
    id
    productId
    userId
    userName
    rating
    comment
    createdAt
  }
`;

// ============== QUERIES ==============

export const GET_PRODUCTS = gql`
  ${PRODUCT_FRAGMENT}
  ${PAGE_INFO_FRAGMENT}
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
          ...ProductFields
        }
      }
      pageInfo {
        ...PageInfoFields
      }
    }
  }
`;

export const GET_PRODUCT = gql`
  ${PRODUCT_FRAGMENT}
  query GetProduct($id: ID!) {
    product(id: $id) {
      ...ProductFields
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

export const GET_REVIEWS = gql`
  ${REVIEW_FRAGMENT}
  ${PAGE_INFO_FRAGMENT}
  query GetReviews($productId: ID!, $first: Int, $after: String) {
    reviews(productId: $productId, first: $first, after: $after) {
      edges {
        cursor
        node {
          ...ReviewFields
        }
      }
      pageInfo {
        ...PageInfoFields
      }
    }
  }
`;

export const GET_SEARCH_SUGGESTIONS = gql`
  query GetSearchSuggestions($query: String!, $limit: Int) {
    searchSuggestions(query: $query, limit: $limit)
  }
`;

// ============== MUTATIONS ==============

export const CREATE_REVIEW = gql`
  ${REVIEW_FRAGMENT}
  mutation CreateReview($input: CreateReviewInput!) {
    createReview(input: $input) {
      ...ReviewFields
    }
  }
`;

export const DELETE_REVIEW = gql`
  mutation DeleteReview($id: ID!) {
    deleteReview(id: $id)
  }
`;

export const CREATE_PRODUCT = gql`
  ${PRODUCT_FRAGMENT}
  mutation CreateProduct($input: CreateProductInput!) {
    createProduct(input: $input) {
      ...ProductFields
    }
  }
`;

export const UPDATE_PRODUCT = gql`
  ${PRODUCT_FRAGMENT}
  mutation UpdateProduct($id: ID!, $input: UpdateProductInput!) {
    updateProduct(id: $id, input: $input) {
      ...ProductFields
    }
  }
`;

export const DELETE_PRODUCT = gql`
  mutation DeleteProduct($id: ID!) {
    deleteProduct(id: $id)
  }
`;
