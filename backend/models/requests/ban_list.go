package requests

type BanRequest struct {
	Email string `json:"email" validate:"required,email"`
}
