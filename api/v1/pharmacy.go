package v1

import (
	"log"
	"net/http"
	"pharmafinder/db"
	"pharmafinder/types"
	"pharmafinder/web"
	"strconv"
	"strings"
)

type PharmaciesController struct {
	repo db.PharmacyRepository
}

func ProvidePharmacyController(repo db.PharmacyRepository) []web.Route {
	controller := &PharmaciesController{repo: repo}
	return controller.GetRoutes()
}

func (handler *PharmaciesController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewRequestsHandler[PharmaciesController](handler.GetPharmacies, "/pharmacies", []string{"GET"}),
	}
}

// GET /api/v1/phamacies?sw=lat,lng&ne=lat,lng
func (handler *PharmaciesController) GetPharmacies(details *web.HttpRequestDetails) (int, interface{}, error) {
	swText := details.Params.Get("sw")
	neText := details.Params.Get("ne")

	swCoords := strings.Split(swText, ",")
	neCoords := strings.Split(neText, ",")

	if len(swCoords) != 2 || len(neCoords) != 2 {
		log.Println("Could not extract latitude and longitude from bounds")
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
		log.Println("Failed to query pharmacies in coordinate bounds")
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, data, nil
}
