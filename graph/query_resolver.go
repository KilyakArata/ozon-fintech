package graph

import (
	"context"
	"github.com/KilyakArata/gqlgen-todos/graph/model"
	"github.com/KilyakArata/gqlgen-todos/models"
)

type queryResolver struct{ *Resolver }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, filter *model.CommentFilter, limit *int, offset *int) ([]*models.Comment, error) {
	return r.CommentsRepo.GetComments(ctx, filter, limit, offset)
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*models.Post, error) {
	return r.PostsRepo.GetPosts(ctx)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	return r.PostsRepo.GetPostByID(ctx, id)
}
