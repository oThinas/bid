package products

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/oThinas/bid/internal/validator"
)

type CreateProductRequest struct {
	SellerID    uuid.UUID `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductRequest) Valid(context.Context) validator.Evaluator {
	var ev validator.Evaluator

	ev.CheckField(validator.NotBlank(req.Name), "name", "this field cannot be empty")

	ev.CheckField(validator.NotBlank(req.Description), "description", "this field cannot be empty")
	ev.CheckField(
		validator.MinChars(req.Description, 10) && validator.MaxChars(req.Description, 255),
		"description",
		"this field must have between 10 and 255 characters",
	)

	ev.CheckField(req.BasePrice > 0, "base_price", "base price must be greater than 0")
	ev.CheckField(time.Until(req.AuctionEnd) >= minAuctionDuration, "auction_end", "this field must be at least 2 hours from now")

	return ev
}
