package crawler

import "sync"

// CrawlGraph stores directed edges: from -> to
type CrawlGraph struct {
	mu    sync.Mutex
	edges map[string][]string
}

func NewCrawlGraph() *CrawlGraph {
	return &CrawlGraph{
		edges: make(map[string][]string),
	}
}

// AddEdge records a directed edge
func (g *CrawlGraph) AddEdge(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges[from] = append(g.edges[from], to)
}

func (g *CrawlGraph) Edges() map[string][]string {
	g.mu.Lock()
	defer g.mu.Unlock()

	copy := make(map[string][]string)
	for k, v := range g.edges {
		copy[k] = append([]string{}, v...)
	}
	return copy
}
