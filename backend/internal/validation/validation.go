package validation

import (
	"fmt"
	"html"
	"strings"
	"unicode/utf8"

	"github.com/kalyankumar/graphql-product-search/graph/model"
)

// MaxLengths defines the maximum allowed lengths for various input fields.
var MaxLengths = map[string]int{
	"name":        200,
	"description": 5000,
	"comment":     2000,
	"search":      100,
	"imageURL":    500,
	"id":          50,
}

// ValidationError represents an input validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// sanitizeString removes leading/trailing whitespace and escapes HTML entities.
func sanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = html.EscapeString(s)
	return s
}

// ValidateCreateProduct validates product creation input.
func ValidateCreateProduct(input model.CreateProductInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return &ValidationError{Field: "name", Message: "product name is required"}
	}
	if utf8.RuneCountInString(input.Name) > MaxLengths["name"] {
		return &ValidationError{Field: "name", Message: fmt.Sprintf("product name must be %d characters or less", MaxLengths["name"])}
	}

	if strings.TrimSpace(input.Description) == "" {
		return &ValidationError{Field: "description", Message: "product description is required"}
	}
	if utf8.RuneCountInString(input.Description) > MaxLengths["description"] {
		return &ValidationError{Field: "description", Message: fmt.Sprintf("product description must be %d characters or less", MaxLengths["description"])}
	}

	if input.Price < 0 {
		return &ValidationError{Field: "price", Message: "price must be non-negative"}
	}
	if input.Price > 1_000_000 {
		return &ValidationError{Field: "price", Message: "price exceeds maximum allowed value"}
	}

	if strings.TrimSpace(input.CategoryID) == "" {
		return &ValidationError{Field: "categoryId", Message: "category ID is required"}
	}

	if strings.TrimSpace(input.ImageURL) == "" {
		return &ValidationError{Field: "imageUrl", Message: "image URL is required"}
	}
	if utf8.RuneCountInString(input.ImageURL) > MaxLengths["imageURL"] {
		return &ValidationError{Field: "imageUrl", Message: "image URL is too long"}
	}
	if !strings.HasPrefix(input.ImageURL, "http://") && !strings.HasPrefix(input.ImageURL, "https://") {
		return &ValidationError{Field: "imageUrl", Message: "image URL must use HTTP or HTTPS protocol"}
	}

	if input.Stock < 0 {
		return &ValidationError{Field: "stock", Message: "stock must be non-negative"}
	}
	if input.Stock > 1_000_000 {
		return &ValidationError{Field: "stock", Message: "stock exceeds maximum allowed value"}
	}

	return nil
}

// ValidateUpdateProduct validates product update input.
func ValidateUpdateProduct(input model.UpdateProductInput) error {
	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			return &ValidationError{Field: "name", Message: "product name cannot be empty"}
		}
		if utf8.RuneCountInString(*input.Name) > MaxLengths["name"] {
			return &ValidationError{Field: "name", Message: fmt.Sprintf("product name must be %d characters or less", MaxLengths["name"])}
		}
	}

	if input.Description != nil {
		if strings.TrimSpace(*input.Description) == "" {
			return &ValidationError{Field: "description", Message: "product description cannot be empty"}
		}
		if utf8.RuneCountInString(*input.Description) > MaxLengths["description"] {
			return &ValidationError{Field: "description", Message: fmt.Sprintf("product description must be %d characters or less", MaxLengths["description"])}
		}
	}

	if input.Price != nil {
		if *input.Price < 0 {
			return &ValidationError{Field: "price", Message: "price must be non-negative"}
		}
		if *input.Price > 1_000_000 {
			return &ValidationError{Field: "price", Message: "price exceeds maximum allowed value"}
		}
	}

	if input.ImageURL != nil {
		if strings.TrimSpace(*input.ImageURL) == "" {
			return &ValidationError{Field: "imageUrl", Message: "image URL cannot be empty"}
		}
		if !strings.HasPrefix(*input.ImageURL, "http://") && !strings.HasPrefix(*input.ImageURL, "https://") {
			return &ValidationError{Field: "imageUrl", Message: "image URL must use HTTP or HTTPS protocol"}
		}
	}

	if input.Stock != nil {
		if *input.Stock < 0 {
			return &ValidationError{Field: "stock", Message: "stock must be non-negative"}
		}
		if *input.Stock > 1_000_000 {
			return &ValidationError{Field: "stock", Message: "stock exceeds maximum allowed value"}
		}
	}

	return nil
}

// ValidateCreateReview validates review creation input.
func ValidateCreateReview(input model.CreateReviewInput) error {
	if strings.TrimSpace(input.ProductID) == "" {
		return &ValidationError{Field: "productId", Message: "product ID is required"}
	}

	if input.Rating < 1 || input.Rating > 5 {
		return &ValidationError{Field: "rating", Message: "rating must be between 1 and 5"}
	}

	if strings.TrimSpace(input.Comment) == "" {
		return &ValidationError{Field: "comment", Message: "comment is required"}
	}
	if utf8.RuneCountInString(input.Comment) > MaxLengths["comment"] {
		return &ValidationError{Field: "comment", Message: fmt.Sprintf("comment must be %d characters or less", MaxLengths["comment"])}
	}

	return nil
}

// ValidateID validates an entity ID.
func ValidateID(id, fieldName string) error {
	if strings.TrimSpace(id) == "" {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s is required", fieldName)}
	}
	if utf8.RuneCountInString(id) > MaxLengths["id"] {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s is too long", fieldName)}
	}
	return nil
}

// ValidateSearchQuery validates a search query string.
func ValidateSearchQuery(query string) error {
	if utf8.RuneCountInString(query) > MaxLengths["search"] {
		return &ValidationError{Field: "query", Message: fmt.Sprintf("search query must be %d characters or less", MaxLengths["search"])}
	}
	return nil
}

// ValidatePagination validates pagination parameters and enforces the max page size.
func ValidatePagination(first *int, last *int, maxPageSize int) error {
	if first != nil {
		if *first < 0 {
			return &ValidationError{Field: "first", Message: "first must be non-negative"}
		}
		if *first > maxPageSize {
			return &ValidationError{Field: "first", Message: fmt.Sprintf("first must not exceed %d", maxPageSize)}
		}
	}
	if last != nil {
		if *last < 0 {
			return &ValidationError{Field: "last", Message: "last must be non-negative"}
		}
		if *last > maxPageSize {
			return &ValidationError{Field: "last", Message: fmt.Sprintf("last must not exceed %d", maxPageSize)}
		}
	}
	return nil
}

// SanitizeCreateProduct sanitizes string fields in product creation input.
func SanitizeCreateProduct(input *model.CreateProductInput) {
	input.Name = sanitizeString(input.Name)
	input.Description = sanitizeString(input.Description)
	input.ImageURL = strings.TrimSpace(input.ImageURL)
}

// SanitizeCreateReview sanitizes string fields in review creation input.
func SanitizeCreateReview(input *model.CreateReviewInput) {
	input.Comment = sanitizeString(input.Comment)
}
