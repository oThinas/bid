package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/oThinas/bid/internal/services"
	"github.com/oThinas/bid/internal/usecase/users"
	"github.com/oThinas/bid/internal/utils"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeJSON[users.CreateUserRequest](r)
	if err != nil {
		_ = utils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.Username, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedUsernameOrEmail) {
			_ = utils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]string{
				"error": "username or email already exists",
			})
			return
		}
	}

	_ = utils.EncodeJSON(w, r, http.StatusCreated, map[string]uuid.UUID{
		"data": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {

}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {

}
