package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/oThinas/bid/internal/services"
)

type Api struct {
	Router         *chi.Mux
	Sessions       *scs.SessionManager
	UserService    services.UserService
	ProductService services.ProductService
}
