package crawler

import (
	"fmt"
	"os"
)

// ExportDOT exports the crawl graph in Graphviz DOT format
func (g *CrawlGraph) ExportDOT(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, "digraph CrawlGraph {")
	if err != nil {
		return err
	}

	edges := g.Edges()
	for from, tos := range edges {
		for _, to := range tos {
			fmt.Fprintf(file, "  \"%s\" -> \"%s\";\n", from, to)
		}
	}

	_, err = fmt.Fprintln(file, "}")
	return err
}
