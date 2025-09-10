package reviews

import (
	"crypto/rand"
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

// `POST /api/v1/pharmacies/{id}/reviews`
func (handler *PharmacyReviewController) PostPharmacyReview(details *web.HttpRequestDetails[dto.PharmacyReviewCreationDTO]) (int, interface{}, error) {
	idStr := details.PathVars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handler.logger.Warn().Msgf("Malformed ID path variable '%s'", idStr)
		return http.StatusBadRequest, types.NewHttpError(http.StatusBadRequest, "Malformed ID path variable"), nil
	}

	review := entity.PharmacyReview{
		PharmacyID:       id,
		PrescriptionType: details.Body.PrescriptionType,
		Stars:            details.Body.Stars,
		HRTKind:          details.Body.HRTKind,
		Nationality:      details.Body.Nationality,
		Review:           details.Body.Review,
		CreatedAt:        types.Time(time.Now().UTC()),
		UpdatedAt:        types.Time(time.Now().UTC()),
		ModificationCode: handler.generateModificationCode(),
	}

	err = handler.repo.Store(&review)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, review, nil
}
