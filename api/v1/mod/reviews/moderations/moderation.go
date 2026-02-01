package moderations

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/service"
	"pharmafinder/utils"
	"pharmafinder/web"

	"github.com/rs/zerolog"
)

type ReviewModerationController struct {
	userRepo db.ModeratorUserRepository
	// TODO: Create and implement db.CommentModerationRepository
	// modRepo db.CommentModerationRepository
	tokenManager service.SessionManager
	logger       zerolog.Logger
}

func ProvideReviewModerationController(
	userRepo db.ModeratorUserRepository,
	tokenManager service.SessionManager,
) []web.Route {
	controller := &ReviewModerationController{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		logger:       utils.GetLogger("API"),
	}

	return controller.GetRoutes()
}

func (handler *ReviewModerationController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewSecureRequestsHandler[ReviewModerationController](handler.GetReviewModerations, "/mod/reviews/{id}/moderations", []string{"GET"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ReviewModerationController](handler.CreateReviewModeration, "/mod/reviews/{id}/moderations", []string{"POST"}, web.NewSecurityChain[dto.CommentReviewModificationDTO]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ReviewModerationController](handler.CreateReviewModeration, "/mod/reviews/{reviewID}/moderations/{modID}", []string{"PUT"}, web.NewSecurityChain[dto.CommentReviewModificationDTO]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ReviewModerationController](handler.CreateReviewModeration, "/mod/reviews/{reviewID}/moderations/{modID}", []string{"DELETE"}, web.NewSecurityChain[dto.CommentReviewModificationDTO]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
	}
}

// Get all review moderations for specific review
//
// Path: `GET /api/v1/mod/reviews/{id}/moderations`
func (handler *ReviewModerationController) GetReviewModerations(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Create a new moderation review for specific comment
//
// Path: `POST /api/v1/mod/reviews/{id}/moderations`
func (handler *ReviewModerationController) CreateReviewModeration(details *web.HttpRequestDetails[dto.CommentReviewModificationDTO]) (int, interface{}, error) {
	return http.StatusCreated, []string{}, nil
}

// Modify moderation review
//
// Path: `PUT /api/v1/mod/reviews/{reviewID}/moderations/{modID}`
func (handler *ReviewModerationController) ModifyReviewModeration(details *web.HttpRequestDetails[dto.CommentReviewModificationDTO]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}

// Delete moderation review
//
// Path: `DELETE /api/v1/mod/reviews/{reviewID}/moderations/{modID}`
func (handler *ReviewModerationController) DeleteReviewModeration(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	return http.StatusOK, []string{}, nil
}
