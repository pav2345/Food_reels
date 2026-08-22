package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                       uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FullName                 string     `gorm:"column:full_name;not null"`
	Email                    string     `gorm:"not null;uniqueIndex"`
	Password                 string     `gorm:"not null"`
	Latitude                 *float64   `gorm:"column:latitude"`
	Longitude                *float64   `gorm:"column:longitude"`
	LastLocationUpdated      *time.Time `gorm:"column:last_location_updated"`
	SavedLatitude            *float64   `gorm:"column:saved_latitude"`
	SavedLongitude           *float64   `gorm:"column:saved_longitude"`
	CurrentLatitude          *float64   `gorm:"column:current_latitude"`
	CurrentLongitude         *float64   `gorm:"column:current_longitude"`
	CurrentLocationUpdatedAt *time.Time `gorm:"column:current_location_updated_at"`
	CreatedAt                time.Time  `gorm:"not null"`
	UpdatedAt                time.Time  `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}

type FoodPartner struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name             string    `gorm:"not null"`
	ContactName      string    `gorm:"column:contact_name;not null"`
	Phone            string    `gorm:"not null"`
	Address          string    `gorm:"not null"`
	Email            string    `gorm:"not null;uniqueIndex"`
	Password         string    `gorm:"not null"`
	Latitude         *float64  `gorm:"column:latitude"`
	Longitude        *float64  `gorm:"column:longitude"`
	DeliveryRadiusKM float64   `gorm:"column:delivery_radius_km;not null;default:6"`
	OpeningTime      string    `gorm:"column:opening_time"`
	ClosingTime      string    `gorm:"column:closing_time"`
	Rating           float64   `gorm:"column:rating;not null;default:0"`
	ProfileImage     string    `gorm:"column:profile_image"`
	CoverImage       string    `gorm:"column:cover_image"`
	FoodImages       string    `gorm:"column:food_images;type:text"`
	MenuImages       string    `gorm:"column:menu_images;type:text"`
	MenuPDF          string    `gorm:"column:menu_pdf"`
}

func (FoodPartner) TableName() string {
	return "food_partners"
}

type Food struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `gorm:"not null"`
	VideoURL      string    `gorm:"column:video_url;not null"`
	Description   *string
	ImageURL      *string   `gorm:"column:image_url"`
	ThumbnailURL  *string   `gorm:"column:thumbnail_url"`
	Category      *string   `gorm:"column:category"`
	Available     bool      `gorm:"column:available;not null;default:true"`
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
