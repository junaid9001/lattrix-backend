package dto

type ProfileUpdateRequet struct {
	Username *string `json:"username"`
	Email    *string `json:"email" validate:"email"`
}
