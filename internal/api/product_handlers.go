package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/oThinas/bid/internal/usecase/products"
	"github.com/oThinas/bid/internal/utils"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeJSON[products.CreateProductRequest](r)
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), AuthenticatedUserID).(uuid.UUID)

	if !ok {
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	id, err := api.ProductService.CreateProduct(
		r.Context(),
		userID,
		data.Name,
		data.Description,
		data.BasePrice,
		data.AuctionEnd,
	)
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	utils.EncodeJSON(w, r, http.StatusOK, map[string]uuid.UUID{
		"data": id,
	})
}
