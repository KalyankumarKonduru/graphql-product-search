package graph

import (
	"context"

	"github.com/kalyankumar/graphql-product-search/graph/model"
	"github.com/kalyankumar/graphql-product-search/internal/database"
)

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
	return database.DB.CreateReview(input, "user-demo", "Demo User"), nil
}

func (r *mutationResolver) DeleteReview(ctx context.Context, id string) (bool, error) {
	return database.DB.DeleteReview(id), nil
}

func (r *productResolver) Category(ctx context.Context, obj *model.Product) (*model.Category, error) {
	return database.DB.GetCategory(obj.CategoryID), nil
}

func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	return database.DB.GetProduct(id), nil
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, last *int, before *string, filter *model.ProductFilterInput, sort *model.ProductSortInput) (*model.ProductConnection, error) {
	return database.DB.GetProducts(filter, sort, first, after, last, before), nil
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

func (r *subscriptionResolver) ProductUpdated(ctx context.Context, id string) (<-chan *model.Product, error) {
	ch := make(chan *model.Product)
	return ch, nil
}

func (r *subscriptionResolver) ReviewAdded(ctx context.Context, productID string) (<-chan *model.Review, error) {
	ch := make(chan *model.Review)
	return ch, nil
}

func (r *Resolver) Mutation() MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Product() ProductResolver           { return &productResolver{r} }
func (r *Resolver) Query() QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
