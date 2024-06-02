package entity

import (
	"time"
)

type Post struct {
	ID             string     `gorm:"primaryKey" json:"id"`
	Title          string     `gorm:"not null" json:"title" validate:"required"`
	Content        string     `gorm:"not null" json:"content" validate:"required"`
	CommentsActive bool       `gorm:"not null" json:"commentsActive" validate:"required"`
	CreatedAt      time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"index" json:"updatedAt"`
	Comments       []*Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}
