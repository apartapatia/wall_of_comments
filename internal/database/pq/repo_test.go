package pq

import (
	"fmt"
	"testing"

	"github.com/apartapatia/wall_of_comments/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&entity.Post{}, &entity.Comment{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestRepo_GetPosts(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo{db: db}

	posts := []*entity.Post{
		{ID: "1", Title: "Post 1", Content: "Content 1"},
		{ID: "2", Title: "Post 2", Content: "Content 2"},
		{ID: "3", Title: "Post 3", Content: "Content 3"},
		{ID: "4", Title: "Post 4", Content: "Content 4"},
	}
	for _, post := range posts {
		_, err := repo.CreatePost(post)
		assert.NoError(t, err)
	}

	retPosts, err := repo.GetPosts()
	assert.NoError(t, err)

	assert.Equal(t, len(posts), len(retPosts))
	for i := range posts {
		assert.Equal(t, posts[i].Title, retPosts[i].Title)
	}
}

func TestRepo_GetPostById(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo{db: db}

	post := &entity.Post{Title: "Post 1", Content: "Content 1"}
	_, err := repo.CreatePost(post)
	assert.NoError(t, err)

	retPost, err := repo.GetPostById(post.ID)
	assert.NoError(t, err)

	assert.NotNil(t, retPost)
	assert.Equal(t, post.Title, retPost.Title)
	assert.Equal(t, post.Content, retPost.Content)
}

func TestRepo_GetCommentsForPost(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo{db: db}

	post := &entity.Post{ID: "1", Title: "Post 1", Content: "Content 1"}
	_, err := repo.CreatePost(post)
	assert.NoError(t, err)

	comments := &entity.Comment{PostID: post.ID, Content: "Content comment 1"}
	_, err = repo.CreateComment(comments)
	assert.NoError(t, err)

	getComment, err := repo.GetCommentsForPost(post.ID)
	assert.NotNil(t, getComment)
	assert.Equal(t, comments.Content, getComment[0].Content)
}

func TestRepo_GetCommentsForPostWithLimitAndOffset(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo{db: db}

	post := &entity.Post{ID: "1", Title: "Post 1", Content: "Content 1"}
	_, err := repo.CreatePost(post)
	assert.NoError(t, err)

	for i := range 100 {
		comment := &entity.Comment{PostID: post.ID, Content: fmt.Sprintf("Content comment %d", i)}
		_, err := repo.CreateComment(comment)
		assert.NoError(t, err)
	}

	limit := 10
	offset := 20
	comments, err := repo.GetCommentsForPostWithLimitAndOffset(post.ID, &limit, &offset)
	assert.NoError(t, err)
	assert.Len(t, comments, limit)

	for i, comment := range comments {
		assert.Equal(t, fmt.Sprintf("Content comment %d", offset+i), comment.Content)
	}
}

func TestRepo_GetCommentById(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo{db: db}

	post := &entity.Post{Title: "Post 1", Content: "Content 1"}
	_, err := repo.CreatePost(post)
	assert.NoError(t, err)

	comment := &entity.Comment{ID: "1", PostID: post.ID, Content: "Content comment 11"}
	_, err = repo.CreateComment(comment)
	assert.NoError(t, err)

	getComment, err := repo.GetCommentById(comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, getComment)
	assert.Equal(t, comment.Content, getComment.Content)
}
