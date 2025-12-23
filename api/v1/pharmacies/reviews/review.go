package reviews

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"pharmafinder/types"
	"pharmafinder/utils"
	"pharmafinder/web"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type PharmacyReviewController struct {
	repo   db.PharmacyReviewRepository
	logger zerolog.Logger
}

func ProvidePharmacyReviewController(repo db.PharmacyReviewRepository) []web.Route {
	controller := &PharmacyReviewController{
		repo:   repo,
		logger: utils.GetLogger("API"),
	}
	return controller.GetRoutes()
}

func (handler *PharmacyReviewController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewRequestsHandler[PharmacyReviewController](handler.PostPharmacyReview, "/pharmacies/{id}/reviews", []string{"POST"}),
		web.NewRequestsHandler[PharmacyReviewController](handler.GetPharmacyReviews, "/pharmacies/{id}/reviews", []string{"GET"}),
		web.NewRequestsHandler[PharmacyReviewController](handler.PatchPharmacyReview, "/pharmacies/{pharmaID}/reviews/{reviewID}", []string{"PATCH"}),
		web.NewRequestsHandler[PharmacyReviewController](handler.DeletePharmacyReview, "/pharmacies/{pharmaID}/reviews/{reviewID}", []string{"DELETE"}),
	}
}

func (handler *PharmacyReviewController) generateModificationCode() string {
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	const n = 16
	builder := strings.Builder{}
	builder.Grow(n)

	for range n {
		val, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		idx := int(val.Int64())
		_, _ = builder.WriteRune(alphabet[idx])
	}

	return builder.String()
}

// Get paged resultset for reviews of given pharmacy
//
// Path: `GET /api/v1/pharmacies/{id}/reviews`
//
// @Summary			Query reviews for pharmacy
// @Description		Endpoint for querying paged resultset of reviews for given pharmacy
// @Tags			Reviews
// @Produce 		json
// @Param			id path integer true "Pharmacy ID"
// @Param			uk query int false "ID of the latest review in previous query set"
// @Param			k query int false "Timestamp of the latest review in previous query set (unix millis)"
// @Param			l query int false "Limit of the query set (defaults to 50)"
// @Param			desc query boolean false "Reverse the order of reviews (default false)"
// @Success 		200 {array} dto.PharmacyReviewsetResultDTO
// @Failure			400 {object} types.HttpError
// @Router			/api/v1/pharmacies/{id}/reviews [get]
func (handler *PharmacyReviewController) GetPharmacyReviews(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	ukStr, kStr, l, desc := db.ExtractPagerQueryParameters(details.Params)
	uk, _ := strconv.ParseInt(ukStr, 10, 64)
	k, _ := strconv.ParseInt(kStr, 10, 64)

	var reviews []entity.PharmacyReview
	if uk == 0 || k == 0 {
		reviews, err = handler.repo.FindReviewForPharmacy(id).Page(nil, nil, l, desc)
	} else {
		reviews, err = handler.repo.FindReviewForPharmacy(id).Page(uk, types.Time(time.UnixMilli(k)), l, desc)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	reviewResult := make([]dto.PharmacyReviewsetResultDTO, len(reviews))
	for i := range reviews {
		reviewResult[i] = dto.PharmacyReviewsetResultDTO{
			ID:               reviews[i].ID,
			PrescriptionType: reviews[i].PrescriptionType,
			Stars:            reviews[i].Stars,
			HRTKind:          reviews[i].HRTKind,
			Nationality:      reviews[i].Nationality,
			Review:           reviews[i].Review,
			CreatedAt:        reviews[i].CreatedAt,
			UpdatedAt:        reviews[i].UpdatedAt,
		}
	}

	return http.StatusOK, reviewResult, nil
}

// Create a new review for provided pharmacy
//
// `POST /api/v1/pharmacies/{id}/reviews`
//
// @Summary			Create a new review for given pharmacy
// @Description		Endpoint for creating a new review for given pharmacy
// @Tags			Reviews
// @Accepts 		json
// @Produce 		json
// @Success			201 {object} entity.PharmacyReview
// @Failure			400 {object} types.HttpError
// @Param			request body dto.PharmacyReviewCreationDTO true "Review creation request body"
// @Param			id path int true "Pharmacy ID"
// @Router			/api/v1/pharmacies/{id}/reviews [post]
func (handler *PharmacyReviewController) PostPharmacyReview(details *web.HttpRequestDetails[dto.PharmacyReviewCreationDTO]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	modCode := handler.generateModificationCode()
	h := sha256.New()
	h.Write([]byte(modCode))
	checksum := h.Sum(nil)

	review := entity.PharmacyReview{
		PharmacyID:       id,
		PrescriptionType: details.Body.PrescriptionType,
		Stars:            details.Body.Stars,
		HRTKind:          details.Body.HRTKind,
		Nationality:      details.Body.Nationality,
		Review:           details.Body.Review,
		CreatedAt:        types.Time(time.Now().UTC()),
		UpdatedAt:        types.Time(time.Now().UTC()),
		ModificationCode: hex.EncodeToString(checksum),
	}

	err = handler.repo.Store(&review)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	review.ModificationCode = modCode
	return http.StatusCreated, review, nil
}

// Modify existing pharmacy review
//
// Path: `PATCH /api/v1/pharmacies/{pharmaID}/reviews/{reviewID}`
//
// @Summary				Modify existing pharmacy review
// @Description 		Endpoint for modifying existing pharmacy review by supplying the modification code
// @Tags				Reviews
// @Accepts 			json
// @Produce				json
// @Param				pharmaID path int true "Pharmacy ID"
// @Param				reviewID path int true "ID of the review to modify"
// @Param				request body dto.PharmacyReviewModificationDTO true "Review modififcation request body"
// @Success				200 {object} dto.PharmacyReviewsetResultDTO
// @Failure				400 {object} types.HttpError
// @Router				/api/v1/pharmacies/{pharmaID}/reviews/{reviewID} [patch]
func (handler *PharmacyReviewController) PatchPharmacyReview(details *web.HttpRequestDetails[dto.PharmacyReviewModificationDTO]) (int, interface{}, error) {
	pharmaIDStr := details.PathVars["pharmaID"]
	reviewIDStr := details.PathVars["reviewID"]

	pharmaID, err := strconv.ParseInt(pharmaIDStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed pharmacy ID path variable '%s'", pharmaIDStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed pharmacy ID path variable"), nil
	}

	reviewID, err := strconv.ParseInt(reviewIDStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed review ID path variable '%s'", pharmaIDStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed review ID path variable"), nil

	}

	review, err := handler.repo.FindReviewByID(pharmaID, reviewID).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	} else if review == nil {
		return http.StatusNotFound, types.NewHttpError(http.StatusNotFound, "Not found"), nil
	}

	// check if provided modifcation code matches the one in the database
	h := sha256.New()
	h.Write([]byte(details.Body.ModificationCode))
	checksum := h.Sum(nil)

	if hex.EncodeToString(checksum) != review.ModificationCode {
		return http.StatusForbidden, types.NewHttpError(http.StatusForbidden, "Invalid modification code"), nil
	}

	review.PrescriptionType = details.Body.PrescriptionType
	review.Stars = details.Body.Stars
	review.HRTKind = details.Body.HRTKind
	review.Nationality = details.Body.Nationality
	review.Review = details.Body.Review
	review.UpdatedAt = types.Time(time.Now().UTC())

	err = handler.repo.Store(review)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, dto.PharmacyReviewsetResultDTO{
		ID:               review.ID,
		PrescriptionType: review.PrescriptionType,
		Stars:            review.Stars,
		HRTKind:          review.HRTKind,
		Nationality:      review.Nationality,
		Review:           review.Review,
		CreatedAt:        review.CreatedAt,
		UpdatedAt:        review.UpdatedAt,
	}, nil
}

