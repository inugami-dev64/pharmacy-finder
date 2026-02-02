package service

import (
	"maps"
	"pharmafinder/db/dto"
	"pharmafinder/db/entity"
	"pharmafinder/types"
	"slices"
	"sync"
	"time"
)

// Thread safe scraper result collector
// which collects information about pharmacy scrape results
type ScraperStatCollector interface {
	CollectScrapeResult(chain entity.PharmacyChain, success bool)
	GetResults() []dto.PharmacyScraperResultDTO
}

type ScraperStatCollectorImpl struct {
	results map[string]dto.PharmacyScraperResultDTO
	m       sync.RWMutex
}

func ProvideScraperStatCollector() ScraperStatCollector {
	return &ScraperStatCollectorImpl{
		results: make(map[string]dto.PharmacyScraperResultDTO),
		m:       sync.RWMutex{},
	}
}

func (c *ScraperStatCollectorImpl) CollectScrapeResult(chain entity.PharmacyChain, success bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if v, ok := c.results[string(chain)]; ok {
		v.LastSuccessfulScrapeTimestamp = &v.Timestamp
		v.Timestamp = types.Time(time.Now().UTC())
		v.Success = success
		c.results[string(chain)] = v
	} else {
		res := dto.PharmacyScraperResultDTO{
			Chain:     chain,
			Timestamp: types.Time(time.Now().UTC()),
			Success:   success,
		}
		c.results[string(chain)] = res
	}
}

func (c *ScraperStatCollectorImpl) GetResults() []dto.PharmacyScraperResultDTO {
	c.m.RLock()
	defer c.m.RUnlock()

	return slices.Collect(maps.Values(c.results))
}
