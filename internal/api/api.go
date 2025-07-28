package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/oThinas/bid/internal/services"
)

type Api struct {
	Router         *chi.Mux
	Sessions       *scs.SessionManager
	WsUpgrader     websocket.Upgrader
	UserService    services.UserService
	ProductService services.ProductService
	BidsService    services.BidsService
	AuctionLobby   services.AuctionLobby
}
