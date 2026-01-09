package crawler

import (
	"log"
	"net/url"
	"time"

	"concurrent-web-crawler/internal/config"
	"concurrent-web-crawler/internal/fetcher"
)

type Crawler struct {
	cfg         config.Config
	jobs        chan Job
	visited     *VisitedStore
	wg          *WaitGroupWrapper
	fetcher     *fetcher.Fetcher
	baseURL     *url.URL
	startTime   time.Time
	rateLimiter *RateLimiter
	statsTicker *time.Ticker
	done        chan struct{}
	robots      *RobotsChecker
	graph       *CrawlGraph
}

// NewCrawler initializes the crawler
func NewCrawler(cfg config.Config) (*Crawler, error) {
	base, err := url.Parse(cfg.SeedURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		cfg:         cfg,
		jobs:        make(chan Job, 100),
		visited:     NewVisitedStore(cfg.MaxPages),
		wg:          &WaitGroupWrapper{},
		fetcher:     fetcher.New(cfg.Timeout),
		baseURL:     base,
		rateLimiter: NewRateLimiter(cfg.RateLimit),
		done:        make(chan struct{}),
		robots:      NewRobotsChecker(cfg.UserAgent),
		graph:       NewCrawlGraph(),
	}, nil
}

// Start begins crawling
func (c *Crawler) Start() {
	c.startTime = time.Now()
	c.startProgressLogger()

	// Start workers
	for i := 1; i <= c.cfg.Workers; i++ {
		worker := NewWorker(
			c.cfg,
			c.fetcher,
			c.visited,
			c.baseURL,
			c.jobs,
			c.wg,
			c.rateLimiter,
			c.robots,
			c.graph,
		)

		go worker.Run(i)
	}

	// Seed URL
	if c.visited.TryVisit(c.cfg.SeedURL, 0) {
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
	c.rateLimiter.Stop()
	c.statsTicker.Stop()
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

	log.Println("[INFO] Pages by depth:")
	stats := c.visited.DepthStats()
	for depth := 0; depth <= c.cfg.MaxDepth; depth++ {
		log.Printf("  Depth %d → %d pages", depth, stats[depth])
	}
	c.printGraphStats()

	err := c.graph.ExportDOT("crawl_graph.dot")
	if err != nil {
		log.Printf("[WARN] Failed to export crawl graph: %v", err)
	} else {
		log.Println("[INFO] Crawl graph exported to crawl_graph.dot")
	}

}

func (c *Crawler) startProgressLogger() {
	c.statsTicker = time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-c.statsTicker.C:
				elapsed := time.Since(c.startTime).Seconds()
				log.Printf(
					"[STATS] Crawled=%d | Queue=%d | Elapsed=%.1fs",
					c.visited.Count(),
					len(c.jobs),
					elapsed,
				)
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Crawler) printGraphStats() {
	edges := c.graph.Edges()

	totalEdges := 0
	for _, tos := range edges {
		totalEdges += len(tos)
	}

	log.Printf("[INFO] Crawl graph nodes: %d", len(edges))
	log.Printf("[INFO] Crawl graph edges: %d", totalEdges)
}
