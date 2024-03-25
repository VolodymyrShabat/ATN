package models

type ForgotPasswordInput struct {
	Email string `json:"email,omitempty" validate:"required"`
}

type ResetPasswordInput struct {
	Password string `json:"password,omitempty" validate:"required"`
}
