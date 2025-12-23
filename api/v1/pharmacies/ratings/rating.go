package ratings

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/types"
	"pharmafinder/utils"
	"pharmafinder/web"
	"strconv"

	"github.com/rs/zerolog"
)

type PharmacyRatingController struct {
	repo   db.PharmacyRepository
	logger zerolog.Logger
}

func ProvidePharmacyRatingController(repo db.PharmacyRepository) []web.Route {
	controller := &PharmacyRatingController{
		repo:   repo,
		logger: utils.GetLogger("API"),
	}

	return controller.GetRoutes()
}

func (handler *PharmacyRatingController) GetRoutes() []web.Route {
	return []web.Route{}
}

// Pharmacy ratings endpoint - aggregated average score
//
// Path: `GET /api/v1/pharmacies/{id}/ratings`
//
// @Summary 		Pharmacy ratings endpoint
// @Description		Queries informaton about average ratings for given pharmacy
// @Tags			Ratings
// @Produce 		json
// @Success			200 {array} dto.PharmacyRatingDTO
// @Failure			400 {object} types.HttpError
// @Param			id path string true "Pharmacy ID"
// @Router			/api/v1/pharmacies/{id}/ratings [get]
func (handler *PharmacyRatingController) GetPharmacyRatings(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	ratings, err := handler.repo.FindPharmacyRatingsByID(id).QueryAll()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, ratings, nil
}
