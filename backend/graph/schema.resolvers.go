package graph

import (
	"context"
	"fmt"

	"github.com/kalyankumar/graphql-product-search/graph/model"
	"github.com/kalyankumar/graphql-product-search/internal/auth"
	"github.com/kalyankumar/graphql-product-search/internal/database"
	"github.com/kalyankumar/graphql-product-search/internal/validation"
)

// --- Mutations ---

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	// Authorization: only admins can create products
	if !auth.RequireRole(ctx, auth.RoleAdmin) {
		return nil, fmt.Errorf("unauthorized: admin role required to create products")
	}

	// Validate & sanitize input
	validation.SanitizeCreateProduct(&input)
	if err := validation.ValidateCreateProduct(input); err != nil {
		return nil, err
	}

	// Verify category exists
	if cat := database.DB.GetCategory(input.CategoryID); cat == nil {
		return nil, fmt.Errorf("category '%s' not found", input.CategoryID)
	}

	return database.DB.CreateProduct(input), nil
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, input model.UpdateProductInput) (*model.Product, error) {
	// Authorization: only admins can update products
	if !auth.RequireRole(ctx, auth.RoleAdmin) {
		return nil, fmt.Errorf("unauthorized: admin role required to update products")
	}

	if err := validation.ValidateID(id, "id"); err != nil {
		return nil, err
	}
	if err := validation.ValidateUpdateProduct(input); err != nil {
		return nil, err
	}

	product := database.DB.UpdateProduct(id, input)
	if product == nil {
		return nil, fmt.Errorf("product '%s' not found", id)
	}
	return product, nil
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	// Authorization: only admins can delete products
	if !auth.RequireRole(ctx, auth.RoleAdmin) {
		return false, fmt.Errorf("unauthorized: admin role required to delete products")
	}

	if err := validation.ValidateID(id, "id"); err != nil {
		return false, err
	}

	if !database.DB.DeleteProduct(id) {
		return false, fmt.Errorf("product '%s' not found", id)
	}
	return true, nil
}

func (r *mutationResolver) CreateReview(ctx context.Context, input model.CreateReviewInput) (*model.Review, error) {
	// Authorization: authenticated users can create reviews
	user := auth.UserFromContext(ctx)
	if user.Role == auth.RoleAnonymous {
		return nil, fmt.Errorf("unauthorized: authentication required to create reviews")
	}

	// Validate & sanitize input
	validation.SanitizeCreateReview(&input)
	if err := validation.ValidateCreateReview(input); err != nil {
		return nil, err
	}

	// Verify product exists
	if p := database.DB.GetProduct(input.ProductID); p == nil {
		return nil, fmt.Errorf("product '%s' not found", input.ProductID)
	}

	return database.DB.CreateReview(input, user.ID, user.Name), nil
}

func (r *mutationResolver) DeleteReview(ctx context.Context, id string) (bool, error) {
	// Authorization: admins can delete any review
	if !auth.RequireRole(ctx, auth.RoleAdmin, auth.RoleUser) {
		return false, fmt.Errorf("unauthorized: authentication required to delete reviews")
	}

	if err := validation.ValidateID(id, "id"); err != nil {
		return false, err
	}

	if !database.DB.DeleteReview(id) {
		return false, fmt.Errorf("review '%s' not found", id)
	}
	return true, nil
}

// --- Queries ---

func (r *productResolver) Category(ctx context.Context, obj *model.Product) (*model.Category, error) {
	return database.DB.GetCategory(obj.CategoryID), nil
}

func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	if err := validation.ValidateID(id, "id"); err != nil {
		return nil, err
	}
	return database.DB.GetProduct(id), nil
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, last *int, before *string, filter *model.ProductFilterInput, sort *model.ProductSortInput) (*model.ProductConnection, error) {
	// Validate pagination
	maxPageSize := 100
	if r.Config != nil {
		maxPageSize = r.Config.MaxPageSize
	}
	if err := validation.ValidatePagination(first, last, maxPageSize); err != nil {
		return nil, err
	}

	// Validate search query length if present
	if filter != nil && filter.Search != nil {
		if err := validation.ValidateSearchQuery(*filter.Search); err != nil {
			return nil, err
		}
	}

	return database.DB.GetProducts(filter, sort, first, after, last, before), nil
}

func (r *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
	return database.DB.GetCategories(), nil
}

func (r *queryResolver) Category(ctx context.Context, id string) (*model.Category, error) {
	if err := validation.ValidateID(id, "id"); err != nil {
		return nil, err
	}
	return database.DB.GetCategory(id), nil
}

func (r *queryResolver) Reviews(ctx context.Context, productID string, first *int, after *string) (*model.ReviewConnection, error) {
	if err := validation.ValidateID(productID, "productId"); err != nil {
		return nil, err
	}
	maxPageSize := 100
	if r.Config != nil {
		maxPageSize = r.Config.MaxPageSize
	}
	if err := validation.ValidatePagination(first, nil, maxPageSize); err != nil {
		return nil, err
	}
	return database.DB.GetReviews(productID, first, after), nil
}

func (r *queryResolver) SearchSuggestions(ctx context.Context, query string, limit *int) ([]string, error) {
	if err := validation.ValidateSearchQuery(query); err != nil {
		return nil, err
	}

	l := 5
	if limit != nil {
		l = *limit
		if l > 20 {
			l = 20 // Cap suggestions
		}
		if l < 1 {
			l = 1
		}
	}
	return database.DB.GetSearchSuggestions(query, l), nil
}

// --- Subscriptions ---

func (r *subscriptionResolver) ProductUpdated(ctx context.Context, id string) (<-chan *model.Product, error) {
	ch := make(chan *model.Product)
	return ch, nil
}

func (r *subscriptionResolver) ReviewAdded(ctx context.Context, productID string) (<-chan *model.Review, error) {
	ch := make(chan *model.Review)
	return ch, nil
}

// --- Resolver types ---

func (r *Resolver) Mutation() MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Product() ProductResolver           { return &productResolver{r} }
func (r *Resolver) Query() QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
