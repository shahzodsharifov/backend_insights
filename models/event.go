package models

import (
	"time"
	"gorm.io/datatypes"
	"github.com/google/uuid"
)

type Event struct {
	ID uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title string `json:"title"`
	Date string `json:"date"`
	Info string `json:"info"`
	EventType string `json:"eventType"`
	Location string `json:"location"`
	Speakers datatypes.JSONSlice[string] `json:"speakers" gorm:"varchar(64)[]"`
	OrganizerID uuid.UUID `json:"organizer_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
}