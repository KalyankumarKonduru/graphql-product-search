package graph

import (
	"context"

	"github.com/kalyankumar/graphql-product-search/graph/model"
	"github.com/kalyankumar/graphql-product-search/internal/database"
)

// Resolver is the root resolver
type Resolver struct{}

// Query resolvers
type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	return database.DB.GetProduct(id), nil
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, last *int, before *string, filter *model.ProductFilterInput, sortInput *model.ProductSortInput) (*model.ProductConnection, error) {
	return database.DB.GetProducts(filter, sortInput, first, after, last, before), nil
}

func (r *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
	return database.DB.GetCategories(), nil
}

func (r *queryResolver) Category(ctx context.Context, id string) (*model.Category, error) {
	return database.DB.GetCategory(id), nil
}

func (r *queryResolver) Reviews(ctx context.Context, productID string, first *int, after *string) (*model.ReviewConnection, error) {
	return database.DB.GetReviews(productID, first, after), nil
}

func (r *queryResolver) SearchSuggestions(ctx context.Context, query string, limit *int) ([]string, error) {
	l := 5
	if limit != nil {
		l = *limit
	}
	return database.DB.GetSearchSuggestions(query, l), nil
}

// Mutation resolvers
type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	return database.DB.CreateProduct(input), nil
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, input model.UpdateProductInput) (*model.Product, error) {
	return database.DB.UpdateProduct(id, input), nil
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	return database.DB.DeleteProduct(id), nil
}

func (r *mutationResolver) CreateReview(ctx context.Context, input model.CreateReviewInput) (*model.Review, error) {
	// In a real app, get user from context/auth
	return database.DB.CreateReview(input, "user-demo", "Demo User"), nil
}

func (r *mutationResolver) DeleteReview(ctx context.Context, id string) (bool, error) {
	return database.DB.DeleteReview(id), nil
}

// Product field resolvers
type productResolver struct{ *Resolver }

func (r *Resolver) Product() ProductResolver {
	return &productResolver{r}
}

func (r *productResolver) Category(ctx context.Context, obj *model.Product) (*model.Category, error) {
	return database.DB.GetCategory(obj.CategoryID), nil
}

// Subscription resolvers (placeholder for demo)
type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

func (r *subscriptionResolver) ProductUpdated(ctx context.Context, id string) (<-chan *model.Product, error) {
	// Placeholder - would use websockets/pubsub in production
	ch := make(chan *model.Product)
	return ch, nil
}

func (r *subscriptionResolver) ReviewAdded(ctx context.Context, productID string) (<-chan *model.Review, error) {
	// Placeholder - would use websockets/pubsub in production
	ch := make(chan *model.Review)
	return ch, nil
}

// Interface definitions for gqlgen
type QueryResolver interface {
	Product(ctx context.Context, id string) (*model.Product, error)
	Products(ctx context.Context, first *int, after *string, last *int, before *string, filter *model.ProductFilterInput, sort *model.ProductSortInput) (*model.ProductConnection, error)
	Categories(ctx context.Context) ([]*model.Category, error)
	Category(ctx context.Context, id string) (*model.Category, error)
	Reviews(ctx context.Context, productID string, first *int, after *string) (*model.ReviewConnection, error)
	SearchSuggestions(ctx context.Context, query string, limit *int) ([]string, error)
}

type MutationResolver interface {
	CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, id string, input model.UpdateProductInput) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) (bool, error)
	CreateReview(ctx context.Context, input model.CreateReviewInput) (*model.Review, error)
	DeleteReview(ctx context.Context, id string) (bool, error)
}

type SubscriptionResolver interface {
	ProductUpdated(ctx context.Context, id string) (<-chan *model.Product, error)
	ReviewAdded(ctx context.Context, productID string) (<-chan *model.Review, error)
}

type ProductResolver interface {
	Category(ctx context.Context, obj *model.Product) (*model.Category, error)
}
