package crawler

import (
	"log"
	"net/url"

	"concurrent-web-crawler/internal/config"
	"concurrent-web-crawler/internal/fetcher"
	"concurrent-web-crawler/internal/parser"
)

// Job represents a crawl task
type Job struct {
	URL   string
	Depth int
}

type Worker struct {
	cfg         config.Config
	fetcher     *fetcher.Fetcher
	visited     *VisitedStore
	baseURL     *url.URL
	jobs        chan Job
	wg          *WaitGroupWrapper
	rateLimiter *RateLimiter
}

// NewWorker creates a worker instance
func NewWorker(
	cfg config.Config,
	f *fetcher.Fetcher,
	v *VisitedStore,
	base *url.URL,
	jobs chan Job,
	wg *WaitGroupWrapper,
	rl *RateLimiter,
) *Worker {
	return &Worker{
		cfg:         cfg,
		fetcher:     f,
		visited:     v,
		baseURL:     base,
		jobs:        jobs,
		wg:          wg,
		rateLimiter: rl,
	}
}

// Run starts the worker loop
func (w *Worker) Run(id int) {
	for job := range w.jobs {
		w.process(job, id)
		w.wg.Done()
	}
}

func (w *Worker) process(job Job, id int) {
	if job.Depth > w.cfg.MaxDepth {
		return
	}

	// log.Printf("[Worker %d] Crawling: %s (depth=%d)", id, job.URL, job.Depth)

	w.rateLimiter.Acquire()

	body, err := w.fetcher.Fetch(job.URL)
	if err != nil {
		log.Printf("[Worker %d] fetch error: %s (%v)", id, job.URL, err)
		return
	}

	links := parser.ExtractLinks(body, w.baseURL, w.cfg.SameDomain)
	accepted := 0
	for _, link := range links {
		if !w.visited.TryVisit(link, job.Depth+1) {
			continue
		}
		accepted++
		w.wg.Add(1)
		w.jobs <- Job{URL: link, Depth: job.Depth + 1}
	}

	// log.Printf(
	// 	"[Worker %d] Extracted=%d | Accepted=%d | Skipped=%d | %s",
	// 	id,
	// 	len(links),
	// 	accepted,
	// 	len(links)-accepted,
	// 	job.URL,
	// )

	for _, link := range links {
		if !w.visited.TryVisit(link, job.Depth+1) {
			continue
		}

		w.wg.Add(1)
		w.jobs <- Job{
			URL:   link,
			Depth: job.Depth + 1,
		}
	}
}
