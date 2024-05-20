package graph

import (
	"context"
	"errors"
	"github.com/KilyakArata/gqlgen-todos/graph/model"
	"github.com/KilyakArata/gqlgen-todos/models"
)

type mutationResolver struct{ *Resolver }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	return r.DeletePostByID(ctx, id)
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input model.UpdatePost) (*models.Post, error) {
	post := &models.Post{
		ID:   id,
		Body: *input.Body,
		Open: *input.Open,
	}
	return r.UpdatePostByID(ctx, post)
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*models.Post, error) {
	if len(input.Author) < 3 {
		return nil, errors.New("name not long enough")
	}

	post := &models.Post{
		Author: input.Author,
		Body:   input.Body,
		Open:   input.Open,
	}

	return r.PostsRepo.CreatPost(ctx, post)
}

// UpdateComment is the resolver for the updateComment field.
func (r *mutationResolver) UpdateComment(ctx context.Context, id string, input model.UpdateComment) (*models.Comment, error) {
	comment := &models.Comment{
		ID:   id,
		Body: *input.Body,
	}
	return r.CommentsRepo.Update(ctx, comment)
}

// DeleteComment is the resolver for the deleteComment field.
func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	return r.CommentsRepo.DeleteComment(ctx, id)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*models.Comment, error) {
	if len(input.Author) < 3 {
		return nil, errors.New("name not long enough")
	}

	comment := &models.Comment{
		Author:          input.Author,
		Body:            input.Body,
		PostId:          input.ParentPost,
		ParentCommentId: *input.ParentComment,
	}

	return r.CreatComment(ctx, comment)
}
