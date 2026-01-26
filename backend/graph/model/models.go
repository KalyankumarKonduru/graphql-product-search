package model

import "time"

// Product represents a product in the catalog
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  string    `json:"categoryId"`
	ImageURL    string    `json:"imageUrl"`
	Stock       int       `json:"stock"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"reviewCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Category represents a product category
type Category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ProductCount int    `json:"productCount"`
}

// Review represents a product review
type Review struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	UserID    string    `json:"userId"`
	UserName  string    `json:"userName"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

// PageInfo for cursor-based pagination
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
	TotalCount      int     `json:"totalCount"`
}

// ProductEdge for connection pattern
type ProductEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Product `json:"node"`
}

// ProductConnection for paginated products
type ProductConnection struct {
	Edges    []*ProductEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

// ReviewEdge for connection pattern
type ReviewEdge struct {
	Cursor string  `json:"cursor"`
	Node   *Review `json:"node"`
}

// ReviewConnection for paginated reviews
type ReviewConnection struct {
	Edges    []*ReviewEdge `json:"edges"`
	PageInfo *PageInfo     `json:"pageInfo"`
}

// Filter and sort inputs
type ProductFilterInput struct {
	Search     *string  `json:"search"`
	CategoryID *string  `json:"categoryId"`
	MinPrice   *float64 `json:"minPrice"`
	MaxPrice   *float64 `json:"maxPrice"`
	MinRating  *float64 `json:"minRating"`
	InStock    *bool    `json:"inStock"`
}

type ProductSortInput struct {
	Field ProductSortField `json:"field"`
	Order SortOrder        `json:"order"`
}

type ProductSortField string

const (
	ProductSortFieldName      ProductSortField = "NAME"
	ProductSortFieldPrice     ProductSortField = "PRICE"
	ProductSortFieldRating    ProductSortField = "RATING"
	ProductSortFieldCreatedAt ProductSortField = "CREATED_AT"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// Mutation inputs
type CreateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"categoryId"`
	ImageURL    string  `json:"imageUrl"`
	Stock       int     `json:"stock"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	CategoryID  *string  `json:"categoryId"`
	ImageURL    *string  `json:"imageUrl"`
	Stock       *int     `json:"stock"`
}

type CreateReviewInput struct {
	ProductID string `json:"productId"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
