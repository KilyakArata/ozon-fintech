package postgres

import (
	"context"
	"github.com/KilyakArata/gqlgen-todos/models"
	"github.com/go-pg/pg/v10"
)

type PostsRepo struct {
	DB *pg.DB
}

func (m *PostsRepo) GetPosts(ctx context.Context) ([]*models.Post, error) {
	var posts []*models.Post
	err := m.DB.WithContext(ctx).Model(&posts).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostsRepo) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	err := m.DB.WithContext(ctx).Model(&post).Where("id = ?", id).First()
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (m *PostsRepo) CreatPost(ctx context.Context, comment *models.Post) (*models.Post, error) {
	_, err := m.DB.WithContext(ctx).Model(comment).Returning("*").Insert()

	return comment, err
}
