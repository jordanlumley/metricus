package agent

import (
	"context"
	"sync"
)

type Options struct {
	ScrapeTargets []ScrapeOptions
}

type Agent struct {
	Options Options
}

func New(opts Options) *Agent {
	return &Agent{
		Options: opts,
	}
}

func (a *Agent) Start() {
	var wg sync.WaitGroup
	for _, target := range a.Options.ScrapeTargets {
		wg.Add(1)
		go func(target ScrapeOptions) {
			defer wg.Done()

			scraper := NewScraper(ScrapeOptions{
				Host:            target.Host,
				IntervalSeconds: target.IntervalSeconds,
				Channel:         target.Channel,
			})

			scraper.Scrape(context.Background())
		}(target)
	}
	wg.Wait()
}
