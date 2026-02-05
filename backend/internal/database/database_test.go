package database

import (
	"testing"

	"github.com/kalyankumar/graphql-product-search/graph/model"
)

func setupTestDB() {
	Initialize()
}

// --- Product CRUD Tests ---

func TestGetProduct(t *testing.T) {
	setupTestDB()

	t.Run("returns existing product", func(t *testing.T) {
		product := DB.GetProduct("prod-1")
		if product == nil {
			t.Fatal("expected product, got nil")
		}
		if product.Name != "MacBook Pro 16\"" {
			t.Errorf("expected 'MacBook Pro 16\"', got '%s'", product.Name)
		}
	})

	t.Run("returns nil for non-existent product", func(t *testing.T) {
		product := DB.GetProduct("prod-nonexistent")
		if product != nil {
			t.Error("expected nil for non-existent product")
		}
	})
}

func TestCreateProduct(t *testing.T) {
	setupTestDB()

	input := model.CreateProductInput{
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		CategoryID:  "cat-1",
		ImageURL:    "https://example.com/img.jpg",
		Stock:       10,
	}

	product := DB.CreateProduct(input)

	if product == nil {
		t.Fatal("expected product, got nil")
	}
	if product.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got '%s'", product.Name)
	}
	if product.Price != 29.99 {
		t.Errorf("expected price 29.99, got %f", product.Price)
	}
	if product.Rating != 0 {
		t.Errorf("expected rating 0, got %f", product.Rating)
	}
	if product.ReviewCount != 0 {
		t.Errorf("expected review count 0, got %d", product.ReviewCount)
	}

	// Verify it's retrievable
	retrieved := DB.GetProduct(product.ID)
	if retrieved == nil {
		t.Fatal("expected to retrieve created product")
	}
}

func TestUpdateProduct(t *testing.T) {
	setupTestDB()

	t.Run("updates existing product", func(t *testing.T) {
		newName := "Updated MacBook"
		newPrice := 2299.00

		updated := DB.UpdateProduct("prod-1", model.UpdateProductInput{
			Name:  &newName,
			Price: &newPrice,
		})

		if updated == nil {
			t.Fatal("expected updated product, got nil")
		}
		if updated.Name != "Updated MacBook" {
			t.Errorf("expected 'Updated MacBook', got '%s'", updated.Name)
		}
		if updated.Price != 2299.00 {
			t.Errorf("expected price 2299.00, got %f", updated.Price)
		}
	})

	t.Run("returns nil for non-existent product", func(t *testing.T) {
		newName := "Test"
		updated := DB.UpdateProduct("prod-nonexistent", model.UpdateProductInput{
			Name: &newName,
		})
		if updated != nil {
			t.Error("expected nil for non-existent product")
		}
	})
}

func TestDeleteProduct(t *testing.T) {
	setupTestDB()

	t.Run("deletes existing product", func(t *testing.T) {
		result := DB.DeleteProduct("prod-1")
		if !result {
			t.Error("expected true for deleting existing product")
		}

		product := DB.GetProduct("prod-1")
		if product != nil {
			t.Error("expected nil after deletion")
		}
	})

	t.Run("returns false for non-existent product", func(t *testing.T) {
		result := DB.DeleteProduct("prod-nonexistent")
		if result {
			t.Error("expected false for deleting non-existent product")
		}
	})
}

// --- Filter & Sort Tests ---

func TestGetProductsWithFilter(t *testing.T) {
	setupTestDB()

	t.Run("filters by search term", func(t *testing.T) {
		search := "MacBook"
		conn := DB.GetProducts(&model.ProductFilterInput{Search: &search}, nil, nil, nil, nil, nil)
		if len(conn.Edges) == 0 {
			t.Error("expected at least one result for 'MacBook' search")
		}
		for _, edge := range conn.Edges {
			if edge.Node.Name != "MacBook Pro 16\"" {
				t.Errorf("unexpected product in search results: %s", edge.Node.Name)
			}
		}
	})

	t.Run("filters by category", func(t *testing.T) {
		catID := "cat-1" // Electronics
		conn := DB.GetProducts(&model.ProductFilterInput{CategoryID: &catID}, nil, nil, nil, nil, nil)
		for _, edge := range conn.Edges {
			if edge.Node.CategoryID != "cat-1" {
				t.Errorf("expected category cat-1, got %s", edge.Node.CategoryID)
			}
		}
	})

	t.Run("filters by price range", func(t *testing.T) {
		minPrice := 100.0
		maxPrice := 500.0
		conn := DB.GetProducts(&model.ProductFilterInput{
			MinPrice: &minPrice,
			MaxPrice: &maxPrice,
		}, nil, nil, nil, nil, nil)

		for _, edge := range conn.Edges {
			if edge.Node.Price < minPrice || edge.Node.Price > maxPrice {
				t.Errorf("product %s price %f outside range [%f, %f]",
					edge.Node.Name, edge.Node.Price, minPrice, maxPrice)
			}
		}
	})

	t.Run("filters by in stock", func(t *testing.T) {
		inStock := true
		conn := DB.GetProducts(&model.ProductFilterInput{InStock: &inStock}, nil, nil, nil, nil, nil)
		for _, edge := range conn.Edges {
			if edge.Node.Stock <= 0 {
				t.Errorf("product %s has stock %d but filter requires in stock",
					edge.Node.Name, edge.Node.Stock)
			}
		}
	})
}

