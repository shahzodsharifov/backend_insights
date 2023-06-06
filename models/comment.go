package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Body string `gorm:"not null"`
	AuthorID uuid.UUID `json:"author_id" gorm:"type:uuid;not null"`
	PostID uuid.UUID `json:"post_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`

}