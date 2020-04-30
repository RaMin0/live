package main

import (
	"encoding/json"
	"flag"
	"fmt"
	htmltemplate "html/template"
	"net/http"
	"os"
	texttemplate "text/template"
)

type Story map[string]struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func main() {
	var (
		flagStoryJSONFilename = flag.String("story", "gopher.json", "The path to the JSON for the story to render")
		flagHTTP              = flag.Bool("http", false, "Run as a web server")
	)
	flag.Parse()

	f, err := os.Open(*flagStoryJSONFilename)
	if err != nil {
		fmt.Printf("Failed to open %s: %v", *flagStoryJSONFilename, err)
		return
	}
	defer f.Close()

	var story Story
	if err := json.NewDecoder(f).Decode(&story); err != nil {
		fmt.Printf("Failed to decode: %v", err)
		return
	}

	if *flagHTTP {
		fmt.Println("Listening on 8080...")
		http.ListenAndServe(":8080", NewStoryMux(story))
	} else {
		runAsCmd(story)
	}
}

type StoryMux struct {
	story Story
}

func (m StoryMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arcName := r.URL.Query().Get("arc")
	if arcName == "" {
		arcName = "intro"
	}

	arc, ok := m.story[arcName]
	if !ok {
		http.Error(w, fmt.Sprintf("arc not found: %s", arcName), http.StatusNotFound)
		return
	}

	tmpl, err := htmltemplate.ParseFiles("arc.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse template: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, arc); err != nil {
		http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}
}

func NewStoryMux(story Story) http.Handler {
	return StoryMux{story}
}

func runAsCmd(story Story) {
	arc := story["intro"]
	tmpl := texttemplate.Must(texttemplate.ParseFiles("arc.txt"))

	for {
		tmpl.Execute(os.Stdout, arc)
		if len(arc.Options) == 0 {
			break
		}

		fmt.Printf("Choice: ")
		var choice int
		for {
			if _, err := fmt.Scanf("%d\n", &choice); err != nil {
				fmt.Printf("Failed to scan: %v", err)
				return
			}
			if choice < 0 || choice >= len(arc.Options) {
				fmt.Printf("Invalid choice: %d. Allowed [0-%d]: ", choice, len(arc.Options)-1)
				continue
			}
			break
		}
		arcName := arc.Options[choice].Arc
		arc = story[arcName]
	}
}
