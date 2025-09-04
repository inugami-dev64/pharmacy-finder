package bg

import (
	"context"

	"github.com/robfig/cron"
	"go.uber.org/fx"
)

// Empty object to represent CronJobs
// mainly used for periodical pharmacy data scraping
type CronJob struct{}

func NewCronJob(scrapers []Scraper, lc fx.Lifecycle) CronJob {
	c := cron.New()

	// Run scrapers on server startup
	for _, scraper := range scrapers {
		scraper.Scrape()
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			for _, scraper := range scrapers {
				c.AddFunc("0 0 1 1,6 *", scraper.Scrape)
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			c.Stop()
			return nil
		},
	})

	return CronJob{}
}
