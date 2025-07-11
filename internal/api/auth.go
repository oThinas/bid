package api

import (
	"net/http"

	"github.com/oThinas/bid/internal/utils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), AuthenticatedUserID) {
			utils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{
				"error": "must be logged in",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
