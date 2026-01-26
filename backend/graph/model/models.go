package model

import "time"

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

type Category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ProductCount int    `json:"productCount"`
}

type Review struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	UserID    string    `json:"userId"`
	UserName  string    `json:"userName"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
	TotalCount      int     `json:"totalCount"`
}

type ProductEdge struct {
	Cursor string   `json:"cursor"`
	Node   *Product `json:"node"`
}

type ProductConnection struct {
	Edges    []*ProductEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

type ReviewEdge struct {
	Cursor string  `json:"cursor"`
	Node   *Review `json:"node"`
}

type ReviewConnection struct {
	Edges    []*ReviewEdge `json:"edges"`
	PageInfo *PageInfo     `json:"pageInfo"`
}

type ProductFilterInput struct {
	Search     *string  `json:"search,omitempty"`
	CategoryID *string  `json:"categoryId,omitempty"`
	MinPrice   *float64 `json:"minPrice,omitempty"`
	MaxPrice   *float64 `json:"maxPrice,omitempty"`
	MinRating  *float64 `json:"minRating,omitempty"`
	InStock    *bool    `json:"inStock,omitempty"`
}

type ProductSortField string

const (
	ProductSortFieldName      ProductSortField = "NAME"
	ProductSortFieldPrice     ProductSortField = "PRICE"
	ProductSortFieldRating    ProductSortField = "RATING"
	ProductSortFieldCreatedAt ProductSortField = "CREATED_AT"
)

func (e ProductSortField) IsValid() bool {
	switch e {
	case ProductSortFieldName, ProductSortFieldPrice, ProductSortFieldRating, ProductSortFieldCreatedAt:
		return true
	}
	return false
}

func (e ProductSortField) String() string {
	return string(e)
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

func (e SortOrder) String() string {
	return string(e)
}

type ProductSortInput struct {
	Field ProductSortField `json:"field"`
	Order SortOrder        `json:"order"`
}

type CreateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"categoryId"`
	ImageURL    string  `json:"imageUrl"`
	Stock       int     `json:"stock"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	CategoryID  *string  `json:"categoryId,omitempty"`
	ImageURL    *string  `json:"imageUrl,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

type CreateReviewInput struct {
	ProductID string `json:"productId"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
