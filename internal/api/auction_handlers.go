package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/oThinas/bid/internal/services"
	"github.com/oThinas/bid/internal/utils"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductID := chi.URLParam(r, "productID")

	productID, err := uuid.Parse(rawProductID)
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "invalid product id",
		})
		return
	}

	_, err = api.ProductService.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			utils.EncodeJSON(w, r, http.StatusNotFound, map[string]string{
				"error": "no product with given id",
			})
			return
		}
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected error, try again later.",
		})
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), AuthenticatedUserID).(uuid.UUID)
	if !ok {
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected error, try again later.",
		})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productID]
	api.AuctionLobby.Unlock()

	if !ok {
		utils.EncodeJSON(w, r, http.StatusBadRequest, map[string]string{
			"error": "the auction has ended",
		})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "could not upgrade connection to a websocket protocol",
		})
		return
	}

	client := services.NewClient(conn, room, userID)

	room.Register <- client
	go client.ReadEventLoop()
	go client.WriteEventLoop()
}
