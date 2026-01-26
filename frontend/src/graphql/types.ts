export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  imageUrl: string;
  stock: number;
  rating: number;
  reviewCount: number;
  category: Category;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  productCount: number;
}

export interface PageInfo {
  hasNextPage: boolean;
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

export interface ProductFilterInput {
  search?: string;
  categoryId?: string;
  minPrice?: number;
  maxPrice?: number;
  minRating?: number;
  inStock?: boolean;
}

export interface ProductSortInput {
  field: 'NAME' | 'PRICE' | 'RATING' | 'CREATED_AT';
  order: 'ASC' | 'DESC';
}

export interface GetProductsResponse {
  products: ProductConnection;
}

export interface GetCategoriesResponse {
  categories: Category[];
}

export interface GetSearchSuggestionsResponse {
  searchSuggestions: string[];
}
