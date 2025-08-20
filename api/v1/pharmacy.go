package v1

import (
	"log"
	"net/http"
	"pharmafinder/db"
	"pharmafinder/types"
	"pharmafinder/utils"
	"strconv"
	"strings"
	"time"
)

type PharmaciesHandler struct {
	repo db.PharmacyRepository
}

func NewPharmaciesHandler(repo db.PharmacyRepository) *PharmaciesHandler {
	return &PharmaciesHandler{repo: repo}
}

func (handler *PharmaciesHandler) Pattern() string {
	return "/pharmacies"
}

func (handler *PharmaciesHandler) Methods() []string {
	return []string{"GET"}
}

// GET /api/v1/phamacies?sw=lat,lng&ne=lat,lng
func (handler *PharmaciesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	swText := query.Get("sw")
	neText := query.Get("ne")

	swCoords := strings.Split(swText, ",")
	neCoords := strings.Split(neText, ",")

	if len(swCoords) != 2 || len(neCoords) != 2 {
		log.Println("Could not extract latitude and longitude from bounds")

		utils.WriteJSONResponse(w, http.StatusBadRequest, types.HttpError{
			StatusCode: http.StatusBadRequest,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "Missing coordinate bounds",
		})
		return
	}

	lat, err := strconv.ParseFloat(swCoords[0], 64)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, types.HttpError{
			StatusCode: http.StatusBadRequest,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "South-west bound latitude is malformed",
		})
		return
	}
	lng, err := strconv.ParseFloat(swCoords[1], 64)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, types.HttpError{
			StatusCode: http.StatusBadRequest,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "South-west bound longitude is malformed",
		})
	}
	sw := types.Point{Lat: float32(lat), Lng: float32(lng)}

	lat, err = strconv.ParseFloat(neCoords[0], 64)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, types.HttpError{
			StatusCode: http.StatusBadRequest,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "North-east bound latitude is malformed",
		})
		return
	}
	lng, err = strconv.ParseFloat(neCoords[1], 64)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, types.HttpError{
			StatusCode: http.StatusBadRequest,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "North-east bound longitude is malformed",
		})
	}
	ne := types.Point{Lat: float32(lat), Lng: float32(lng)}

	data, err := handler.repo.FindPharmaciesInCoordinateBounds(sw, ne).QueryAll()
	if err != nil {
		log.Println("Failed to query pharmacies in coordinate bounds")
		utils.WriteJSONResponse(w, http.StatusInternalServerError, types.HttpError{
			StatusCode: http.StatusInternalServerError,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "Internal server error",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}