func TestGetProductsWithSort(t *testing.T) {
	setupTestDB()

	t.Run("sorts by price ascending", func(t *testing.T) {
		conn := DB.GetProducts(nil, &model.ProductSortInput{
			Field: model.ProductSortFieldPrice,
			Order: model.SortOrderAsc,
		}, nil, nil, nil, nil)

		for i := 1; i < len(conn.Edges); i++ {
			if conn.Edges[i].Node.Price < conn.Edges[i-1].Node.Price {
				t.Errorf("products not sorted by price ascending: %f < %f",
					conn.Edges[i].Node.Price, conn.Edges[i-1].Node.Price)
			}
		}
	})

	t.Run("sorts by price descending", func(t *testing.T) {
		conn := DB.GetProducts(nil, &model.ProductSortInput{
			Field: model.ProductSortFieldPrice,
			Order: model.SortOrderDesc,
		}, nil, nil, nil, nil)

		for i := 1; i < len(conn.Edges); i++ {
			if conn.Edges[i].Node.Price > conn.Edges[i-1].Node.Price {
				t.Errorf("products not sorted by price descending: %f > %f",
					conn.Edges[i].Node.Price, conn.Edges[i-1].Node.Price)
			}
		}
	})

	t.Run("sorts by name ascending", func(t *testing.T) {
		conn := DB.GetProducts(nil, &model.ProductSortInput{
			Field: model.ProductSortFieldName,
			Order: model.SortOrderAsc,
		}, nil, nil, nil, nil)

		for i := 1; i < len(conn.Edges); i++ {
			if conn.Edges[i].Node.Name < conn.Edges[i-1].Node.Name {
				t.Errorf("products not sorted by name ascending: %s < %s",
					conn.Edges[i].Node.Name, conn.Edges[i-1].Node.Name)
			}
		}
	})
}

// --- Pagination Tests ---

func TestPagination(t *testing.T) {
	setupTestDB()

	t.Run("first N items", func(t *testing.T) {
		first := 3
		conn := DB.GetProducts(nil, nil, &first, nil, nil, nil)

		if len(conn.Edges) != 3 {
			t.Errorf("expected 3 edges, got %d", len(conn.Edges))
		}
		if !conn.PageInfo.HasNextPage {
			t.Error("expected HasNextPage to be true")
		}
	})

	t.Run("cursor-based pagination", func(t *testing.T) {
		first := 3
		page1 := DB.GetProducts(nil, nil, &first, nil, nil, nil)

		if len(page1.Edges) != 3 {
			t.Fatalf("expected 3 edges in page 1, got %d", len(page1.Edges))
		}

		// Get page 2 using the cursor from page 1
		cursor := page1.PageInfo.EndCursor
		page2 := DB.GetProducts(nil, nil, &first, cursor, nil, nil)

		if len(page2.Edges) == 0 {
			t.Error("expected edges in page 2")
		}

		// Ensure no overlap between pages
		page1IDs := make(map[string]bool)
		for _, e := range page1.Edges {
			page1IDs[e.Node.ID] = true
		}
		for _, e := range page2.Edges {
			if page1IDs[e.Node.ID] {
				t.Errorf("product %s appears in both page 1 and page 2", e.Node.ID)
			}
		}
	})

	t.Run("total count is correct", func(t *testing.T) {
		first := 2
		conn := DB.GetProducts(nil, nil, &first, nil, nil, nil)
		if conn.PageInfo.TotalCount != 10 {
			t.Errorf("expected total count 10, got %d", conn.PageInfo.TotalCount)
		}
	})
}

