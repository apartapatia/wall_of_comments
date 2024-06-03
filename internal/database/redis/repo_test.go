package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() (*Repo, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	validate := validator.New()

	repo := &Repo{
		db:       client,
		validate: validate,
	}

	return repo, s
}

func TestRepo_GetPosts(t *testing.T) {
	repo, s := setupTestDB()
	defer s.Close()

	s.HSet("post:1", "id", "1", "title", "Post 1", "content", "Content 1",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	s.HSet("post:2", "id", "2", "title", "Post 2", "content", "Content 2",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	s.HSet("post:3", "id", "3", "title", "Post 3", "content", "Content 3",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	s.HSet("post:4", "id", "4", "title", "Post 4", "content", "Content 4",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	posts, err := repo.GetPosts()

	assert.NoError(t, err)
	assert.Len(t, posts, 4)
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, "Content 1", posts[0].Content)

	assert.Equal(t, "Post 4", posts[3].Title)
	assert.Equal(t, "Content 4", posts[3].Content)
}

func TestRepo_GetPostByID(t *testing.T) {
	repo, s := setupTestDB()
	defer s.Close()

	s.HSet("post:1", "id", "1", "title", "Post 1", "content", "Content 1",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	post, err := repo.GetPostById("1")

	assert.NoError(t, err)
	assert.Equal(t, "Post 1", post.Title)
	assert.Equal(t, "Content 1", post.Content)
}

func TestRepo_GetCommentsForPost(t *testing.T) {
	repo, s := setupTestDB()
	defer s.Close()

	s.HSet("post:1", "id", "1", "title", "Post 1", "content", "Content 1",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	s.HSet("comment:1", "id", "1", "postId", "1", "content", "Content comment 1",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:2", "id", "2", "postId", "1", "content", "Content comment 2",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")

	comments, err := repo.GetCommentsForPost("1")

	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, "Content comment 1", comments[0].Content)
	assert.Equal(t, "Content comment 2", comments[1].Content)
}

func TestRepo_GetCommentsForPostWithLimitAndOffset(t *testing.T) {
	repo, s := setupTestDB()
	defer s.Close()

	s.HSet("post:1", "id", "1", "title", "Post 1", "content", "Content 1",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	s.HSet("comment:1", "id", "1", "postId", "1", "content", "Content comment 1",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:2", "id", "2", "postId", "1", "content", "Content comment 2",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:3", "id", "3", "postId", "1", "content", "Content comment 3",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:4", "id", "4", "postId", "1", "content", "Content comment 4",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:5", "id", "5", "postId", "1", "content", "Content comment 5",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:6", "id", "6", "postId", "1", "content", "Content comment 6",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:7", "id", "7", "postId", "1", "content", "Content comment 7",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:8", "id", "8", "postId", "1", "content", "Content comment 8",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")

	limit := 5
	offset := 1
	comments, err := repo.GetCommentsForPostWithLimitAndOffset("1", &limit, &offset)
	assert.NoError(t, err)
	assert.Len(t, comments, limit)

	for i, comment := range comments {
		assert.Equal(t, fmt.Sprintf("Content comment %d", offset+i+1), comment.Content)
	}
}

func TestRepo_GetCommentById(t *testing.T) {
	repo, s := setupTestDB()
	defer s.Close()

	s.HSet("post:1", "id", "1", "title", "Post 1", "content", "Content 1",
		"commentsActive", "1", "createdAt", time.Now().Format(time.RFC3339),
		"updatedAt", time.Now().Format(time.RFC3339), "comments", "[]")

	s.HSet("comment:1", "id", "1", "postId", "1", "content", "Content comment 1",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:2", "id", "2", "postId", "1", "content", "Content comment 2",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:3", "id", "3", "postId", "1", "content", "Content comment 3",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:4", "id", "4", "postId", "1", "content", "Content comment 4",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:5", "id", "5", "postId", "1", "content", "Content comment 5",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:6", "id", "6", "postId", "1", "content", "Content comment 6",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:7", "id", "7", "postId", "1", "content", "Content comment 7",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")
	s.HSet("comment:8", "id", "8", "postId", "1", "content", "Content comment 8",
		"createdAt", time.Now().Format(time.RFC3339), "updatedAt", time.Now().Format(time.RFC3339), "replies", "[]", "parentId", "")

	comment, err := repo.GetCommentById("4")
	assert.NoError(t, err)
	assert.Equal(t, comment.Content, "Content comment 4")
}
