package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FullName  string    `gorm:"column:full_name;not null"`
	Email     string    `gorm:"not null;uniqueIndex"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

type FoodPartner struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	ContactName string    `gorm:"column:contact_name;not null"`
	Phone       string    `gorm:"not null"`
	Address     string    `gorm:"not null"`
	Email       string    `gorm:"not null;uniqueIndex"`
	Password    string    `gorm:"not null"`
}

func (FoodPartner) TableName() string {
	return "food_partners"
}

type Food struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `gorm:"not null"`
	VideoURL      string    `gorm:"column:video_url;not null"`
	Description   *string
	FoodPartnerID uuid.UUID `gorm:"column:food_partner_id;type:uuid;not null;index"`
	Likes         int       `gorm:"not null;default:0"`
	Saves         int       `gorm:"not null;default:0"`
	Version       int       `gorm:"column:version;not null;default:0"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`

	FoodPartner FoodPartner `gorm:"foreignKey:FoodPartnerID"`
}

func (Food) TableName() string {
	return "foods"
}

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_likes_user_food"`
	FoodID    uuid.UUID `gorm:"column:food_id;type:uuid;not null;uniqueIndex:idx_likes_user_food;index"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Like) TableName() string {
	return "likes"
}

type Save struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_saves_user_food;index"`
	FoodID    uuid.UUID `gorm:"column:food_id;type:uuid;not null;uniqueIndex:idx_saves_user_food"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Food Food `gorm:"foreignKey:FoodID"`
}

func (Save) TableName() string {
	return "saves"
}
