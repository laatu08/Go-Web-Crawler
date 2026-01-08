package main

import (
	"fmt"
	"log"
	"os"

	"concurrent-web-crawler/internal/config"
	"concurrent-web-crawler/internal/crawler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crawler <seed-url>")
		os.Exit(1)
	}

	seedURL := os.Args[1]

	cfg := config.DefaultConfig(seedURL)

	c, err := crawler.NewCrawler(cfg)
	if err != nil {
		log.Fatalf("failed to initialize crawler: %v", err)
	}

	fmt.Println("🚀 Starting crawl:", cfg.SeedURL)
	c.Start()
	fmt.Println("✅ Crawl finished")
}
