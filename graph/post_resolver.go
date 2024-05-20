package graph

import (
	"context"
	"github.com/KilyakArata/gqlgen-todos/models"
)

type postResolver struct{ *Resolver }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *models.Post) ([]*models.Comment, error) {
	return r.CommentsRepo.GetCommentsForPost(ctx, obj)
}
