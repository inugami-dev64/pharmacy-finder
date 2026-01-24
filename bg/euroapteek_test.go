package bg_test

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"pharmafinder/bg"
	"pharmafinder/db/entity"
	"pharmafinder/mock"
	"pharmafinder/types"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:embed _embeds/euroapteek.html
var euroapteekHtml embed.FS

var euroapteekPharmacies map[int64]entity.Pharmacy = map[int64]entity.Pharmacy{
	int64(-7238096502453610823): {
		PharmacyID:  -7238096502453610823,
		Chain:       string(entity.CHAIN_EUROAPTEEK),
		Name:        "Liivaku Apteek",
		Address:     "J. Sütiste tee 28",
		City:        "Tallinn",
		County:      "Harjumaa",
		PostalCode:  "13411",
		PhoneNumber: "+37282820101",
		ModTime:     types.Time(time.UnixMilli(0)),
		Latitude:    59.397629,
		Longitude:   24.69058,
	},
	int64(-2073188454510069133): {
		PharmacyID:  -2073188454510069133,
		Chain:       string(entity.CHAIN_EUROAPTEEK),
		Name:        "Nõmme Tee Apteek",
		Address:     "Nõmme tee 23a",
		City:        "Tallinn",
		County:      "Harjumaa",
		PostalCode:  "11311",
		PhoneNumber: "+37282820102",
		ModTime:     types.Time(time.UnixMilli(0)),
		Latitude:    59.41785,
		Longitude:   24.72165,
	},
}

func TestEuroapteekScraper_EmptyDB(t *testing.T) {
	// make a local copy of the map
	m2 := make(map[int64]entity.Pharmacy)
	for k := range euroapteekPharmacies {
		m2[k] = euroapteekPharmacies[k]
	}

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		Times(3).
		DoAndReturn(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "www.omniva.ee" {
				search := req.URL.Query().Get("search")
				switch search {
				case "J. Sütiste tee 28, Tallinn, Harjumaa":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"J. Sütiste tee 28, Mustamäe linnaosa, Tallinn, Harju maakond, 13411","zipCode":"13411"}]}`)),
					}, nil
				case "Nõmme tee 23a, Tallinn, Harjumaa":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"Nõmme tee 23a, Kristiine linnaosa, Tallinn, Harju maakond, 11311","zipCode":"11311"}]}`)),
					}, nil
				}
			}

			file, _ := euroapteekHtml.Open("_embeds/euroapteek.html")
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(file),
			}, nil
		})

	queryMock := mock.NewMockQuery[entity.Pharmacy](ctrl)
	queryMock.EXPECT().
		QueryAll().
		Return(nil, nil)

	repoMock := mock.NewMockPharmacyRepository(ctrl)
	repoMock.EXPECT().
		FindPharmaciesByChain(gomock.Eq(entity.CHAIN_EUROAPTEEK)).
		Return(queryMock)
	repoMock.EXPECT().
		StoreAll(gomock.Any()).
		DoAndReturn(func(pharmacies []entity.Pharmacy) error {
			assert.Equal(t, 2, len(pharmacies))

			for _, pharmacy := range pharmacies {
				if v, ok := m2[pharmacy.PharmacyID]; ok {
					assert.Equal(t, v, pharmacy)
					delete(m2, pharmacy.PharmacyID)
				} else {
					assert.Fail(t, fmt.Sprintf("No pharmacy with ID %v was found", pharmacy.PharmacyID))
				}
			}

			assert.Equal(t, 0, len(m2))
			return nil
		})

	scraper := bg.ProvideEuroapteekScraper(repoMock, httpMock)
	scraper.Scrape()
}

func TestEuroapteekScraper_Existing(t *testing.T) {
	// make a local copy of the map
	m2 := make(map[int64]entity.Pharmacy)
	for k := range euroapteekPharmacies {
		pharmacy := euroapteekPharmacies[k]
		if k == -7238096502453610823 {
			pharmacy.ID = 1
		}
		m2[k] = pharmacy
	}

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		Times(2).
		DoAndReturn(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "www.omniva.ee" {
				search := req.URL.Query().Get("search")
				switch search {
				case "J. Sütiste tee 28, Tallinn, Harjumaa":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"J. Sütiste tee 28, Mustamäe linnaosa, Tallinn, Harju maakond, 13411","zipCode":"13411"}]}`)),
					}, nil
				case "Nõmme tee 23a, Tallinn, Harjumaa":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"Nõmme tee 23a, Kristiine linnaosa, Tallinn, Harju maakond, 11311","zipCode":"11311"}]}`)),
					}, nil
				}
			}

			file, _ := euroapteekHtml.Open("_embeds/euroapteek.html")
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(file),
			}, nil
		})

	queryMock := mock.NewMockQuery[entity.Pharmacy](ctrl)
	queryMock.EXPECT().
		QueryAll().
		Return([]entity.Pharmacy{m2[-7238096502453610823]}, nil)

	repoMock := mock.NewMockPharmacyRepository(ctrl)
	repoMock.EXPECT().
		FindPharmaciesByChain(gomock.Eq(entity.CHAIN_EUROAPTEEK)).
		Return(queryMock)
	repoMock.EXPECT().
		StoreAll(gomock.Any()).
		DoAndReturn(func(pharmacies []entity.Pharmacy) error {
			assert.Equal(t, 1, len(pharmacies))

			for _, pharmacy := range pharmacies {
				if v, ok := m2[pharmacy.PharmacyID]; ok {
					assert.Equal(t, v, pharmacy)
					delete(m2, pharmacy.PharmacyID)
				} else {
					assert.Fail(t, "No pharmacy with ID %v was found", pharmacy.PharmacyID)
				}
			}

			assert.Equal(t, 1, len(m2))
			return nil
		})

	scraper := bg.ProvideEuroapteekScraper(repoMock, httpMock)
	scraper.Scrape()
}
