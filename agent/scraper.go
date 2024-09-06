package agent

import (
	"context"
	"fmt"
	"time"
)

type ScrapeOptions struct {
	Host            string
	IntervalSeconds int
	Channel         chan<- []byte
}

type Scraper struct {
	Options ScrapeOptions

	client *SturdyClient
}

func NewScraper(opts ScrapeOptions) *Scraper {
	client := NewSturdyHTTPClient()
	client.SetBaseURL(opts.Host)

	return &Scraper{opts, client}
}

func (s *Scraper) Scrape(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(s.Options.IntervalSeconds) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := s.scrapeTarget(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (s *Scraper) scrapeTarget(ctx context.Context) error {
	// scrape the target
	response, err := s.client.R().
		SetContext(ctx).
		Get("/metrics")
	if err != nil {
		return fmt.Errorf("failed performing request to get metrics: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("error getting metrics: %s", response.Body())
	}

	s.Options.Channel <- response.Body()

	return nil
}
