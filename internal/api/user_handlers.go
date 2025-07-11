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
	data, problems, err := utils.DecodeJSON[users.LoginUserRequest](r)
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.EncodeJSON(w, r, http.StatusBadRequest, map[string]string{
				"error": "invalid email or password",
			})
			return
		}

		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), AuthenticatedUserID, id)
	utils.EncodeJSON(w, r, http.StatusOK, map[string]string{
		"data": "logged in successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		utils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), AuthenticatedUserID)
	utils.EncodeJSON(w, r, http.StatusOK, map[string]string{
		"data": "logged out successfully",
	})
}
