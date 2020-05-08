package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ramin0/live/go/quiet_hn/hn"
)

var (
	cache     = map[int]hn.Item{}
	cacheLock sync.RWMutex
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		// only get the first numStories*1.25 items
		ids = ids[:int(float64(numStories)*1.25)]

		storiesChan := make(chan orderedItem)
		var wg sync.WaitGroup
		for i, id := range ids {
			// add a goroutine for the wg to wait for
			wg.Add(1)
			go func(id, i int) {
				// notify the wg that this goroutine is done
				defer wg.Done()

				if _, ok := cacheRead(id); !ok {
					item, err := client.GetItem(id)
					if err != nil {
						return
					}
					cacheWrite(id, item)
				}

				hnItem, _ := cacheRead(id)
				item := parseHNItem(hnItem)
				if isStoryLink(item) {
					storiesChan <- orderedItem{item, i}
				}
			}(id, i)
		}

		// start a monitoring goroutine to wait for the wg
		go func() {
			wg.Wait()
			close(storiesChan)
		}()

		var stories []orderedItem
		for orderedItem := range storiesChan {
			stories = append(stories, orderedItem)
		}
		sort.Slice(stories, func(i, j int) bool {
			return stories[i].idx < stories[j].idx
		})

		// storiesMap := map[int]orderedItem{}
		// for orderedItem := range storiesChan {
		// 	storiesMap[orderedItem.idx] = orderedItem
		// }
		// var stories []orderedItem
		// for i := 0; i < len(ids); i++ {
		// 	if orderedItem, ok := storiesMap[i]; ok {
		// 		stories = append(stories, orderedItem)
		// 	}
		// }

		data := templateData{
			Stories: stories[:numStories],
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type orderedItem struct {
	item
	idx int
}

type templateData struct {
	Stories []orderedItem
	Time    time.Duration
}

func cacheRead(id int) (hn.Item, bool) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	hnItem, ok := cache[id]
	return hnItem, ok
}

func cacheWrite(id int, hnItem hn.Item) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[id] = hnItem
}
