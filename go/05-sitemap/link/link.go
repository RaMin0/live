package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func Parse(r io.Reader) ([]Link, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	as := make(chan *html.Node)
	go findAnchors(root, as)
	var links []Link
	for a := range as {
		l := Link{
			Text: extractText(a),
			Href: extractHref(a),
		}
		links = append(links, l)
	}
	return links, nil
}

func findAnchors(n *html.Node, as chan *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		as <- n
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findAnchors(c, as)
	}

	if n.Parent == nil {
		close(as)
	}
}

func extractText(a *html.Node) string {
	var text string
	for c := a.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
			continue
		}
		text += extractText(c)
	}
	return strings.TrimSpace(text)
}

func extractHref(a *html.Node) string {
	for _, attr := range a.Attr {
		if attr.Key != "href" {
			continue
		}
		return attr.Val
	}
	return ""
}
