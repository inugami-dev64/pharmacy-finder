package bg_test

import (
	"embed"
	"io"
	"net/http"
	"pharmafinder/bg"
	"pharmafinder/db/entity"
	"pharmafinder/mock"
	"pharmafinder/types"
	"pharmafinder/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:embed _embeds/apotheka.json
var apothekaJson embed.FS

var apothekaPharmacies map[int]entity.Pharmacy = map[int]entity.Pharmacy{
	1: {
		ID:          0,
		PharmacyID:  1,
		Chain:       string(entity.CHAIN_APOTHEKA),
		Name:        "AKADEEMIA KONSUMI APTEEK",
		Address:     "Akadeemia tee 35",
		City:        "Tallinn",
		County:      "Harjumaa",
		PostalCode:  "12618",
		Email:       "kajaapt@apotheka.ee",
		PhoneNumber: "+3726587701",
		ModTime:     types.Time(utils.Unwrap(time.Parse("2006-01-02 15:04:05", "2021-02-01 09:08:38"))),
		Latitude:    59.403729,
		Longitude:   24.655573,
	},
	5: {
		ID:          0,
		PharmacyID:  5,
		Chain:       string(entity.CHAIN_APOTHEKA),
		Name:        "ASTRI KESKUSE APTEEK",
		Address:     "Tallinna mnt 41",
		City:        "Narva",
		County:      "Ida-Virumaa",
		PostalCode:  "20605",
		Email:       "tempharu@apotheka.ee",
		PhoneNumber: "+3723573071",
		ModTime:     types.Time(utils.Unwrap(time.Parse("2006-01-02 15:04:05", "2021-02-01 09:08:38"))),
		Latitude:    59.380785,
		Longitude:   28.174233,
	},
}

func TestApothekaScraper_EmptyDB(t *testing.T) {
	// make a local copy of the map
	m2 := make(map[int]entity.Pharmacy)
	for k := range apothekaPharmacies {
		m2[k] = apothekaPharmacies[k]
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
				case "Akadeemia tee 35, Tallinn, Harju maakond":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"Akadeemia tee 35, Mustamäe linnaosa, Tallinn, Harju maakond, 12618","zipCode":"12618"}]}`)),
					}, nil
				case "Tallinna mnt 41, Narva, Ida-Viru maakond":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"Tallinna mnt 41, Narva linn, Ida-Viru maakond, 20605","zipCode":"20605"}]}`)),
					}, nil
				}
			}

			file, _ := apothekaJson.Open("_embeds/apotheka.json")
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
		FindPharmaciesByChain(gomock.Eq(entity.CHAIN_APOTHEKA)).
		Return(queryMock)
	repoMock.EXPECT().
		StoreAll(gomock.Any()).
		DoAndReturn(func(pharmacies []entity.Pharmacy) error {
			assert.Equal(t, 2, len(pharmacies))

			for _, pharmacy := range pharmacies {
				if v, ok := m2[int(pharmacy.PharmacyID)]; ok {
					assert.Equal(t, v, pharmacy)
					delete(m2, int(pharmacy.PharmacyID))
				} else {
					assert.Fail(t, "No pharmacy with ID %v was found", pharmacy.PharmacyID)
				}
			}

			assert.Equal(t, 0, len(m2))
			return nil
		})

	scraper := bg.ProvideApothekaScraper(repoMock, httpMock)
	scraper.Scrape()
}

func TestApothekaScraper_Existing(t *testing.T) {
	// make a local copy of the map
	m2 := make(map[int]entity.Pharmacy)
	for k := range apothekaPharmacies {
		pharmacy := apothekaPharmacies[k]
		if k == 1 {
			pharmacy.ID = 1
		}
		m2[k] = pharmacy
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
				case "Akadeemia tee 35, Tallinn, Harju maakond":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"Akadeemia tee 35, Mustamäe linnaosa, Tallinn, Harju maakond, 12618","zipCode":"12618"}]}`)),
					}, nil
				case "Tallinna mnt 41, Narva, Ida-Viru maakond":
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"addresses":[{"address":"Tallinna mnt 41, Narva linn, Ida-Viru maakond, 20605","zipCode":"20605"}]}`)),
					}, nil
				}
			}

			file, _ := apothekaJson.Open("_embeds/apotheka.json")
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(file),
			}, nil
		})

	queryMock := mock.NewMockQuery[entity.Pharmacy](ctrl)
	queryMock.EXPECT().
		QueryAll().
		Return([]entity.Pharmacy{m2[1]}, nil)

	repoMock := mock.NewMockPharmacyRepository(ctrl)
	repoMock.EXPECT().
		FindPharmaciesByChain(gomock.Eq(entity.CHAIN_APOTHEKA)).
		Return(queryMock)
	repoMock.EXPECT().
		StoreAll(gomock.Any()).
		DoAndReturn(func(pharmacies []entity.Pharmacy) error {
			assert.Equal(t, 1, len(pharmacies))

			for _, pharmacy := range pharmacies {
				if v, ok := m2[int(pharmacy.PharmacyID)]; ok {
					assert.Equal(t, v, pharmacy)
					delete(m2, int(pharmacy.PharmacyID))
				} else {
					assert.Fail(t, "No pharmacy with ID %v was found", pharmacy.PharmacyID)
				}
			}

			assert.Equal(t, 1, len(m2))
			return nil
		})

	scraper := bg.ProvideApothekaScraper(repoMock, httpMock)
	scraper.Scrape()
}
