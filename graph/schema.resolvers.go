package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apartapatia/wall_of_comments/graph/model"
	"github.com/apartapatia/wall_of_comments/internal/entity"
	"github.com/google/uuid"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsDisabled bool) (*model.Post, error) {
	post := &entity.Post{
		ID:             uuid.New().String(),
		Title:          title,
		Content:        content,
		CommentsActive: !commentsDisabled,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	savedPost, err := createAndSaveEntity(post, r.Repo.CreatePost)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:             savedPost.ID,
		Title:          savedPost.Title,
		Content:        savedPost.Content,
		CommentsActive: savedPost.CommentsActive,
		CreatedAt:      savedPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      savedPost.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*model.Comment, error) {
	comment := &entity.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	savedComment, err := createAndSaveEntity(comment, r.Repo.CreateComment)
	if err != nil {
		return nil, err
	}

	if parentID != nil {
		parentComment, err := r.Repo.GetCommentById(*parentID)
		if parentComment.PostID != savedComment.PostID || err != nil {
			return nil, ErrParentCommentNotFound
		}
	}

	return buildCommentModel(savedComment), nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.Repo.GetPosts()
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	var result []*model.Post
	for _, post := range posts {
		comments, err := r.Repo.GetCommentsForPost(post.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get comments for post with ID %s: %w", post.ID, err)
		}

		commentModels, err := buildCommentTree(comments)
		if err != nil {
			return nil, err
		}

		postModel := &model.Post{
			ID:             post.ID,
			Title:          post.Title,
			Content:        post.Content,
			CommentsActive: post.CommentsActive,
			CreatedAt:      post.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      post.UpdatedAt.Format(time.RFC3339),
			Comments:       commentModels,
		}

		result = append(result, postModel)
	}

	return result, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	post, err := r.Repo.GetPostById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	comments, err := r.Repo.GetCommentsForPost(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for post: %w", err)
	}

	commentModels, err := buildCommentTree(comments)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:             post.ID,
		Title:          post.Title,
		Content:        post.Content,
		CommentsActive: post.CommentsActive,
		CreatedAt:      post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      post.UpdatedAt.Format(time.RFC3339),
		Comments:       commentModels,
	}, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) ([]*model.Comment, error) {
	comments, err := r.Repo.GetCommentsForPostWithLimitAndOffset(postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for post with ID %s: %w", postID, err)
	}

	var commentModels []*model.Comment
	for _, comment := range comments {
		commentModels = append(commentModels, buildCommentModel(comment))
	}

	return commentModels, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var ErrParentCommentNotFound = errors.New("parent comment not found")

func buildCommentModel(comment *entity.Comment) *model.Comment {
	return &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		Replies:   []*model.Comment{},
	}
}
func buildCommentTree(comments []*entity.Comment) ([]*model.Comment, error) {
	commentMap := make(map[string]*model.Comment)
	for _, comment := range comments {
		commentModel := buildCommentModel(comment)
		commentMap[comment.ID] = commentModel
	}

	var commentModels []*model.Comment
	for _, comment := range commentMap {
		if comment.ParentID == nil {
			commentModels = append(commentModels, comment)
		} else {
			if parentComment, exists := commentMap[*comment.ParentID]; exists {
				parentComment.Replies = append(parentComment.Replies, comment)
			}
		}
	}

	return commentModels, nil
}
func createAndSaveEntity[T any](entity T, saveFunc func(T) (T, error)) (T, error) {
	savedEntity, err := saveFunc(entity)
	if err != nil {
		return savedEntity, fmt.Errorf("failed to create entity: %w", err)
	}
	return savedEntity, nil
}
