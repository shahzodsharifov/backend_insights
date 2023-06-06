package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	PostID uuid.UUID `json:"post_id" gorm:"type:uuid;not null"`
}



type Post struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Thumbnail *string  `gorm:"not null;default:'default.png'"`
	Headline string `json:"headline"`
	Subtitle string `json:"subtitle"`
	Body string `json:"body"`
	AuthorID uuid.UUID `json:"author_id" gorm:"type:uuid;not null"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID"`
	Likes []Like `json:"likes" gorm:"foreignKey:PostID"`
	Category string `json:"category"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}