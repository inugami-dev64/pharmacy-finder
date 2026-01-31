package auth

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"pharmafinder/service"
	"pharmafinder/types"
	"pharmafinder/utils"
	"pharmafinder/web"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type ModeratorAuthController struct {
	repo           db.ModeratorUserRepository
	hasher         service.PasswordHasher
	sessionManager service.SessionManager
	logger         zerolog.Logger
}

func ProvideModeratorAuthContorller(
	repo db.ModeratorUserRepository,
	hasher service.PasswordHasher,
	sessionManager service.SessionManager,
) []web.Route {
	controller := &ModeratorAuthController{
		repo:           repo,
		hasher:         hasher,
		sessionManager: sessionManager,
		logger:         utils.GetLogger("API"),
	}

	return controller.GetRoutes()
}

func (handler *ModeratorAuthController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewRequestsHandler[ModeratorAuthController](handler.ModeratorLogin, "/mod/auth/login", []string{"POST"}),
		web.NewSecureRequestsHandler[ModeratorAuthController](handler.AdminRegistration, "/mod/auth/register", []string{"POST"}, web.NewSecurityChain[dto.ModeratorUserRegistrationDTO]().RuleWhenNoAdminUser(handler.repo)),
	}
}

// Authenticate moderator user and return its details along session token
//
// Path: `POST /api/v1/mod/auth/login`
//
// @Summary 		Authenticate moderator user
// @Description		Endpoint for performing login authentication for moderator users
// @Tags			Auth
// @Produce			json
// @Param			request body dto.ModeratorUserLoginDTO true "Login request body"
// @Success 		200 {object} dto.AuthenticatedModeratorUserResponseDTO
// @Failure			400 {object} types.HttpError
// @Router			/api/v1/mod/auth/login [post]
func (handler *ModeratorAuthController) ModeratorLogin(details *web.HttpRequestDetails[dto.ModeratorUserLoginDTO]) (int, interface{}, error) {
	user, err := handler.repo.FindUserByUsernameOrEmail(details.Body.Username).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if user == nil {
		return http.StatusBadRequest, types.NewHttpError(http.StatusNotFound, "User with provided username or email was not found"), nil
	}

	if !handler.hasher.Verify(details.Body.Password, user.Password) {
		return http.StatusBadRequest, types.NewHttpError(http.StatusForbidden, "Invalid password"), nil
	}

	user.LastLoginTimestamp = types.Time(time.Now().UTC())
	handler.repo.Store(user)

	userDTO := dto.AuthenticatedModeratorUserResponseDTO{
		ID:                    user.ID,
		Username:              user.Username,
		Email:                 user.Email,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		RegistrationTimestamp: user.RegistrationTimestamp,
		LastLoginTimestamp:    user.LastLoginTimestamp,
		Administrator:         user.Administrator,
		Session: dto.SessionTokenResponseDTO{
			Token:    handler.sessionManager.NewSessionToken(user),
			ValidFor: int64(service.SESSION_TTL.Seconds()),
		},
	}

	return http.StatusOK, userDTO, nil
}

// First time admin user registration endpoint
//
// Path: `POST /api/v1/mod/auth/register/admin“
//
// @Summary 		Register the initial admin account
// @Description		Endpoint which can be used to register the initial administrator account when admin accounts are not present
// @Tags			Auth
// @Produce			json
// @Param			request body dto.ModeratorUserRegistrationDTO true "Registration request body"
// @Success 		200 {object} dto.AuthenticatedModeratorUserResponseDTO
// @Failure			400 {object} types.HttpError
// @Failure			403 {object} types.HttpError
// @Router			/api/v1/mod/auth/register [post]
func (handler *ModeratorAuthController) AdminRegistration(details *web.HttpRequestDetails[dto.ModeratorUserRegistrationDTO]) (int, interface{}, error) {
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
		Administrator:         true,
	}

	err = handler.repo.Store(&user)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		handler.logger.Warn().Msgf("User with username '%s' and email '%s' already exists", details.Body.Username, details.Body.Email)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Username or email is already in use"), nil
	} else if err != nil {
		handler.logger.Warn().Msgf("Failed to register an admin account: %v", err)
		return http.StatusInternalServerError, nil, err
	}

	userDTO := dto.AuthenticatedModeratorUserResponseDTO{
		ID:                    user.ID,
		Username:              user.Username,
		Email:                 user.Email,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		RegistrationTimestamp: user.RegistrationTimestamp,
		LastLoginTimestamp:    user.LastLoginTimestamp,
		Administrator:         user.Administrator,
		Session: dto.SessionTokenResponseDTO{
			Token:    handler.sessionManager.NewSessionToken(&user),
			ValidFor: int64(service.SESSION_TTL.Seconds()),
		},
	}

	return http.StatusOK, userDTO, nil
}
