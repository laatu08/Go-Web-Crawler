package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"concurrent-web-crawler/internal/config"
	"concurrent-web-crawler/internal/crawler"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
Concurrent Web Crawler

Usage:
  crawler [flags] <seed-url>

Example:
  crawler --workers 8 --rate 2 --depth 3 https://en.wikipedia.org/wiki/India

Flags:
`)
		flag.PrintDefaults()
	}
	// ---------- CLI flags ----------
	workers := flag.Int("workers", 0, "number of concurrent workers")
	rate := flag.Int("rate", 0, "requests per second (global)")
	depth := flag.Int("depth", -1, "maximum crawl depth")
	maxPages := flag.Int("max-pages", 0, "maximum pages to crawl")
	timeout := flag.Int("timeout", 0, "HTTP timeout in seconds")

	sameDomain := flag.Bool("same-domain", true, "restrict crawling to the same domain")
	userAgent := flag.String("user-agent", "", "custom User-Agent string")
	output := flag.String("output", "", "output file for crawl graph (DOT)")

	flag.Parse()

	// ---------- Positional args ----------
	if flag.NArg() < 1 {
		fmt.Println("Usage: crawler <seed-url> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	seedURL := flag.Arg(0)

	// ---------- Defaults ----------
	cfg := config.DefaultConfig(seedURL)

	// ---------- Overrides ----------
	if *workers > 0 {
		cfg.Workers = *workers
	}
	if *rate > 0 {
		cfg.RateLimit = *rate
	}
	if *depth >= 0 {
		cfg.MaxDepth = *depth
	}
	if *maxPages > 0 {
		cfg.MaxPages = *maxPages
	}
	if *timeout > 0 {
		cfg.Timeout = time.Duration(*timeout) * time.Second
	}

	// NEW overrides
	cfg.SameDomain = *sameDomain

	if *userAgent != "" {
		cfg.UserAgent = *userAgent
	}

	if *output != "" {
		cfg.OutputFile = *output
	}

	// ---------- Validation ----------
	if cfg.Workers <= 0 {
		log.Fatal("workers must be > 0")
	}
	if cfg.RateLimit <= 0 {
		log.Fatal("rate limit must be > 0")
	}
	if cfg.MaxPages <= 0 {
		log.Fatal("max-pages must be > 0")
	}

	// ---------- Start crawler ----------
	c, err := crawler.NewCrawler(cfg)
	if err != nil {
		log.Fatalf("failed to initialize crawler: %v", err)
	}

	fmt.Println("🚀 Starting crawl:", cfg.SeedURL)
	fmt.Printf(
		"[CONFIG] workers=%d rate=%d depth=%d maxPages=%d timeout=%s sameDomain=%v output=%s\n",
		cfg.Workers,
		cfg.RateLimit,
		cfg.MaxDepth,
		cfg.MaxPages,
		cfg.Timeout,
		cfg.SameDomain,
		cfg.OutputFile,
	)

	c.Start()
	fmt.Println("✅ Crawl finished")
}
