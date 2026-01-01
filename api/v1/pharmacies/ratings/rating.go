package ratings

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/types"
	"pharmafinder/utils"
	"pharmafinder/web"
	"strconv"
	"strings"

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
	return []web.Route{
		web.NewRequestsHandler[PharmacyRatingController](handler.GetAllPharmacyRatings, "/pharmacies/ratings", []string{"GET"}),
		web.NewRequestsHandler[PharmacyRatingController](handler.GetPharmacyRatingsByPharmacy, "/pharmacies/{id}/ratings", []string{"GET"}),
	}
}

// Pharmacy ratings endpoint - aggregated average score
//
// Path: `GET /api/v1/pharmacies/{id}/ratings`
//
// @Summary 		Pharmacy ratings endpoint
// @Description		Queries information about average ratings for given pharmacy
// @Tags			Ratings
// @Produce 		json
// @Success			200 {array} dto.PharmacyRatingDTO
// @Failure			400 {object} types.HttpError
// @Param			id path integer true "Pharmacy ID"
// @Router			/api/v1/pharmacies/{id}/ratings [get]
func (handler *PharmacyRatingController) GetPharmacyRatingsByPharmacy(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
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

// Pharmacy ratings endpoint for all pharmacies - aggregated average score
//
// Path: `GET /api/v1/pharmacies/ratings`
//
// @Summary			Get all pharmacy ratings
// @Description		Queries information about average ratings for all pharmacies in the database
// @Tags			Ratings
// @Produce 		json
// @Success			200 {array} dto.PharmacyTierRatingDTO
// @Failure			400 {object} types.HttpError
// @Param			sw query string false "South-west bound coordinates in 'lat,lng' syntax"
// @Param			ne query string false "North-east bound coordinates in 'lat,lng' syntax"
// @Router 			/api/v1/pharmacies/ratings [get]
func (handler *PharmacyRatingController) GetAllPharmacyRatings(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	swStrCoords := strings.Split(details.Params.Get("sw"), ",")
	neStrCoords := strings.Split(details.Params.Get("ne"), ",")
	sw := types.Point{Lat: -90, Lng: -90}
	ne := types.Point{Lat: 90, Lng: 90}
	if len(swStrCoords) == 2 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(swStrCoords[0]), 64); err == nil {
			sw.Lat = float32(v)
		}
		if v, err := strconv.ParseFloat(strings.TrimSpace(swStrCoords[1]), 64); err == nil {
			sw.Lng = float32(v)
		}
	}

	if len(neStrCoords) == 2 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(neStrCoords[0]), 64); err == nil {
			ne.Lat = float32(v)
		}
		if v, err := strconv.ParseFloat(strings.TrimSpace(neStrCoords[1]), 64); err == nil {
			ne.Lng = float32(v)
		}
	}

	ratings, err := handler.repo.FindPharmacyRatings(sw, ne).QueryAll()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, ratings, nil
}
