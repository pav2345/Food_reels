package dto

type RegisterUserRequest struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
