package moderations

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
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type ReviewModerationController struct {
	userRepo     db.ModeratorUserRepository
	modRepo      db.CommentModerationRepository
	tokenManager service.SessionManager
	logger       zerolog.Logger
}

func ProvideReviewModerationController(
	userRepo db.ModeratorUserRepository,
	modRepo db.CommentModerationRepository,
	tokenManager service.SessionManager,
) []web.Route {
	controller := &ReviewModerationController{
		userRepo:     userRepo,
		modRepo:      modRepo,
		tokenManager: tokenManager,
		logger:       utils.GetLogger("API"),
	}

	return controller.GetRoutes()
}

func (handler *ReviewModerationController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewSecureRequestsHandler[ReviewModerationController](handler.CreateReviewModeration, "/mod/reviews/{id}/moderations", []string{"POST"}, web.NewSecurityChain[dto.CommentReviewModificationDTO]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ReviewModerationController](handler.GetReviewModerations, "/mod/reviews/{id}/moderations", []string{"GET"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ReviewModerationController](handler.ModifyReviewModeration, "/mod/reviews/{reviewID}/moderations/{modID}", []string{"PATCH"}, web.NewSecurityChain[dto.CommentReviewModificationDTO]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
		web.NewSecureRequestsHandler[ReviewModerationController](handler.DeleteReviewModeration, "/mod/reviews/{reviewID}/moderations/{modID}", []string{"DELETE"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
	}
}

// Get all review moderations for specific review
//
// Path: `GET /api/v1/mod/reviews/{id}/moderations`
//
// @Summary 		Query moderations for specified review
// @Description 	Endpoint for querying moderations for specified review
// @Tags			Moderation
// @Produce			json
// @Security 		Bearer
// @Param			id path int true "ID of the review whose moderations to query for"
// @Success			200 {array} dto.CommentReviewResultDTO
// @Failure			400 {object} types.HttpError
// @Failure			403 {object} types.HttpError
// @Router			/api/v1/mod/reviews/{id}/moderations [get]
func (handler *ReviewModerationController) GetReviewModerations(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	moderations, err := handler.modRepo.FindCommentModerationsForReview(id).QueryAll()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, moderations, nil
}

// Create a new moderation review for specific comment
//
// Path: `POST /api/v1/mod/reviews/{id}/moderations`
//
// @Summary 		Create a new moderation for review comment
// @Description		Endpoint for creating a new moderation for specified review
// @Tags 			Moderation
// @Produce			json
// @Security		Bearer
// @Param			id path int true "ID of the review whose moderation to create"
// @Param			request body dto.CommentReviewModificationDTO true "Moderation review creation request body"
// @Success			200 {object} entity.CommentReview
// @Failure			400 {object} types.HttpError
// @Failure			403 {object} types.HttpError
// @Router			/api/v1/mod/reviews/{id}/moderations [post]
func (handler *ReviewModerationController) CreateReviewModeration(details *web.HttpRequestDetails[dto.CommentReviewModificationDTO]) (int, interface{}, error) {
	id, err := strconv.ParseInt(details.PathVars["id"], 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", details.PathVars["id"])
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	moderation := entity.CommentReview{
		CommentID:         id,
		ModeratorID:       details.AuthenticatedUser.ID,
		Result:            details.Body.Result,
		ModeratorComment:  details.Body.ModeratorComment,
		MarkedForDeletion: details.Body.MarkedForDeletion,
		ReviewedAt:        types.Time(time.Now().UTC()),
	}

	err = handler.modRepo.Store(&moderation)
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == pq.ErrorCode("23505") {
			handler.logger.Warn().Msgf("Failed to create a new moderation review for comment '%d' by user '%s': moderation review has already been created by that user", moderation.CommentID, details.AuthenticatedUser.Username)
			return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, fmt.Sprintf("Moderation review by this user for comment '%d' already exists", id)), nil
		} else if pqErr.Code == pq.ErrorCode("23514") && pqErr.Constraint == "comment_reviews_approval" {
			handler.logger.Warn().Msgf("User '%s' tried to create an approving moderation review and mark it for deletion")
			return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Cannot make approving moderation review and mark the comment for deletion at the same time"), nil
		} else {
			return http.StatusInternalServerError, nil, err
		}
	} else if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusCreated, moderation, nil
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
