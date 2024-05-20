package postgres

import (
	"context"
	"fmt"
	"github.com/KilyakArata/gqlgen-todos/graph/model"
	"github.com/KilyakArata/gqlgen-todos/models"
	"github.com/go-pg/pg/v10"
)

type CommentsRepo struct {
	DB *pg.DB
}

func (m *CommentsRepo) GetComments(ctx context.Context, filter *model.CommentFilter, limit, offset *int) ([]*models.Comment, error) {
	var comments []*models.Comment

	query := m.DB.WithContext(ctx).Model(&comments).Order("id")

	if filter != nil {
		if filter.Author != nil && *filter.Author != "" {
			query.Where("author ILIKE ?", fmt.Sprintf("%%s%%", *filter.Author))
		}
		if filter.Body != nil && *filter.Body != "" {
			query.Where("body ILIKE ?", fmt.Sprintf("%%s%%", *filter.Body))
		}
	}

	if limit != nil {
		query.Limit(*limit)
	}

	if offset != nil {
		query.Offset(*offset)
	}

	err := query.Select()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentsRepo) GetCommentsForPost(ctx context.Context, post *models.Post) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := m.DB.WithContext(ctx).
		Model(&comments).
		Where("post_id = ?", post.ID).
		Where("parent_comment_id IS NULL").
		Order("id").
		Select()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentsRepo) GetCommentsForComment(ctx context.Context, comment *models.Comment) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := m.DB.WithContext(ctx).Model(&comments).Where("parent_comment_id = ?", comment.ID).Order("id").Select()
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentsRepo) DeleteComment(ctx context.Context, id string) (bool, error) {
	var comment models.Comment
	comment.ID = id
	list, err := m.GetCommentsForComment(ctx, &comment)
	if err != nil {
		return false, err
	}
	for _, comment := range list {
		_, err = m.DeleteComment(ctx, comment.ID)
		if err != nil {
			return false, err
		}
	}
	_, err = m.DB.WithContext(ctx).Model(&comment).Where("id = ?", id).Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *CommentsRepo) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	_, err := m.DB.WithContext(ctx).Model(comment).Where("id = ?", comment.ID).Update()
	return comment, err
}
