package graph

import (
	"context"
	"errors"
	"github.com/KilyakArata/gqlgen-todos/models"
	"github.com/KilyakArata/gqlgen-todos/postgres"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostsRepo    postgres.PostsRepo
	CommentsRepo postgres.CommentsRepo
}

func (m *Resolver) DeletePostByID(ctx context.Context, id string) (bool, error) {
	var post models.Post
	post.ID = id
	list, err := m.CommentsRepo.GetCommentsForPost(ctx, &post)
	if err != nil {
		return false, err
	}
	for _, comment := range list {
		_, err = m.CommentsRepo.DeleteComment(ctx, comment.ID)
		if err != nil {
			return false, err
		}
	}
	_, err = m.PostsRepo.DB.WithContext(ctx).Model(&post).Where("id = ?", id).Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *Resolver) CreatComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	post, err := m.PostsRepo.GetPostByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	if post.Open != true {
		return nil, errors.New("post is closed")
	}
	_, err = m.CommentsRepo.DB.WithContext(ctx).Model(comment).Returning("*").Insert()

	notifySubscribers(comment.PostId, comment)

	return comment, err
}

func (m *Resolver) UpdatePostByID(ctx context.Context, post *models.Post) (*models.Post, error) {
	_, err := m.PostsRepo.DB.WithContext(ctx).Model(post).Where("id = ?", post.ID).Update()

	if post.Open != true {
		list, err := m.CommentsRepo.GetCommentsForPost(ctx, post)
		if err != nil {
			return nil, err
		}
		for _, comment := range list {
			_, err = m.CommentsRepo.DeleteComment(ctx, comment.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return post, err
}
