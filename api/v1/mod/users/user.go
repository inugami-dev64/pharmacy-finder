package users

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/service"
	"pharmafinder/utils"
	"pharmafinder/web"

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

// Get all moderator user details
//
// Path: `GET /api/v1/mod/users`
func (handler *ModeratorUserController) GetAllUsers(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Get currently authenticated user's details
//
// Path: `GET /api/v1/mod/users/me`
func (handler *ModeratorUserController) GetAuthenticatedUserProfile(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Create a new moderator account
//
// Path: `POST /api/v1/mod/users`
func (handler *ModeratorUserController) CreateNewModeratorUser(details *web.HttpRequestDetails[dto.ModeratorUserRegistrationDTO]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
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
