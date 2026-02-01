package users

import (
	"fmt"
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"pharmafinder/service"
	"pharmafinder/types"
	"pharmafinder/utils"
	"pharmafinder/web"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type ModeratorUserController struct {
	repo         db.ModeratorUserRepository
	hasher       service.PasswordHasher
	tokenManager service.SessionManager
	logger       zerolog.Logger
}

func ProvideModeratorUserController(
	repo db.ModeratorUserRepository,
	hasher service.PasswordHasher,
	tokenManager service.SessionManager,
) []web.Route {
	controller := &ModeratorUserController{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
		logger:       utils.GetLogger("API"),
	}

	return controller.GetRoutes()
}

func (handler *ModeratorUserController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewSecureRequestsHandler[ModeratorUserController](handler.GetAllUsers, "/mod/users", []string{"GET"}, web.NewSecurityChain[web.EmptyBody]().RuleAdmin(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.GetAuthenticatedUserProfile, "/mod/users/me", []string{"GET"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.CreateNewModeratorUser, "/mod/users", []string{"POST"}, web.NewSecurityChain[dto.ModeratorUserRegistrationDTO]().RuleAdmin(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.UpdateCurrentModeratorUser, "/mod/users/me", []string{"PATCH"}, web.NewSecurityChain[dto.ModeratorUserUpdateDTO]().RuleAuthenticated(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.UpdateModeratorUser, "/mod/users/{id}", []string{"PATCH"}, web.NewSecurityChain[dto.AdminUserUpdateDTO]().RuleAdmin(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.DeleteCurrentModeratorUser, "/mod/users/me", []string{"DELETE"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.DeleteModeratorUser, "/mod/users/{id}", []string{"DELETE"}, web.NewSecurityChain[web.EmptyBody]().RuleAdmin(handler.repo, handler.tokenManager)),
	}
}

// Get all moderator users
//
// Path: `GET /api/v1/mod/users`
//
// @Summary 		Query a list of all moderator users
// @Description		Endpoint for listing data about all moderator users
// @Tags			Users
// @Security		Bearer
// @Produce			json
// @Success			200 {array} dto.ModeratorUserProfileDTO
// @Failure			403 {object} types.HttpError
// @Router			/api/v1/mod/users [get]
func (handler *ModeratorUserController) GetAllUsers(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	data, err := handler.repo.FindAllUsers().QueryAll()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, data, nil
}

// Get currently authenticated user's details
//
// Path: `GET /api/v1/mod/users/me`
//
// @Summary			Query currently authenticated user's profile information
// @Description		Endpoint for listing profile data about currently authenticated user
// @Tags			Users
// @Security		Bearer
// @Produce 		json
// @Success			200 {object} dto.ModeratorUserProfileDTO
// @Failure			403 {object} types.HttpError
// @Router			/api/v1/mod/users/me [get]
func (handler *ModeratorUserController) GetAuthenticatedUserProfile(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	dto := dto.ModeratorUserProfileDTO{
		ID:                    details.AuthenticatedUser.ID,
		Username:              details.AuthenticatedUser.Username,
		Email:                 details.AuthenticatedUser.Email,
		FirstName:             details.AuthenticatedUser.FirstName,
		LastName:              details.AuthenticatedUser.LastName,
		RegistrationTimestamp: details.AuthenticatedUser.RegistrationTimestamp,
		LastLoginTimestamp:    details.AuthenticatedUser.LastLoginTimestamp,
		Administrator:         details.AuthenticatedUser.Administrator,
	}

	return http.StatusOK, dto, nil
}

// Create a new moderator account
//
// Path: `POST /api/v1/mod/users`
//
// @Summary 		Create a new moderator account
// @Description		Endpoint which allows administrator accounts to create new moderator accounts
// @Tags			Users
// @Security		Bearer
// @Produce 		json
// @Param			request body dto.ModeratorUserRegistrationDTO true "Moderator user registration body"
// @Success			201 {object} dto.ModeratorUserProfileDTO
// @Failure			400 {object} types.HttpError
// @Router			/api/v1/mod/users [post]
func (handler *ModeratorUserController) CreateNewModeratorUser(details *web.HttpRequestDetails[dto.ModeratorUserRegistrationDTO]) (int, interface{}, error) {
	pwdHash, err := handler.hasher.CreatePasswordHash(details.Body.Password)
	if err == service.ErrPasswordTooLong {
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Password is too long"), nil
	} else if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	user := entity.ModeratorUser{
		Username:              details.Body.Username,
		Email:                 details.Body.Email,
		Password:              pwdHash,
		FirstName:             details.Body.FirstName,
		LastName:              details.Body.LastName,
		RegistrationTimestamp: types.Time(time.Now().UTC()),
		LastLoginTimestamp:    types.Time(time.Now().UTC()),
		Administrator:         false,
	}

	err = handler.repo.Store(&user)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		handler.logger.Warn().Msgf("User with username '%s' and email '%s' already exists", details.Body.Username, details.Body.Email)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Username or email is already in use"), nil
	} else if err != nil {
		handler.logger.Warn().Msgf("Failed to register an admin account: %v", err)
		return http.StatusInternalServerError, nil, err
	}

	userDTO := dto.ModeratorUserProfileDTO{
		ID:                    user.ID,
		Username:              user.Username,
		Email:                 user.Email,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		RegistrationTimestamp: user.RegistrationTimestamp,
		LastLoginTimestamp:    user.LastLoginTimestamp,
		Administrator:         user.Administrator,
	}

	return http.StatusCreated, userDTO, nil
}
}

// Modify currently authenticated moderators account
//
// Path: `PATCH /api/v1/mod/users/me`
func (handler *ModeratorUserController) UpdateCurrentModeratorUser(details *web.HttpRequestDetails[dto.ModeratorUserUpdateDTO]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Modify someones user account
//
// Path: `PATCH /api/v1/mod/users/{id}`
func (handler *ModeratorUserController) UpdateModeratorUser(details *web.HttpRequestDetails[dto.AdminUserUpdateDTO]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Delete my user account
//
// Path: `DELETE /api/v1/mod/users/me`
func (handler *ModeratorUserController) DeleteCurrentModeratorUser(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Delete someones user account
//
// Path: `DELETE /api/v1/mod/users/{id}`
func (handler *ModeratorUserController) DeleteModeratorUser(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}
