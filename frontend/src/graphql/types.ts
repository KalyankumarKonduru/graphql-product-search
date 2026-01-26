// GraphQL Types - matches backend schema

export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  imageUrl: string;
  stock: number;
  rating: number;
  reviewCount: number;
  createdAt: string;
  updatedAt: string;
  category: Category;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  productCount: number;
}

export interface Review {
  id: string;
  productId: string;
  userId: string;
  userName: string;
  rating: number;
  comment: string;
  createdAt: string;
}

export interface PageInfo {
  hasNextPage: boolean;
  hasPreviousPage: boolean;
  startCursor: string | null;
  endCursor: string | null;
  totalCount: number;
}

export interface ProductEdge {
  cursor: string;
  node: Product;
}

export interface ProductConnection {
  edges: ProductEdge[];
  pageInfo: PageInfo;
}

export interface ReviewEdge {
  cursor: string;
  node: Review;
}

export interface ReviewConnection {
  edges: ReviewEdge[];
  pageInfo: PageInfo;
}

// Filter and Sort types
export interface ProductFilterInput {
  search?: string;
  categoryId?: string;
  minPrice?: number;
  maxPrice?: number;
  minRating?: number;
  inStock?: boolean;
}

export enum ProductSortField {
  NAME = 'NAME',
  PRICE = 'PRICE',
  RATING = 'RATING',
  CREATED_AT = 'CREATED_AT',
}

export enum SortOrder {
  ASC = 'ASC',
  DESC = 'DESC',
}

export interface ProductSortInput {
  field: ProductSortField;
  order: SortOrder;
}

// Mutation inputs
export interface CreateReviewInput {
  productId: string;
  rating: number;
  comment: string;
}

export interface CreateProductInput {
  name: string;
  description: string;
  price: number;
  categoryId: string;
  imageUrl: string;
  stock: number;
}

export interface UpdateProductInput {
  name?: string;
  description?: string;
  price?: number;
  categoryId?: string;
  imageUrl?: string;
  stock?: number;
}

// Query response types
export interface GetProductsResponse {
  products: ProductConnection;
}

export interface GetProductResponse {
  product: Product | null;
}

export interface GetCategoriesResponse {
  categories: Category[];
}

export interface GetReviewsResponse {
  reviews: ReviewConnection;
}

export interface GetSearchSuggestionsResponse {
  searchSuggestions: string[];
}
