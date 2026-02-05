package validation

import (
	"strings"
	"testing"

	"github.com/kalyankumar/graphql-product-search/graph/model"
)

func TestValidateCreateProduct(t *testing.T) {
	validInput := model.CreateProductInput{
		Name:        "Test Product",
		Description: "A test product description",
		Price:       29.99,
		CategoryID:  "cat-1",
		ImageURL:    "https://example.com/img.jpg",
		Stock:       10,
	}

	t.Run("valid input passes", func(t *testing.T) {
		if err := ValidateCreateProduct(validInput); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("empty name fails", func(t *testing.T) {
		input := validInput
		input.Name = ""
		err := ValidateCreateProduct(input)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("negative price fails", func(t *testing.T) {
		input := validInput
		input.Price = -1
		err := ValidateCreateProduct(input)
		if err == nil {
			t.Error("expected error for negative price")
		}
	})

	t.Run("price over max fails", func(t *testing.T) {
		input := validInput
		input.Price = 2_000_000
		err := ValidateCreateProduct(input)
		if err == nil {
			t.Error("expected error for price over max")
		}
	})

	t.Run("invalid image URL fails", func(t *testing.T) {
		input := validInput
		input.ImageURL = "ftp://example.com/img.jpg"
		err := ValidateCreateProduct(input)
		if err == nil {
			t.Error("expected error for non-HTTP image URL")
		}
	})

	t.Run("negative stock fails", func(t *testing.T) {
		input := validInput
		input.Stock = -1
		err := ValidateCreateProduct(input)
		if err == nil {
			t.Error("expected error for negative stock")
		}
	})

	t.Run("name too long fails", func(t *testing.T) {
		input := validInput
		input.Name = strings.Repeat("a", 201)
		err := ValidateCreateProduct(input)
		if err == nil {
			t.Error("expected error for name exceeding max length")
		}
	})
}

func TestValidateCreateReview(t *testing.T) {
	t.Run("valid input passes", func(t *testing.T) {
		input := model.CreateReviewInput{
			ProductID: "prod-1",
			Rating:    5,
			Comment:   "Great product!",
		}
		if err := ValidateCreateReview(input); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("rating 0 fails", func(t *testing.T) {
		input := model.CreateReviewInput{
			ProductID: "prod-1",
			Rating:    0,
			Comment:   "Comment",
		}
		if err := ValidateCreateReview(input); err == nil {
			t.Error("expected error for rating 0")
		}
	})

	t.Run("rating 6 fails", func(t *testing.T) {
		input := model.CreateReviewInput{
			ProductID: "prod-1",
			Rating:    6,
			Comment:   "Comment",
		}
		if err := ValidateCreateReview(input); err == nil {
			t.Error("expected error for rating 6")
		}
	})

	t.Run("empty comment fails", func(t *testing.T) {
		input := model.CreateReviewInput{
			ProductID: "prod-1",
			Rating:    3,
			Comment:   "",
		}
		if err := ValidateCreateReview(input); err == nil {
			t.Error("expected error for empty comment")
		}
	})

	t.Run("comment too long fails", func(t *testing.T) {
		input := model.CreateReviewInput{
			ProductID: "prod-1",
			Rating:    3,
			Comment:   strings.Repeat("x", 2001),
		}
		if err := ValidateCreateReview(input); err == nil {
			t.Error("expected error for comment exceeding max length")
		}
	})
}

func TestValidateID(t *testing.T) {
	t.Run("valid ID passes", func(t *testing.T) {
		if err := ValidateID("prod-1", "id"); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("empty ID fails", func(t *testing.T) {
		if err := ValidateID("", "id"); err == nil {
			t.Error("expected error for empty ID")
		}
	})

	t.Run("ID too long fails", func(t *testing.T) {
		if err := ValidateID(strings.Repeat("a", 51), "id"); err == nil {
			t.Error("expected error for ID exceeding max length")
		}
	})
}

func TestValidatePagination(t *testing.T) {
	t.Run("nil values pass", func(t *testing.T) {
		if err := ValidatePagination(nil, nil, 100); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("valid first passes", func(t *testing.T) {
		first := 10
		if err := ValidatePagination(&first, nil, 100); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("negative first fails", func(t *testing.T) {
		first := -1
		if err := ValidatePagination(&first, nil, 100); err == nil {
			t.Error("expected error for negative first")
		}
	})

	t.Run("first over max fails", func(t *testing.T) {
		first := 101
		if err := ValidatePagination(&first, nil, 100); err == nil {
			t.Error("expected error for first exceeding max page size")
		}
	})
}

func TestValidateSearchQuery(t *testing.T) {
	t.Run("valid query passes", func(t *testing.T) {
		if err := ValidateSearchQuery("macbook"); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("query too long fails", func(t *testing.T) {
		if err := ValidateSearchQuery(strings.Repeat("a", 101)); err == nil {
			t.Error("expected error for query exceeding max length")
		}
	})
}

func TestSanitizeCreateProduct(t *testing.T) {
	input := &model.CreateProductInput{
		Name:        "  <script>alert('xss')</script>  ",
		Description: "  Normal description  ",
		ImageURL:    "  https://example.com/img.jpg  ",
	}

	SanitizeCreateProduct(input)

	if strings.Contains(input.Name, "<script>") {
		t.Error("expected HTML to be escaped in name")
	}
	if input.Name != "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;" {
		t.Errorf("unexpected sanitized name: %s", input.Name)
	}
	if input.Description != "Normal description" {
		t.Errorf("expected trimmed description, got: '%s'", input.Description)
	}
	if input.ImageURL != "https://example.com/img.jpg" {
		t.Errorf("expected trimmed URL, got: '%s'", input.ImageURL)
	}
}
