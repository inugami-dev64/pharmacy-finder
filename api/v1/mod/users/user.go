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
		web.NewSecureRequestsHandler[ModeratorUserController](handler.GetUserByID, "/mod/users/{id}", []string{"GET"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.CreateNewModeratorUser, "/mod/users", []string{"POST"}, web.NewSecurityChain[dto.ModeratorUserRegistrationDTO]().RuleAdmin(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.UpdateCurrentModeratorUser, "/mod/users/me", []string{"PATCH"}, web.NewSecurityChain[dto.ModeratorUserUpdateDTO]().RuleAuthenticated(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.UpdateModeratorUser, "/mod/users/{id}", []string{"PATCH"}, web.NewSecurityChain[dto.AdminUserUpdateDTO]().RuleAdmin(handler.repo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ModeratorUserController](handler.DeleteCurrentModeratorUser, "/mod/users/me", []string{"DELETE"}, web.NewSecurityChain[dto.ModeratorUserDeletionDTO]().RuleAuthenticated(handler.repo, handler.tokenManager)),
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

// Find a moderator user by specified ID
//
// Path: `GET /api/v1/mod/users/{id}`
//
// @Summary 		Query data about a specific moderator user
// @Description		Endpoint for querying data about a specific moderator user
// @Tags			Users
// @Security		Bearer
// @Produce			json
// @Success			200 {object} dto.ModeratorUserProfileDTO
// @Failure			400 {object} types.HttpError
// @Failure			403 {object} types.HttpError
// @Failure			404 {object} types.HttpError
// @Router			/api/v1/mod/users/{id} [get]
func (handler *ModeratorUserController) GetUserByID(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed user ID %v", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID variable"), nil
	}

	user, err := handler.repo.FindUserByID(types.UUID(id)).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, nil
	} else if user == nil {
		return http.StatusNotFound, types.NewHttpError(http.StatusNotFound, "Not found"), nil
	}

	dto := dto.ModeratorUserProfileDTO{
		ID:                    user.ID,
		Username:              user.Username,
		Email:                 user.Email,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		RegistrationTimestamp: user.RegistrationTimestamp,
		LastLoginTimestamp:    user.LastLoginTimestamp,
		Administrator:         user.Administrator,
	}

	return http.StatusOK, dto, nil
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

func (handler *ModeratorUserController) updateUser(user *entity.ModeratorUser, body *dto.AdminUserUpdateDTO, me bool) (int, interface{}, error) {
	// Update password if present
	if body.Password != "" {
		var err error
		user.Password, err = handler.hasher.CreatePasswordHash(body.Password)

		if err == service.ErrPasswordTooLong {
			if me {
				handler.logger.Warn().Msgf("User '%s' tried to change their password to something too long", user.Username)
			} else {
				handler.logger.Warn().Msgf("Requested password update for user '%s' failed, password too long", user.Username)
			}
			return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Password is too long"), nil
		} else if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	// Update other data fields
	user.Email = body.Email
	user.FirstName = body.FirstName
	user.LastName = body.LastName
	if !me && body.Administrator != nil {
		user.Administrator = *body.Administrator
	}

	err := handler.repo.Store(user)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		if me {
			handler.logger.Warn().Msgf("User '%s' tried to change their email to '%s', which is already in use", user.Username, user.Email)
		} else {
			handler.logger.Warn().Msgf("Requested email change for user '%s' failed, email already in use")
		}
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Email address is already in use"), nil
	} else if err != nil {
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

	return http.StatusOK, userDTO, nil
}

// Modify currently authenticated moderators account
//
// Path: `PATCH /api/v1/mod/users/me`
//
// @Summary			Modify currently authenticated user's details
// @Description		Endpoint for modifying currently authenticated user's details (requires password verification)
// @Tags			Users
// @Security		Bearer
// @Produce 		json
// @Param			request body dto.ModeratorUserUpdateDTO true "Moderator user update request body"
// @Success			200 {object} dto.ModeratorUserProfileDTO
// @Failure			400 {object} types.HttpError
// @Failure			403 {object} types.HttpError
// @Router			/api/v1/mod/users/me [patch]
func (handler *ModeratorUserController) UpdateCurrentModeratorUser(details *web.HttpRequestDetails[dto.ModeratorUserUpdateDTO]) (int, interface{}, error) {
	user := details.AuthenticatedUser
	if !handler.hasher.Verify(details.Body.CurrentPassword, user.Password) {
		handler.logger.Warn().Msgf("User '%s' tried to modify their profile with invalid password", user.Username)
		return http.StatusForbidden, types.NewHttpError(http.StatusForbidden, "Invalid password"), nil
	}

	return handler.updateUser(user, &details.Body.AdminUserUpdateDTO, true)
}

// Modify someone's user account
//
// Path: `PATCH /api/v1/mod/users/{id}`
//
// @Summary 		Modify other user's data
// @Description		Endpoint for administrators so that they could update data for other users
// @Tags			Users
// @Security		Bearer
// @Produce 		json
// @Param			request body dto.AdminUserUpdateDTO true "User update request body"
// @Param			id path string true "User ID"
// @Success			200 {object} dto.ModeratorUserProfileDTO
// @Failure			400 {object} types.HttpError
// @Failure			404 {object} types.HttpError
// @Router			/api/v1/mod/users/{id} [patch]
func (handler *ModeratorUserController) UpdateModeratorUser(details *web.HttpRequestDetails[dto.AdminUserUpdateDTO]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	user, err := handler.repo.FindUserByID(types.UUID(id)).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	} else if user == nil {
		handler.logger.Warn().Msgf("User with ID '%s' was not found", id.String())
		return http.StatusNotFound, types.NewHttpError(http.StatusNotFound, fmt.Sprintf("User with ID '%s' does not exist", id.String())), nil
	}

	return handler.updateUser(user, &details.Body, false)
}

// Delete my user account
//
// Path: `DELETE /api/v1/mod/users/me`
//
// @Summary			Delete currently authenticated moderator account
// @Description 	Endpoint for deleting currently authenticated user's account
// @Tags			Users
// @Security		Bearer
// @Produce 		json
// @Param			request body dto.ModeratorUserDeletionDTO true "User deletion request body"
// @Success			200 {object} dto.ModeratorUserProfileDTO
// @Failure			400 {object} types.HttpError
// @Router 			/api/v1/mod/users/me [delete]
func (handler *ModeratorUserController) DeleteCurrentModeratorUser(details *web.HttpRequestDetails[dto.ModeratorUserDeletionDTO]) (int, interface{}, error) {
	if !handler.hasher.Verify(details.Body.CurrentPassword, details.AuthenticatedUser.Password) {
		handler.logger.Warn().Msgf("User '%s' tried to delete their account with invalid password", details.AuthenticatedUser.Username)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Invalid password"), nil
	}

	// TODO: Verify that if user is admin, that they are not the only one

	user, err := handler.repo.Delete(details.AuthenticatedUser.ID).Query()
	if err != nil {
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

	return http.StatusOK, userDTO, nil
}

// Delete someones user account
//
// Path: `DELETE /api/v1/mod/users/{id}`
//
// @Summary 		Delete someone's moderator account
// @Description 	Endpoint for administrators so that they could delete someone else's moderator accounts
// @Tags			Users
// @Security		Bearer
// @Produce			json
// @Param			id path string true "User ID to delete"
// @Success			200 {object} dto.ModeratorUserProfileDTO
// @Failure			404 {object} types.HttpError
// @Router			/api/v1/mod/users/{id} [delete]
func (handler *ModeratorUserController) DeleteModeratorUser(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID parameter '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID parameter"), nil
	}

	user, err := handler.repo.Delete(types.UUID(id)).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	} else if user == nil {
		return http.StatusNotFound, types.NewHttpError(http.StatusNotFound, "Not found"), nil
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

	return http.StatusOK, userDTO, nil
}
