package bg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"pharmafinder/db"

	"github.com/anaskhan96/soup"
)

const BENU_ENDPOINT = "https://www.benu.ee/leia-apteek"

type BenuScraper struct {
	repo db.PharmacyRepository
}

func ProvideBenuScraper(repo db.PharmacyRepository) Scraper {
	return BenuScraper{repo: repo}
}

func (scraper BenuScraper) Scrape() {
	req, err := http.NewRequest("GET", BENU_ENDPOINT, nil)
	if err != nil {
		log.Println("Failed to create a new request for BENU scraper")
		return
	}
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to make a request to %s: %v", BENU_ENDPOINT, err)
		return
	}

	// make sure that the server responded with status code 200
	if resp.StatusCode != 200 {
		log.Printf("Benu endpoint responded with non-200 status code %d", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body from BENU endpoint request")
		return
	}

	doc := soup.HTMLParse(string(body))
	script := doc.Find("div", "class", "bnContainer").Find("script")
	if script.Error != nil {
		log.Printf("Failed to extract script tag from BENU website's HTML body")
		return
	}

	txt := script.Text()
	fmt.Println(txt)
}
