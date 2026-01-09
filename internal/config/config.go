package config

import "time"

type Config struct {
	SeedURL    string        // starting url
	Workers    int           // number of concurrent worker
	MaxDepth   int           // how deep to crawl
	MaxPages   int           // safety limit
	Timeout    time.Duration // HTTP timeout
	SameDomain bool          // crawl only same domain
	RateLimit  int           // requests per second
	UserAgent  string
	OutputFile string
}

// DefaultConfig returns defaults
func DefaultConfig(seed string) Config {
	return Config{
		SeedURL:    seed,
		Workers:    5,
		MaxDepth:   2,
		MaxPages:   100,
		Timeout:    10 * time.Second,
		SameDomain: true,
		RateLimit:  2,
		UserAgent:  "GoWebCrawler/1.0 (+https://example.com/bot-info)",
		OutputFile: "crawl_graph.dot",
	}
}
