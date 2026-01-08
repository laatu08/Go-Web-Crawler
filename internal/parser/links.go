package parser

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks parses HTML and returns normalized URLs
func ExtractLinks(body io.Reader, baseURL *url.URL, sameDomain bool) []string {
	var links []string

	doc, err := html.Parse(body)
	if err != nil {
		return links
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := normalize(attr.Val, baseURL, sameDomain)
					if link != "" {
						links = append(links, link)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)
	return links
}

func sameSite(host, baseHost string) bool {
	return host == baseHost || strings.HasSuffix(host, "."+baseHost)
}

func isValidWikiArticle(u *url.URL) bool {
	if !strings.HasPrefix(u.Path, "/wiki/") {
		return false
	}

	// Exclude special namespaces
	blacklist := []string{
		":",
		"/wiki/Main_Page",
	}

	for _, b := range blacklist {
		if strings.Contains(u.Path, b) {
			return false
		}
	}

	return true
}

func normalize(href string, base *url.URL, sameDomain bool) string {
	href = strings.TrimSpace(href)

	// Ignore empty, fragments, mailto, javascript
	if href == "" ||
		strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "javascript:") {
		return ""
	}

	u, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// Resolve relative URLs
	abs := base.ResolveReference(u)

	if strings.Contains(base.Host, "wikipedia.org") {
		if !isValidWikiArticle(abs) {
			return ""
		}
	}

	// Enforce same-domain crawling if required
	// if sameDomain && abs.Host != base.Host {
	// 	return ""
	// }

	if sameDomain && !sameSite(abs.Host, base.Host) {
		return ""
	}

	abs.Fragment = "" // Remove fragments
	return abs.String()
}
