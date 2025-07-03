package users

import (
	"context"

	"github.com/oThinas/bid/internal/validator"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserRequest) Valid(context.Context) validator.Evaluator {
	var ev validator.Evaluator

	ev.CheckField(validator.NotBlank(req.Username), "username", "this field cannot be empty")
	ev.CheckField(validator.MinChars(req.Password, 8), "password", "this field must have at least 8 characters")

	ev.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	ev.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "this field must be a valid email address")

	ev.CheckField(validator.NotBlank(req.Bio), "bio", "this field cannot be empty")
	ev.CheckField(
		validator.MinChars(req.Bio, 10) && validator.MaxChars(req.Bio, 255),
		"bio",
		"this field must have between 10 and 255 characters",
	)

	return ev
}
