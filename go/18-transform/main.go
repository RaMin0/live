package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/", displayUploadForm)
	http.HandleFunc("/upload", handleUploadForm)
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("tmp"))))
	log.Print("Listening on 3000...")
	http.ListenAndServe(":3000", nil)
}

func displayUploadForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "upload.html")
}

func handleUploadForm(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	mode := r.FormValue("mode")
	num, err := strconv.Atoi(r.FormValue("num"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	img, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("mode=%s, num=%d img=%d", mode, num, len(img))

	cwd := filepath.Join("tmp", uuid.New().String())
	log.Printf("cwd=%s", cwd)
	if err := os.MkdirAll(cwd, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ioutil.WriteFile(filepath.Join(cwd, "img.jpg"), img, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := make([]string, num)
	for n := 1; n <= num; n++ {
		path, err := run(filepath.Join(cwd, "img.jpg"), mode, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results[n-1] = path
	}

	tmpl := template.Must(template.ParseFiles("results.html"))
	if err := tmpl.Execute(w, results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func run(inPath string, mode string, n int) (string, error) {
	outDir := filepath.Dir(inPath)
	outName := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	outName += fmt.Sprintf("%05d", n)
	outName += filepath.Ext(inPath)
	outPath := filepath.Join(outDir, outName)
	cmd := exec.Command("primitive",
		"-i", inPath,
		"-o", outPath,
		"-n", strconv.Itoa(n),
		"-m", mode)
	log.Printf("Run: %s %v", cmd.Path, strings.Join(cmd.Args[1:], " "))
	out, err := cmd.Output()
	if err != nil {
		exitErr := err.(*exec.ExitError)
		return "", fmt.Errorf("%v: %s", err, exitErr.Stderr)
	}
	log.Print(string(out))
	return outPath, nil
}
