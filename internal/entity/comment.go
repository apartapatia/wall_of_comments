package entity

import (
	"time"
)

type Comment struct {
	ID        string     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    string     `gorm:"not null" json:"postId" validate:"required"`
	ParentID  *string    `gorm:"index" json:"parentId,omitempty"`
	Content   string     `gorm:"not null;size:2000" json:"content" validate:"required,max=2000"`
	CreatedAt time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"index" json:"updatedAt"`
	Replies   []*Comment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"replies,omitempty"`
}
