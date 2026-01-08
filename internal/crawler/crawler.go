package crawler

import (
	"log"
	"net/url"
	"time"

	"concurrent-web-crawler/internal/config"
	"concurrent-web-crawler/internal/fetcher"
)

type Crawler struct {
	cfg       config.Config
	jobs      chan Job
	visited   *VisitedStore
	wg        *WaitGroupWrapper
	fetcher   *fetcher.Fetcher
	baseURL   *url.URL
	startTime time.Time
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
	c.startTime = time.Now()

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

	// Wait for crawl completion
	c.wg.Wait()
	close(c.jobs)

	c.printStats()
}


func (c *Crawler) printStats() {
	elapsed := time.Since(c.startTime).Seconds()

	log.Printf(
		"[INFO] Pages crawled: %d / %d",
		c.visited.Count(),
		c.cfg.MaxPages,
	)

	log.Printf(
		"[INFO] Queue length: %d",
		len(c.jobs),
	)

	log.Printf(
		"[INFO] Crawl complete in %.2fs",
		elapsed,
	)
}
