package crawler

import "time"

// RateLimiter controls request rate globally
type RateLimiter struct {
	ticker *time.Ticker
}

// NewRateLimiter allows `rate` requests per second
func NewRateLimiter(rate int) *RateLimiter {
	interval := time.Second / time.Duration(rate)
	return &RateLimiter{
		ticker: time.NewTicker(interval),
	}
}

// Acquire blocks until a token is available
func (r *RateLimiter) Acquire() {
	<-r.ticker.C
}

func (r *RateLimiter) Stop() {
	r.ticker.Stop()
}
