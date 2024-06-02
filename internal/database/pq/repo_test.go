package pq

import (
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
		{Title: "Post 1", Content: "Content 1"},
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

	post := &entity.Post{Title: "Test Post", Content: "Test Content"}
	_, err := repo.CreatePost(post)
	assert.NoError(t, err)

	retPost, err := repo.GetPostById(post.ID)
	assert.NoError(t, err)

	assert.NotNil(t, retPost)
	assert.Equal(t, post.Title, retPost.Title)
	assert.Equal(t, post.Content, retPost.Content)
}
