package database

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kalyankumar/graphql-product-search/graph/model"
)

// Database is an in-memory store for demo purposes
type Database struct {
	mu         sync.RWMutex
	products   map[string]*model.Product
	categories map[string]*model.Category
	reviews    map[string]*model.Review
}

// Global database instance
var DB *Database

// Initialize creates and seeds the database
func Initialize() {
	DB = &Database{
		products:   make(map[string]*model.Product),
		categories: make(map[string]*model.Category),
		reviews:    make(map[string]*model.Review),
	}
	DB.seed()
}

func (db *Database) seed() {
	// Seed categories
	categories := []*model.Category{
		{ID: "cat-1", Name: "Electronics", Slug: "electronics", ProductCount: 0},
		{ID: "cat-2", Name: "Clothing", Slug: "clothing", ProductCount: 0},
		{ID: "cat-3", Name: "Home & Garden", Slug: "home-garden", ProductCount: 0},
		{ID: "cat-4", Name: "Sports", Slug: "sports", ProductCount: 0},
		{ID: "cat-5", Name: "Books", Slug: "books", ProductCount: 0},
	}

	for _, c := range categories {
		db.categories[c.ID] = c
	}

	// Seed products
	products := []*model.Product{
		{ID: "prod-1", Name: "MacBook Pro 16\"", Description: "Apple M3 Pro chip, 18GB RAM, 512GB SSD", Price: 2499.00, CategoryID: "cat-1", ImageURL: "https://picsum.photos/seed/macbook/400/400", Stock: 15, Rating: 4.8, ReviewCount: 124, CreatedAt: time.Now().AddDate(0, -3, 0), UpdatedAt: time.Now()},
		{ID: "prod-2", Name: "Sony WH-1000XM5", Description: "Industry-leading noise canceling headphones", Price: 349.99, CategoryID: "cat-1", ImageURL: "https://picsum.photos/seed/sony/400/400", Stock: 42, Rating: 4.7, ReviewCount: 89, CreatedAt: time.Now().AddDate(0, -2, 0), UpdatedAt: time.Now()},
		{ID: "prod-3", Name: "Nike Air Max 270", Description: "Men's running shoes with Max Air unit", Price: 150.00, CategoryID: "cat-4", ImageURL: "https://picsum.photos/seed/nike/400/400", Stock: 28, Rating: 4.5, ReviewCount: 203, CreatedAt: time.Now().AddDate(0, -1, 0), UpdatedAt: time.Now()},
		{ID: "prod-4", Name: "Levi's 501 Original", Description: "Classic straight fit jeans", Price: 69.50, CategoryID: "cat-2", ImageURL: "https://picsum.photos/seed/levis/400/400", Stock: 65, Rating: 4.4, ReviewCount: 156, CreatedAt: time.Now().AddDate(0, -4, 0), UpdatedAt: time.Now()},
		{ID: "prod-5", Name: "Kindle Paperwhite", Description: "6.8\" display, adjustable warm light", Price: 139.99, CategoryID: "cat-1", ImageURL: "https://picsum.photos/seed/kindle/400/400", Stock: 33, Rating: 4.6, ReviewCount: 312, CreatedAt: time.Now().AddDate(0, -5, 0), UpdatedAt: time.Now()},
		{ID: "prod-6", Name: "Dyson V15 Detect", Description: "Cordless vacuum with laser dust detection", Price: 749.99, CategoryID: "cat-3", ImageURL: "https://picsum.photos/seed/dyson/400/400", Stock: 12, Rating: 4.9, ReviewCount: 78, CreatedAt: time.Now().AddDate(0, -1, -15), UpdatedAt: time.Now()},
		{ID: "prod-7", Name: "The Pragmatic Programmer", Description: "Your Journey to Mastery, 20th Anniversary Edition", Price: 49.99, CategoryID: "cat-5", ImageURL: "https://picsum.photos/seed/pragmatic/400/400", Stock: 100, Rating: 4.8, ReviewCount: 445, CreatedAt: time.Now().AddDate(-1, 0, 0), UpdatedAt: time.Now()},
		{ID: "prod-8", Name: "Samsung Galaxy S24 Ultra", Description: "AI-powered smartphone with S Pen", Price: 1299.99, CategoryID: "cat-1", ImageURL: "https://picsum.photos/seed/samsung/400/400", Stock: 22, Rating: 4.6, ReviewCount: 167, CreatedAt: time.Now().AddDate(0, 0, -20), UpdatedAt: time.Now()},
		{ID: "prod-9", Name: "Patagonia Down Sweater", Description: "Lightweight, windproof insulation", Price: 229.00, CategoryID: "cat-2", ImageURL: "https://picsum.photos/seed/patagonia/400/400", Stock: 18, Rating: 4.7, ReviewCount: 92, CreatedAt: time.Now().AddDate(0, -2, -10), UpdatedAt: time.Now()},
		{ID: "prod-10", Name: "YETI Rambler 26oz", Description: "Vacuum insulated bottle with chug cap", Price: 40.00, CategoryID: "cat-4", ImageURL: "https://picsum.photos/seed/yeti/400/400", Stock: 55, Rating: 4.9, ReviewCount: 234, CreatedAt: time.Now().AddDate(0, -3, -5), UpdatedAt: time.Now()},
		{ID: "prod-11", Name: "iPad Pro 12.9\"", Description: "M2 chip, Liquid Retina XDR display", Price: 1099.00, CategoryID: "cat-1", ImageURL: "https://picsum.photos/seed/ipad/400/400", Stock: 19, Rating: 4.8, ReviewCount: 198, CreatedAt: time.Now().AddDate(0, -1, -5), UpdatedAt: time.Now()},
		{ID: "prod-12", Name: "Herman Miller Aeron", Description: "Ergonomic office chair, size B", Price: 1395.00, CategoryID: "cat-3", ImageURL: "https://picsum.photos/seed/aeron/400/400", Stock: 8, Rating: 4.7, ReviewCount: 67, CreatedAt: time.Now().AddDate(0, -6, 0), UpdatedAt: time.Now()},
		{ID: "prod-13", Name: "Clean Code", Description: "A Handbook of Agile Software Craftsmanship", Price: 39.99, CategoryID: "cat-5", ImageURL: "https://picsum.photos/seed/cleancode/400/400", Stock: 85, Rating: 4.6, ReviewCount: 523, CreatedAt: time.Now().AddDate(-2, 0, 0), UpdatedAt: time.Now()},
		{ID: "prod-14", Name: "Adidas Ultraboost 22", Description: "High-performance running shoes", Price: 190.00, CategoryID: "cat-4", ImageURL: "https://picsum.photos/seed/adidas/400/400", Stock: 31, Rating: 4.5, ReviewCount: 178, CreatedAt: time.Now().AddDate(0, -2, -20), UpdatedAt: time.Now()},
		{ID: "prod-15", Name: "Instant Pot Duo 7-in-1", Description: "Electric pressure cooker, 6 quart", Price: 89.95, CategoryID: "cat-3", ImageURL: "https://picsum.photos/seed/instantpot/400/400", Stock: 47, Rating: 4.7, ReviewCount: 412, CreatedAt: time.Now().AddDate(0, -4, -10), UpdatedAt: time.Now()},
	}

	for _, p := range products {
		db.products[p.ID] = p
		if cat, ok := db.categories[p.CategoryID]; ok {
			cat.ProductCount++
		}
	}

	// Seed reviews
	reviews := []*model.Review{
		{ID: "rev-1", ProductID: "prod-1", UserID: "user-1", UserName: "John D.", Rating: 5, Comment: "Best laptop I've ever owned!", CreatedAt: time.Now().AddDate(0, 0, -5)},
		{ID: "rev-2", ProductID: "prod-1", UserID: "user-2", UserName: "Sarah M.", Rating: 4, Comment: "Great performance, but pricey.", CreatedAt: time.Now().AddDate(0, 0, -3)},
		{ID: "rev-3", ProductID: "prod-2", UserID: "user-3", UserName: "Mike R.", Rating: 5, Comment: "Noise canceling is incredible!", CreatedAt: time.Now().AddDate(0, 0, -7)},
		{ID: "rev-4", ProductID: "prod-7", UserID: "user-4", UserName: "Emily K.", Rating: 5, Comment: "A must-read for every developer.", CreatedAt: time.Now().AddDate(0, 0, -14)},
		{ID: "rev-5", ProductID: "prod-6", UserID: "user-5", UserName: "Alex P.", Rating: 5, Comment: "The laser feature is amazing!", CreatedAt: time.Now().AddDate(0, 0, -2)},
	}

	for _, r := range reviews {
		db.reviews[r.ID] = r
	}
}

