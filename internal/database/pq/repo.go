package pq

import (
	"fmt"

	"github.com/apartapatia/wall_of_comments/internal/entity"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func (p Repo) GetPosts() ([]*entity.Post, error) {
	var posts []*entity.Post
	if err := p.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (p Repo) CreatePost(post *entity.Post) (*entity.Post, error) {
	if err := p.db.Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (p Repo) GetPostById(id string) (*entity.Post, error) {
	post := &entity.Post{}
	if err := p.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (p Repo) CreateComment(comment *entity.Comment) (*entity.Comment, error) {
	if err := p.db.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (p Repo) GetCommentById(id string) (*entity.Comment, error) {
	var comment entity.Comment
	if err := p.db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (p Repo) GetCommentsForPost(postID string) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := p.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return []*entity.Comment{}, nil
	}
	return comments, nil
}

func (p Repo) GetCommentsForPostWithLimitAndOffset(postID string, limit *int, offset *int) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	query := p.db.Where("post_id = ?", postID)

	if limit != nil && offset != nil {
		query = query.Limit(*limit).Offset(*offset)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments for post with ID %s: %w", postID, err)
	}

	return comments, nil
}
