package crawler

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/temoto/robotstxt"
)

type RobotsChecker struct {
	mu     sync.Mutex
	cache  map[string]*robotstxt.RobotsData
	client *http.Client
	agent  string
}

func NewRobotsChecker(agent string) *RobotsChecker {
	return &RobotsChecker{
		cache: make(map[string]*robotstxt.RobotsData),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		agent: agent,
	}
}

func (r *RobotsChecker) Allowed(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	host := u.Scheme + "://" + u.Host

	r.mu.Lock()
	data, exists := r.cache[host]
	r.mu.Unlock()

	if !exists {
		resp, err := r.client.Get(host + "/robots.txt")
		if err != nil {
			// If robots.txt is unreachable, default to ALLOW
			return true
		}
		defer resp.Body.Close()

		data, err = robotstxt.FromResponse(resp)
		if err != nil {
			return true
		}

		r.mu.Lock()
		r.cache[host] = data
		r.mu.Unlock()
	}

	group := data.FindGroup(r.agent)
	return group.Test(u.Path)
}
