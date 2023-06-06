package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string     `gorm:"type:varchar(100);not null"`
	Username  string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Bio 	  string  `gorm:"type:varchar(100); default:"Qo'shimcha" ma'lumot; not null"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(100);not null"`
	Role      *string    `gorm:"type:varchar(50);default:'user';not null"`
	Photo     *string    `gorm:"not null;default:'default.png'"`
	Verified  *bool      `gorm:"not null;default:false"`
	Followers []UserRelationship `json:"followers" gorm:"foreignKey:FollowingID"`
	Following []UserRelationship `json:"following" gorm:"foreignKey:FollowerID"`
	Posts []Post `json:"posts" gorm:"foreignKey:AuthorID"`
	Vaccancies []Vaccancy `json:"vaccancies" gorm:"foreignKey:EmployerID"`
	Events []Event `json:"events" gorm:"foreignKey:OrganizerID"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

type UserRelationship struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FollowerID   uuid.UUID `gorm:"type:uuid;not null"`
	FollowingID  uuid.UUID `gorm:"type:uuid;not null"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T)[]*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err !=nil {
		for _,err:= range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors= append(errors, &element)
		}
	}

	return errors
}


type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Username      string   `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	Photo           string `json:"photo"`
	Role string `json:"role"`
}

type CategoryInput struct {
	Category string `json:"category"`
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type FollowInput struct {
	FollowerID uuid.UUID `json:"followerID"`
	FollowingID uuid.UUID `json:"followingID"`
}

type FindUserInput struct {
	UserID uuid.UUID `json:"userId"`
}

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username: user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Role:      *user.Role,
		Photo:     *user.Photo,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Followers []UserRelationship `json:"followers" gorm:"foreignKey:FollowingID"`
	Following []UserRelationship `json:"following" gorm:"foreignKey:FollowerID"`
	Posts []Post `json:"posts" gorm:"foreignKey:AuthorID"`
	Vaccancies []Vaccancy `json:"vaccancies" gorm:"foreignKey:EmployerID"`
	Events []Event `json:"events" gorm:"foreignKey:OrganizerID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}




type Users struct {
	Users []User `json:"users"`
}