// GetProduct returns a product by ID
func (db *Database) GetProduct(id string) *model.Product {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.products[id]
}

// GetProducts returns filtered and paginated products
func (db *Database) GetProducts(filter *model.ProductFilterInput, sortInput *model.ProductSortInput, first *int, after *string, last *int, before *string) *model.ProductConnection {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Convert map to slice
	var products []*model.Product
	for _, p := range db.products {
		products = append(products, p)
	}

	// Apply filters
	products = db.filterProducts(products, filter)

	// Apply sorting
	db.sortProducts(products, sortInput)

	// Apply pagination
	return db.paginateProducts(products, first, after, last, before)
}

func (db *Database) filterProducts(products []*model.Product, filter *model.ProductFilterInput) []*model.Product {
	if filter == nil {
		return products
	}

	var filtered []*model.Product
	for _, p := range products {
		// Search filter (name or description)
		if filter.Search != nil && *filter.Search != "" {
			search := strings.ToLower(*filter.Search)
			if !strings.Contains(strings.ToLower(p.Name), search) &&
				!strings.Contains(strings.ToLower(p.Description), search) {
				continue
			}
		}

		// Category filter
		if filter.CategoryID != nil && p.CategoryID != *filter.CategoryID {
			continue
		}

		// Price range filter
		if filter.MinPrice != nil && p.Price < *filter.MinPrice {
			continue
		}
		if filter.MaxPrice != nil && p.Price > *filter.MaxPrice {
			continue
		}

		// Rating filter
		if filter.MinRating != nil && p.Rating < *filter.MinRating {
			continue
		}

		// Stock filter
		if filter.InStock != nil && *filter.InStock && p.Stock <= 0 {
			continue
		}

		filtered = append(filtered, p)
	}

	return filtered
}

