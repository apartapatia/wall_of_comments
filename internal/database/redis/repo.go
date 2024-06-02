package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/apartapatia/wall_of_comments/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
)

var ErrNotActive = errors.New("post comments are not active")

type Repo struct {
	db       *redis.Client
	validate *validator.Validate
}

func (rp *Repo) GetPosts() ([]*entity.Post, error) {
	keys, err := rp.db.Keys("post:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get post keys from Redis: %w", err)
	}

	var posts []*entity.Post
	for _, key := range keys {
		data, err := rp.db.HGetAll(key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get post from Redis: %w", err)
		}
		post, err := mapToPost(data)
		if err != nil {
			return nil, fmt.Errorf("failed map to post: %w", err)
		}

		comments, err := rp.GetCommentsForPost(post.ID)
		if err != nil {
			return nil, fmt.Errorf("failed get comments for post: %w", err)
		}

		if comments == nil {
			comments = []*entity.Comment{}
		}

		post.Comments = comments

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].UpdatedAt.Before(posts[j].CreatedAt)
	})

	return posts, nil
}

func (rp *Repo) CreatePost(post *entity.Post) (*entity.Post, error) {
	if err := rp.validate.Struct(post); err != nil {
		return nil, fmt.Errorf("failed to validate post: %w", err)
	}

	redisID := fmt.Sprintf("post:%s", post.ID)
	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	data, err := postToMap(post)
	if err != nil {
		return nil, fmt.Errorf("failed map to post: %w", err)
	}

	_, err = rp.db.HMSet(redisID, data).Result()
	if err != nil {
		return nil, fmt.Errorf("failed set post to Redis: %w", err)
	}

	return post, nil
}

func (rp *Repo) GetPostById(id string) (*entity.Post, error) {
	key := fmt.Sprintf("post:%s", id)
	data, err := rp.db.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get post from Redis: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("post with id %s not found", id)
	}

	post, err := mapToPost(data)
	if err != nil {
		return nil, fmt.Errorf("failed to map post: %w", err)
	}

	comments, err := rp.GetCommentsForPost(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for post: %w", err)
	}

	if comments == nil {
		comments = []*entity.Comment{}
	}

	post.Comments = comments

	return post, nil
}

func (rp *Repo) GetCommentsForPost(postID string) ([]*entity.Comment, error) {
	keys, err := rp.db.Keys(fmt.Sprintf("comment:*")).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment keys from Redis: %w", err)
	}

	var comments []*entity.Comment
	for _, key := range keys {
		data, err := rp.db.HGetAll(key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get comment from Redis: %w", err)
		}
		comment, err := mapToComment(data)
		if err != nil {
			return nil, fmt.Errorf("failed map to comment: %w", err)
		}

		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (rp *Repo) CreateComment(comment *entity.Comment) (*entity.Comment, error) {
	if err := rp.validate.Struct(comment); err != nil {
		return nil, fmt.Errorf("failed to validate comment: %w", err)
	}

	post, err := rp.GetPostById(comment.PostID)
	if err != nil {
		return nil, err
	}

	if !post.CommentsActive {
		return nil, ErrNotActive
	}

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = comment.CreatedAt
	comment.PostID = post.ID
	redisID := fmt.Sprintf("comment:%s", comment.ID)

	data, err := commentToMap(comment)
	if err != nil {
		return nil, fmt.Errorf("failed comment to map: %w", err)
	}

	_, err = rp.db.HMSet(redisID, data).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to set comment to Redis: %w", err)
	}

	if comment.ParentID != nil {
		parentComment, err := rp.GetCommentById(*comment.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed get comment by id: %w", err)
		}
		parentComment.Replies = append(parentComment.Replies, comment)
		data, err := commentToMap(parentComment)
		if err != nil {
			return nil, fmt.Errorf("failed comment to map: %w", err)
		}
		_, err = rp.db.HMSet(parentComment.ID, data).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to set comment to Redis: %w", err)
		}
	}

	post.Comments = append(post.Comments, comment)
	data, err = postToMap(post)
	if err != nil {
		return nil, fmt.Errorf("failed post to map: %w", err)
	}
	_, err = rp.db.HMSet(post.ID, data).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to set post to Redis: %w", err)
	}

	return comment, nil
}

