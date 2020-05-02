package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ramin0/live/go/sitemap/link"
)

func main() {
	flagURL := flag.String("url", "", "The URL to create a sitemap for.")
	flagDepth := flag.Int("depth", 2, "The depth of the links tree.")
	flagXMLFilename := flag.String("xml", "sitemap.xml", "The name of the sitemap XML file.")
	flag.Parse()

	if *flagURL == "" {
		log.Fatal("Missing -url flag")
	}

	sitemap, err := buildSitemap(*flagURL, *flagDepth)
	if err != nil {
		log.Fatalf("Failed to build sitemap for %s: %v", *flagURL, err)
	}

	if err := generateSitemap(sitemap, *flagXMLFilename); err != nil {
		log.Fatalf("Failed to generate sitemap in %s: %v", *flagXMLFilename, err)
	}

	log.Printf("Generated sitemap with %d link(s) for %s in %s", len(sitemap), *flagURL, *flagXMLFilename)
}

func buildSitemap(baseURL string, depth int) ([]string, error) {
	// use a map to ensure uniqueness of parsed urls
	urlsMap := map[string]bool{}

	urls := []string{baseURL}
	for d := 0; d < depth; d++ {
		var newURLs []string
		for _, url := range urls {
			subURLs, err := getURLs(url)
			if err != nil {
				return nil, err
			}
			var uniqueSubURLs []string
			for _, subURL := range subURLs {
				// if we parsed this url before, skip it
				if urlsMap[subURL] {
					continue
				}
				// otherwise, add it to the unique urls map
				urlsMap[subURL] = true
				uniqueSubURLs = append(uniqueSubURLs, subURL)
			}
			newURLs = append(newURLs, uniqueSubURLs...)
		}
		urls = newURLs
	}
	return urls, nil
}

func getURLs(pageURL string) ([]string, error) {
	pageURL = strings.TrimSuffix(pageURL, "/")

	// fetch the html page for this url
	res, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// parse the page and get all links
	links, err := link.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	// we only care about the href of the links
	var urls []string
	for _, l := range links {
		urls = append(urls, l.Href)
	}

	// filter out non-domain links
	var domainURLs []string
	for _, url := range urls {
		// http://google.com or https://ramin0.me
		if strings.HasPrefix(url, "http") &&
			!strings.HasPrefix(url, pageURL) {
			continue
		}
		// https://domain.com/...
		if strings.HasPrefix(url, pageURL) {
			domainURLs = append(domainURLs, url)
			continue
		}
		// mailto:email@example.com
		if strings.Contains(url, ":") {
			continue
		}
		// remove # portion of the url
		// https://domain.com/path/to/page#some-link
		if i := strings.Index(url, "#"); i != -1 {
			url = url[:i]
		}
		// prefix with a / if not already there
		if url == "" || url[0] != '/' {
			url = "/" + url
		}
		// convert /path/to/page to https://domain.com/path/to/page
		url = pageURL + url
		domainURLs = append(domainURLs, url)
	}
	return domainURLs, nil
}

// Generated using https://www.onlinetool.io/xmltogo
type SitemapXML struct {
	XMLName xml.Name        `xml:"urlset"`
	Xmlns   string          `xml:"xmlns,attr"`
	URLs    []SitemapXMLURL `xml:"url"`
}
type SitemapXMLURL struct {
	Loc string `xml:"loc"`
}

func generateSitemap(urls []string, pathToXML string) error {
	var sitemap SitemapXML
	sitemap.Xmlns = "https://www.sitemaps.org/schemas/sitemap/0.9"
	for _, url := range urls {
		sitemap.URLs = append(sitemap.URLs, SitemapXMLURL{
			Loc: url,
		})
	}

	// Alternative: xml.NewEncoder(f).Encode(&sitemap)
	sitemapBytes, err := xml.MarshalIndent(&sitemap, "", "\t")
	if err != nil {
		return err
	}
	xmlData := []byte(xml.Header + string(sitemapBytes))
	return ioutil.WriteFile(pathToXML, xmlData, os.ModePerm)
}
