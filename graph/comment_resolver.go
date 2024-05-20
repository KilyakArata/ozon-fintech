package graph

import (
	"context"
	"github.com/KilyakArata/gqlgen-todos/models"
)

type commentResolver struct{ *Resolver }

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Comments is the resolver for the comments field.
func (r *commentResolver) Comments(ctx context.Context, obj *models.Comment) ([]*models.Comment, error) {
	return r.CommentsRepo.GetCommentsForComment(ctx, obj)
}