func (db *Database) sortProducts(products []*model.Product, sortInput *model.ProductSortInput) {
	if sortInput == nil {
		// Default sort by created date desc
		sort.Slice(products, func(i, j int) bool {
			return products[i].CreatedAt.After(products[j].CreatedAt)
		})
		return
	}

	sort.Slice(products, func(i, j int) bool {
		var less bool
		switch sortInput.Field {
		case model.ProductSortFieldName:
			less = products[i].Name < products[j].Name
		case model.ProductSortFieldPrice:
			less = products[i].Price < products[j].Price
		case model.ProductSortFieldRating:
			less = products[i].Rating < products[j].Rating
		case model.ProductSortFieldCreatedAt:
			less = products[i].CreatedAt.Before(products[j].CreatedAt)
		default:
			less = products[i].CreatedAt.Before(products[j].CreatedAt)
		}

		if sortInput.Order == model.SortOrderDesc {
			return !less
		}
		return less
	})
}

func (db *Database) paginateProducts(products []*model.Product, first *int, after *string, last *int, before *string) *model.ProductConnection {
	totalCount := len(products)

	// Find cursor positions
	startIdx := 0
	endIdx := totalCount

	if after != nil {
		cursor := decodeCursor(*after)
		for i, p := range products {
			if p.ID == cursor {
				startIdx = i + 1
				break
			}
		}
	}

	if before != nil {
		cursor := decodeCursor(*before)
		for i, p := range products {
			if p.ID == cursor {
				endIdx = i
				break
			}
		}
	}

	// Slice products
	if startIdx > endIdx {
		startIdx = endIdx
	}
	products = products[startIdx:endIdx]

	// Apply first/last limits
	if first != nil && *first < len(products) {
		products = products[:*first]
	}
	if last != nil && *last < len(products) {
		products = products[len(products)-*last:]
	}

	// Build edges
	edges := make([]*model.ProductEdge, len(products))
	for i, p := range products {
		edges[i] = &model.ProductEdge{
			Cursor: encodeCursor(p.ID),
			Node:   p,
		}
	}

	// Build page info
	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &model.ProductConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     endIdx < totalCount,
			HasPreviousPage: startIdx > 0,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
			TotalCount:      totalCount,
		},
	}
}

