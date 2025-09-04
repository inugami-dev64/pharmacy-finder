package bg_test

import (
	"embed"
	"io"
	"net/http"
	"pharmafinder/bg"
	"pharmafinder/db/entity"
	"pharmafinder/mock"
	"pharmafinder/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:embed _embeds/benu.html
var benuHtml embed.FS

var benuPharmacies map[int]entity.Pharmacy = map[int]entity.Pharmacy{
	491: {
		ID:          0,
		PharmacyID:  491,
		Chain:       string(entity.CHAIN_BENU),
		Name:        "Veskimöldre BENU Apteek",
		Address:     "Instituudi tee 132",
		City:        "Saue vald",
		County:      "Harjumaa",
		PostalCode:  "76403",
		Email:       "benu.5656@benu.ee",
		PhoneNumber: "+3726888055",
		ModTime:     types.Time(unwrap(time.Parse("2006-01-02 15:04:05", "2025-07-02 08:36:31"))),
		Latitude:    59.35778,
		Longitude:   24.60182,
	},
	406: {
		ID:          0,
		PharmacyID:  406,
		Chain:       string(entity.CHAIN_BENU),
		Name:        "Kohila apteek",
		Address:     "Lõuna 2",
		City:        "Kohila",
		County:      "Raplamaa",
		PostalCode:  "79804",
		Email:       "kohilaapteek1@gmail.com",
		PhoneNumber: "+3724833574",
		ModTime:     types.Time(unwrap(time.Parse("2006-01-02 15:04:05", "2025-06-02 12:54:48"))),
		Latitude:    59.16742,
		Longitude:   24.74963,
	},
	33: {
		ID:          0,
		PharmacyID:  33,
		Chain:       string(entity.CHAIN_BENU),
		Name:        "Lasnamäe Tervisemaja Apteek",
		Address:     "Linnamäe tee 3, Lasnamäe",
		City:        "Tallinn",
		County:      "Harjumaa",
		PostalCode:  "13912",
		Email:       "benu.5154@benu.ee",
		PhoneNumber: "+3726091998",
		ModTime:     types.Time(unwrap(time.Parse("2006-01-02 15:04:05", "2025-09-02 19:51:29"))),
		Latitude:    59.44924,
		Longitude:   24.86303,
	},
	465: {
		ID:          0,
		PharmacyID:  465,
		Chain:       string(entity.CHAIN_BENU),
		Name:        "Jõhvi Tsentraali apteek",
		Address:     "Keskväljak 4",
		City:        "Jõhvi",
		County:      "Ida-Virumaa",
		PostalCode:  "41531",
		Email:       "benu.5642@benu.ee",
		PhoneNumber: "+3726622001",
		ModTime:     types.Time(unwrap(time.Parse("2006-01-02 15:04:05", "2025-07-07 16:06:07"))),
		Latitude:    59.35835,
		Longitude:   27.41395,
	},
}

func TestBenuScraper_EmptyDB(t *testing.T) {
	// make a local copy of the map
	m2 := make(map[int]entity.Pharmacy)
	for k := range benuPharmacies {
		m2[k] = benuPharmacies[k]
	}

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(req *http.Request) (*http.Response, error) {
			file, _ := benuHtml.Open("_embeds/benu.html")
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
		FindPharmaciesByChain(gomock.Eq(entity.CHAIN_BENU)).
		Return(queryMock)
	repoMock.EXPECT().
		StoreAll(gomock.Any()).
		DoAndReturn(func(pharmacies []entity.Pharmacy) error {
			assert.Equal(t, 4, len(pharmacies))

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

	scraper := bg.ProvideBenuScraper(repoMock, httpMock)
	scraper.Scrape()
}

func TestBenuScraper_Existing(t *testing.T) {
	// make a local copy of the map
	m2 := make(map[int]entity.Pharmacy)
	counter := 1
	existingPharmacies := make([]entity.Pharmacy, 0)
	for k := range benuPharmacies {
		if k != 33 {
			pharmacy := benuPharmacies[k]
			pharmacy.ID = int64(counter)
			counter++
			m2[k] = pharmacy
			existingPharmacies = append(existingPharmacies, pharmacy)
		} else {
			m2[k] = benuPharmacies[k]
		}
	}

	ctrl := gomock.NewController(t)
	httpMock := mock.NewMockHttpClient(ctrl)
	httpMock.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(req *http.Request) (*http.Response, error) {
			file, _ := benuHtml.Open("_embeds/benu.html")
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(file),
			}, nil
		})

	queryMock := mock.NewMockQuery[entity.Pharmacy](ctrl)
	queryMock.EXPECT().
		QueryAll().
		Return(existingPharmacies, nil)

	repoMock := mock.NewMockPharmacyRepository(ctrl)
	repoMock.EXPECT().
		FindPharmaciesByChain(gomock.Eq(entity.CHAIN_BENU)).
		Return(queryMock)
	repoMock.EXPECT().
		StoreAll(gomock.Any()).
		DoAndReturn(func(pharmacies []entity.Pharmacy) error {
			assert.Equal(t, 1, len(pharmacies))
			assert.Equal(t, int64(0), pharmacies[0].ID)
			assert.Equal(t, int64(33), pharmacies[0].PharmacyID)
			return nil
		})

	scraper := bg.ProvideBenuScraper(repoMock, httpMock)
	scraper.Scrape()
}