func (rp *Repo) GetCommentById(id string) (*entity.Comment, error) {
	key := fmt.Sprintf("comment:%s", id)
	data, err := rp.db.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment key from Redis: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("comment with id %s not found", id)
	}

	comment, err := mapToComment(data)
	if err != nil {
		return nil, fmt.Errorf("failed map to comment: %w", err)
	}

	return comment, nil
}

func (rp *Repo) GetCommentsForPostWithLimitAndOffset(postID string, limit *int, offset *int) ([]*entity.Comment, error) {
	keys, err := rp.db.Keys("comment:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment keys from Redis: %w", err)
	}

	var comments []*entity.Comment

	var startIndex, endIndex int

	if limit != nil && offset != nil {
		startIndex = *offset
		endIndex = *offset + *limit
		if endIndex > len(keys) {
			endIndex = len(keys)
		}
	} else {
		startIndex = 0
		endIndex = len(keys)
	}

	for _, key := range keys {
		data, err := rp.db.HGetAll(key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get comment from Redis: %w", err)
		}
		comment, err := mapToComment(data)
		if err != nil {
			return nil, fmt.Errorf("failed to map to comment: %w", err)
		}

		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	if startIndex >= len(comments) {
		return []*entity.Comment{}, nil
	}
	if endIndex > len(comments) {
		endIndex = len(comments)
	}
	comments = comments[startIndex:endIndex]

	return comments, nil
}

func postToMap(post *entity.Post) (map[string]interface{}, error) {
	comments, err := json.Marshal(post.Comments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comments: %w", err)
	}

	return map[string]interface{}{
		"id":             post.ID,
		"title":          post.Title,
		"content":        post.Content,
		"commentsActive": post.CommentsActive,
		"createdAt":      post.CreatedAt.Format(time.RFC3339),
		"updatedAt":      post.UpdatedAt.Format(time.RFC3339),
		"comments":       string(comments),
	}, nil
}

func commentToMap(comment *entity.Comment) (map[string]interface{}, error) {
	replies, err := json.Marshal(comment.Replies)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal replies: %w", err)
	}

	result := map[string]interface{}{
		"id":        comment.ID,
		"postId":    comment.PostID,
		"content":   comment.Content,
		"createdAt": comment.CreatedAt.Format(time.RFC3339),
		"updatedAt": comment.UpdatedAt.Format(time.RFC3339),
		"replies":   string(replies),
	}

	if comment.ParentID != nil {
		result["parentId"] = *comment.ParentID
	} else {
		result["parentId"] = ""
	}

	return result, nil
}

func mapToPost(data map[string]string) (*entity.Post, error) {
	createdAt, err := time.Parse(time.RFC3339, data["createdAt"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, data["updatedAt"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse updatedAt: %w", err)
	}

	var comments []*entity.Comment
	if err := json.Unmarshal([]byte(data["comments"]), &comments); err != nil {
		comments = []*entity.Comment{}
	}

	return &entity.Post{
		ID:             data["id"],
		Title:          data["title"],
		Content:        data["content"],
		CommentsActive: data["commentsActive"] == "1",
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Comments:       comments,
	}, nil
}

func mapToComment(data map[string]string) (*entity.Comment, error) {
	createdAt, err := time.Parse(time.RFC3339, data["createdAt"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, data["updatedAt"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse updatedAt: %w", err)
	}

	var replies []*entity.Comment
	if err := json.Unmarshal([]byte(data["replies"]), &replies); err != nil {
		replies = []*entity.Comment{}
	}

	var parentID *string
	if data["parentId"] != "" {
		parentIDValue := data["parentId"]
		parentID = &parentIDValue
	}

	return &entity.Comment{
		ID:        data["id"],
		PostID:    data["postId"],
		ParentID:  parentID,
		Content:   data["content"],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Replies:   replies,
	}, nil
}