// --- Category Tests ---

func TestGetCategories(t *testing.T) {
	setupTestDB()

	categories := DB.GetCategories()
	if len(categories) != 5 {
		t.Errorf("expected 5 categories, got %d", len(categories))
	}
}

func TestGetCategory(t *testing.T) {
	setupTestDB()

	t.Run("returns existing category", func(t *testing.T) {
		cat := DB.GetCategory("cat-1")
		if cat == nil {
			t.Fatal("expected category, got nil")
		}
		if cat.Name != "Electronics" {
			t.Errorf("expected 'Electronics', got '%s'", cat.Name)
		}
	})

	t.Run("returns nil for non-existent category", func(t *testing.T) {
		cat := DB.GetCategory("cat-nonexistent")
		if cat != nil {
			t.Error("expected nil for non-existent category")
		}
	})
}

// --- Review Tests ---

func TestCreateReview(t *testing.T) {
	setupTestDB()

	input := model.CreateReviewInput{
		ProductID: "prod-1",
		Rating:    5,
		Comment:   "Great product!",
	}

	review := DB.CreateReview(input, "user-test", "Test User")

	if review == nil {
		t.Fatal("expected review, got nil")
	}
	if review.Rating != 5 {
		t.Errorf("expected rating 5, got %d", review.Rating)
	}
	if review.UserName != "Test User" {
		t.Errorf("expected username 'Test User', got '%s'", review.UserName)
	}

	// Verify product rating was updated
	product := DB.GetProduct("prod-1")
	if product.ReviewCount <= 2 { // Should be more than the initial 2 reviews
		t.Errorf("expected review count > 2, got %d", product.ReviewCount)
	}
}

func TestGetReviews(t *testing.T) {
	setupTestDB()

	conn := DB.GetReviews("prod-1", nil, nil)
	if len(conn.Edges) == 0 {
		t.Error("expected reviews for prod-1")
	}

	// Verify all reviews belong to prod-1
	for _, edge := range conn.Edges {
		if edge.Node.ProductID != "prod-1" {
			t.Errorf("expected productID 'prod-1', got '%s'", edge.Node.ProductID)
		}
	}
}

func TestDeleteReview(t *testing.T) {
	setupTestDB()

	t.Run("deletes existing review", func(t *testing.T) {
		result := DB.DeleteReview("rev-1")
		if !result {
			t.Error("expected true for deleting existing review")
		}
	})

	t.Run("returns false for non-existent review", func(t *testing.T) {
		result := DB.DeleteReview("rev-nonexistent")
		if result {
			t.Error("expected false for deleting non-existent review")
		}
	})
}

// --- Search Suggestions Tests ---

func TestGetSearchSuggestions(t *testing.T) {
	setupTestDB()

	t.Run("returns matching suggestions", func(t *testing.T) {
		suggestions := DB.GetSearchSuggestions("mac", 5)
		if len(suggestions) == 0 {
			t.Error("expected suggestions for 'mac'")
		}
	})

	t.Run("respects limit", func(t *testing.T) {
		suggestions := DB.GetSearchSuggestions("", 2)
		if len(suggestions) > 2 {
			t.Errorf("expected at most 2 suggestions, got %d", len(suggestions))
		}
	})

	t.Run("returns empty for no matches", func(t *testing.T) {
		suggestions := DB.GetSearchSuggestions("zzzznonexistent", 5)
		if len(suggestions) != 0 {
			t.Errorf("expected 0 suggestions, got %d", len(suggestions))
		}
	})
}

// --- Concurrent Access Tests ---

func TestConcurrentAccess(t *testing.T) {
	setupTestDB()

	done := make(chan bool, 100)

	// Concurrent reads
	for i := 0; i < 50; i++ {
		go func() {
			DB.GetProducts(nil, nil, nil, nil, nil, nil)
			done <- true
		}()
	}

	// Concurrent writes
	for i := 0; i < 50; i++ {
		go func() {
			input := model.CreateProductInput{
				Name:        "Concurrent Product",
				Description: "Test concurrent creation",
				Price:       9.99,
				CategoryID:  "cat-1",
				ImageURL:    "https://example.com/img.jpg",
				Stock:       1,
			}
			DB.CreateProduct(input)
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
