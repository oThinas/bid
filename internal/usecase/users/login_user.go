package users

import (
	"context"

	"github.com/oThinas/bid/internal/validator"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var ev validator.Evaluator

	ev.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "this field must be a valid email address")
	ev.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")

	return ev
}
