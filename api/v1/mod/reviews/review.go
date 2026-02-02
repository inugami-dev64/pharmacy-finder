package reviews

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/service"
	"pharmafinder/types"
	"pharmafinder/utils"
	"pharmafinder/web"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type ModeratorReviewController struct {
	userRepo     db.ModeratorUserRepository
	reviewRepo   db.PharmacyReviewRepository
	tokenManager service.SessionManager
	logger       zerolog.Logger
}

func ProvideModeratorReviewController(
	userRepo db.ModeratorUserRepository,
	reviewRepo db.PharmacyReviewRepository,
	tokenManager service.SessionManager,
) []web.Route {
	controller := &ModeratorReviewController{
		userRepo:     userRepo,
		reviewRepo:   reviewRepo,
		tokenManager: tokenManager,
		logger:       utils.GetLogger("API"),
	}

	return controller.GetRoutes()
}

func (handler *ModeratorReviewController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewSecureRequestsHandler[ModeratorReviewController](handler.GetModerationReviews, "/mod/reviews", []string{"GET"}, web.NewSecurityChain[web.EmptyBody]().RuleAuthenticated(handler.userRepo, handler.tokenManager)),
	}
}

// Get paged resultset for moderated reviews
//
// Path: `GET /api/v1/mod/reviews`
//
// @Summary			Query reviews for moderation
// @Description		Endpoint for querying a paged resultset of reviews for moderation
// @Tags			Moderation
// @Produce			json
// @Param			uk query int false "ID of the latest review in previous query set"
// @Param			k query int false "Timestamp of the latest review in previous query set (unix millis)"
// @Param			l query int false "Limit of the query set (defaults to 50)"
// @Param			desc query boolean false "Reverse the order of reviews (default false)"
// @Param 			unmoderated query boolean false "Show unmoderated reviews (default true)"
// @Param			moderated query boolean false "Show moderated reviews (default: false)"
// @Security		Bearer
// @Success 		200 {array} dto.ModerationPharmacyReview
// @Failure			400 {object} types.HttpError
// @Router			/api/v1/mod/reviews [get]
func (handler *ModeratorReviewController) GetModerationReviews(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	showUnmoderated := true
	showModerated := false

	if unmoderated, err := strconv.ParseBool(details.Params.Get("unmoderated")); err == nil {
		showUnmoderated = unmoderated
	}
	if moderated, err := strconv.ParseBool(details.Params.Get("moderated")); err == nil {
		showModerated = moderated
	}

	ukStr, kStr, l, desc := db.ExtractPagerQueryParameters(details.Params)
	uk, _ := strconv.ParseInt(ukStr, 10, 64)
	k, _ := strconv.ParseInt(kStr, 10, 64)

	var reviews []dto.ModerationPharmacyReview
	var err error
	if uk == 0 || k == 0 {
		reviews, err = handler.reviewRepo.FindReviewsForModeration(showUnmoderated, showModerated).Page(nil, nil, l, desc)
	} else {
		reviews, err = handler.reviewRepo.FindReviewsForModeration(showUnmoderated, showModerated).Page(uk, types.Time(time.UnixMilli(k)), l, desc)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, reviews, nil
}
