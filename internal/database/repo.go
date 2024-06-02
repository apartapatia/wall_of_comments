package database

import (
	"github.com/apartapatia/wall_of_comments/internal/entity"
)

type Repo interface {
	GetPosts() ([]*entity.Post, error)
	CreatePost(post *entity.Post) (*entity.Post, error)
	GetPostById(id string) (*entity.Post, error)
	CreateComment(comment *entity.Comment) (*entity.Comment, error)
	GetCommentById(id string) (*entity.Comment, error)
	GetCommentsForPost(postID string) ([]*entity.Comment, error)
	GetCommentsForPostWithLimitAndOffset(postID string, limit *int, offset *int) ([]*entity.Comment, error)
}
