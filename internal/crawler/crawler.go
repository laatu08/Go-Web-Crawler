package crawler

import (
	"net/url"

	"concurrent-web-crawler/internal/config"
	"concurrent-web-crawler/internal/fetcher"
)

type Crawler struct {
	cfg     config.Config
	jobs    chan Job
	visited *VisitedStore
	wg      *WaitGroupWrapper
	fetcher *fetcher.Fetcher
	baseURL *url.URL
}

// NewCrawler initializes the crawler
func NewCrawler(cfg config.Config) (*Crawler, error) {
	base, err := url.Parse(cfg.SeedURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		cfg:     cfg,
		jobs:    make(chan Job, 100),
		visited: NewVisitedStore(cfg.MaxPages),
		wg:      &WaitGroupWrapper{},
		fetcher: fetcher.New(cfg.Timeout),
		baseURL: base,
	}, nil
}

// Start begins crawling
func (c *Crawler) Start() {
	// Start workers
	for i := 1; i <= c.cfg.Workers; i++ {
		worker := NewWorker(
			c.cfg,
			c.fetcher,
			c.visited,
			c.baseURL,
			c.jobs,
			c.wg,
		)
		go worker.Run(i)
	}

	// Seed URL
	if c.visited.TryVisit(c.cfg.SeedURL) {
		c.wg.Add(1)
		c.jobs <- Job{
			URL:   c.cfg.SeedURL,
			Depth: 0,
		}
	}

	// Wait until all jobs finish
	c.wg.Wait()
	close(c.jobs)
}
