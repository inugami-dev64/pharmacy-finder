package pharmacies

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

type PharmaciesController struct {
	repo   db.PharmacyRepository
	logger zerolog.Logger
}

func ProvidePharmacyController(repo db.PharmacyRepository) []web.Route {
	controller := &PharmaciesController{
		repo:   repo,
		logger: utils.GetLogger("API"),
	}
	return controller.GetRoutes()
}

func (handler *PharmaciesController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewRequestsHandler[PharmaciesController](handler.GetPharmacies, "/pharmacies", []string{"GET"}),
	}
}

// Pharmacy retriever endpoint
//
// GET /api/v1/phamacies?sw=lat,lng&ne=lat,lng
//
// @Summary			Get all pharmacies in coordinate bounds
// @Description 	Endpoint for querying all pharmacies in specified coordinate bounds
// @Tags			Pharmacy
// @Produce 		json
// @Success 		200 {array} entity.Pharmacy
// @Failure			400 {object} types.HttpError
// @Param			sw query string true "South-west coordinates of the bound, syntax: lat,lng"
// @Param			ne query string true "North-east coordinates of the bound, syntax: lat,lng"
// @Router			/api/v1/pharmacies [get]
func (handler *PharmaciesController) GetPharmacies(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	swText := details.Params.Get("sw")
	neText := details.Params.Get("ne")

	swCoords := strings.Split(swText, ",")
	neCoords := strings.Split(neText, ",")

	if len(swCoords) != 2 || len(neCoords) != 2 {
		handler.logger.Warn().Msg("Could not extract latitude and longitude from bounds")
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Missing coordinate bounds"), nil
	}

	lat, err := strconv.ParseFloat(swCoords[0], 64)
	if err != nil {
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "South-west bound latitude is malformed"), nil
	}
	lng, err := strconv.ParseFloat(swCoords[1], 64)
	if err != nil {
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "South-west bound longitude is malformed"), nil
	}
	sw := types.Point{Lat: float32(lat), Lng: float32(lng)}

	lat, err = strconv.ParseFloat(neCoords[0], 64)
	if err != nil {
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "North-east bound latitude is malformed"), nil
	}
	lng, err = strconv.ParseFloat(neCoords[1], 64)
	if err != nil {
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "North-east bound longitude is malformed"), nil
	}
	ne := types.Point{Lat: float32(lat), Lng: float32(lng)}

	data, err := handler.repo.FindPharmaciesInCoordinateBounds(sw, ne).QueryAll()
	if err != nil {
		handler.logger.Warn().Msgf("Failed to query pharmacies in coordinate bounds")
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, data, nil
}
