package dto

type RegisterUserRequest struct {
	FullName  string   `json:"fullName" validate:"required"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type LoginRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type RegisterFoodPartnerRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Address     string `json:"address" validate:"required"`
	ContactName string `json:"contactName" validate:"required"`
}

type FoodActionRequest struct {
	FoodID string `json:"foodId" validate:"required"`
}

type CreateFoodRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

type FeedRequest struct {
	Latitude  *float64
	Longitude *float64
	IPAddress string
	Cursor    string
	Limit     int
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type UpdateDeliveryRadiusRequest struct {
	DeliveryRadiusKM float64 `json:"deliveryRadiusKm" binding:"required,gt=0"`
}

type UpdateWorkingHoursRequest struct {
	OpeningTime string `json:"openingTime" binding:"required"`
	ClosingTime string `json:"closingTime" binding:"required"`
}

type UpdateFoodPartnerProfileRequest struct {
	Name         *string  `json:"name"`
	ContactName  *string  `json:"contactName"`
	Phone        *string  `json:"phone"`
	Address      *string  `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	ProfileImage *string  `json:"profileImage"`
	CoverImage   *string  `json:"coverImage"`
}
