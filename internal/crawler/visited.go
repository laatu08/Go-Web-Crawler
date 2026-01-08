package crawler

import "sync"

// VisitedStore keeps track of visited URLs safely
type VisitedStore struct {
	mu       sync.Mutex
	visited  map[string]struct{}
	maxPages int
	count    int
}

func NewVisitedStore(maxPages int) *VisitedStore {
	return &VisitedStore{
		visited:  make(map[string]struct{}),
		maxPages: maxPages,
	}
}

// TryVisit checks and marks a URL as visited atomically
// Returns false if already visited or maxPages reached
func (v *VisitedStore) TryVisit(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, exists := v.visited[url]; exists {
		return false
	}

	if v.count >= v.maxPages {
		return false
	}

	v.visited[url] = struct{}{}
	v.count++
	return true
}

func (v *VisitedStore) Count() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.count
}
