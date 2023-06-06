package models

import (
	"time"
	"gorm.io/datatypes"
	"github.com/google/uuid"

)

type Vaccancy struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title string `json:"title"`
	Location string `json:"location"`
	Salary string `json:"salary"`
	EmployerID uuid.UUID `json:"employer_id" gorm:"type:uuid;not null"`
	Type string `json:"type" gorm:"not null"`
	Requirements datatypes.JSONSlice[string] `json:"requirements" gorm:"varchar(64)[]"`
	Conditions datatypes.JSONSlice[string] `json:"conditions" gorm:"varchar(64)[]"`
	Info string `json:"info"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}