// Delete existing pharmacy review
//
// Path: `DELETE /api/v1/pharmacies/{pharmaID}/reviews/{reviewID}`
//
// @Summary			Hard-deletes existing pharmacy review
// @Description		Endpoint for hard-deleting existing pharmacy review
// @Tags			Reviews
// @Produce 		json
// @Param			pharmaID path int true "Pharmacy ID"
// @Param			reviewID path int true "ID of the review to delete"
// @Security		Bearer
// @Success 		200 {object} dto.PharmacyReviewsetResultDTO
// @Failure			400 {object} types.HttpError
// @Failure			403	{object} types.HttpError
// @Router			/api/v1/pharmacies/{pharmaID}/reviews/{reviewID} [delete]
func (handler *PharmacyReviewController) DeletePharmacyReview(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	pharmaIDStr := details.PathVars["pharmaID"]
	reviewIDStr := details.PathVars["reviewID"]

	pharmaID, err := strconv.ParseInt(pharmaIDStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed pharmacy ID path variable '%s'", pharmaIDStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed pharmacy ID path variable"), nil
	}

	reviewID, err := strconv.ParseInt(reviewIDStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed review ID path variable '%s'", pharmaIDStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed review ID path variable"), nil

	}

	review, err := handler.repo.FindReviewByID(pharmaID, reviewID).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	} else if review == nil {
		return http.StatusNotFound, types.NewHttpError(http.StatusNotFound, "Not found"), nil
	}

	auth := details.Header.Get("Authorization")
	splits := strings.Split(auth, " ")
	var bearer string
	if len(splits) > 1 {
		bearer = strings.Trim(splits[1], " \t")
	}

	// check if provided modifcation code matches the one in the database
	h := sha256.New()
	h.Write([]byte(bearer))
	checksum := h.Sum(nil)

	if hex.EncodeToString(checksum) != review.ModificationCode {
		return http.StatusForbidden, types.NewHttpError(http.StatusForbidden, "Invalid modification code"), nil
	}

	review, err = handler.repo.Delete(reviewID).Query()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, dto.PharmacyReviewsetResultDTO{
		ID:               review.ID,
		PrescriptionType: review.PrescriptionType,
		Stars:            review.Stars,
		HRTKind:          review.HRTKind,
		Nationality:      review.Nationality,
		Review:           review.Review,
		CreatedAt:        review.CreatedAt,
		UpdatedAt:        review.UpdatedAt,
	}, nil
}