// GetCategories returns all categories
func (db *Database) GetCategories() []*model.Category {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var categories []*model.Category
	for _, c := range db.categories {
		categories = append(categories, c)
	}
	return categories
}

// GetCategory returns a category by ID
func (db *Database) GetCategory(id string) *model.Category {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.categories[id]
}

// GetReviews returns reviews for a product
func (db *Database) GetReviews(productID string, first *int, after *string) *model.ReviewConnection {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var reviews []*model.Review
	for _, r := range db.reviews {
		if r.ProductID == productID {
			reviews = append(reviews, r)
		}
	}

	// Sort by date desc
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})

	totalCount := len(reviews)

	// Apply cursor pagination
	startIdx := 0
	if after != nil {
		cursor := decodeCursor(*after)
		for i, r := range reviews {
			if r.ID == cursor {
				startIdx = i + 1
				break
			}
		}
	}

	reviews = reviews[startIdx:]
	if first != nil && *first < len(reviews) {
		reviews = reviews[:*first]
	}

	// Build edges
	edges := make([]*model.ReviewEdge, len(reviews))
	for i, r := range reviews {
		edges[i] = &model.ReviewEdge{
			Cursor: encodeCursor(r.ID),
			Node:   r,
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &model.ReviewConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     startIdx+len(reviews) < totalCount,
			HasPreviousPage: startIdx > 0,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
			TotalCount:      totalCount,
		},
	}
}

// CreateProduct creates a new product
func (db *Database) CreateProduct(input model.CreateProductInput) *model.Product {
	db.mu.Lock()
	defer db.mu.Unlock()

	product := &model.Product{
		ID:          "prod-" + uuid.New().String()[:8],
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		ImageURL:    input.ImageURL,
		Stock:       input.Stock,
		Rating:      0,
		ReviewCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	db.products[product.ID] = product

	if cat, ok := db.categories[input.CategoryID]; ok {
		cat.ProductCount++
	}

	return product
}

// UpdateProduct updates an existing product
func (db *Database) UpdateProduct(id string, input model.UpdateProductInput) *model.Product {
	db.mu.Lock()
	defer db.mu.Unlock()

	product, ok := db.products[id]
	if !ok {
		return nil
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}
	if input.Stock != nil {
		product.Stock = *input.Stock
	}
	product.UpdatedAt = time.Now()

	return product
}

// DeleteProduct deletes a product
func (db *Database) DeleteProduct(id string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.products[id]; !ok {
		return false
	}
	delete(db.products, id)
	return true
}

// CreateReview creates a new review
func (db *Database) CreateReview(input model.CreateReviewInput, userID, userName string) *model.Review {
	db.mu.Lock()
	defer db.mu.Unlock()

	review := &model.Review{
		ID:        "rev-" + uuid.New().String()[:8],
		ProductID: input.ProductID,
		UserID:    userID,
		UserName:  userName,
		Rating:    input.Rating,
		Comment:   input.Comment,
		CreatedAt: time.Now(),
	}

	db.reviews[review.ID] = review

	// Update product rating
	if product, ok := db.products[input.ProductID]; ok {
		product.ReviewCount++
		// Recalculate average rating
		var totalRating int
		var count int
		for _, r := range db.reviews {
			if r.ProductID == input.ProductID {
				totalRating += r.Rating
				count++
			}
		}
		if count > 0 {
			product.Rating = float64(totalRating) / float64(count)
		}
	}

	return review
}

// DeleteReview deletes a review
func (db *Database) DeleteReview(id string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.reviews[id]; !ok {
		return false
	}
	delete(db.reviews, id)
	return true
}

// GetSearchSuggestions returns autocomplete suggestions
func (db *Database) GetSearchSuggestions(query string, limit int) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 5
	}

	query = strings.ToLower(query)
	suggestions := make(map[string]bool)

	for _, p := range db.products {
		if strings.Contains(strings.ToLower(p.Name), query) {
			suggestions[p.Name] = true
		}
	}

	result := make([]string, 0, limit)
	for s := range suggestions {
		if len(result) >= limit {
			break
		}
		result = append(result, s)
	}

	return result
}

// Helper functions for cursor encoding
func encodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%s", id)))
}

func decodeCursor(cursor string) string {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(string(decoded), "cursor:")
}
