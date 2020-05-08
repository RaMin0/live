package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	filenamePattern       *regexp.Regexp
	renamedFilenameFormat string
)

func main() {
	var (
		flagDir     = flag.String("dir", ".", "The path to the directory to work with")
		flagPattern = flag.String("pattern", ".*", "The pattern of files to match")
		flagFormat  = flag.String("format", "%s", "The format of the renamed file")
	)
	flag.Parse()

	filenamePattern = regexp.MustCompile(*flagPattern)
	renamedFilenameFormat = *flagFormat

	if err := filepath.Walk(*flagDir, walkFn); err != nil {
		log.Fatalf("Failed to walk %q: %v", *flagDir, err)
	}
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	fileExtension := filepath.Ext(info.Name())
	filename := strings.TrimSuffix(info.Name(), fileExtension)
	stringMatches := filenamePattern.FindStringSubmatch(filename)
	if len(stringMatches) == 0 {
		return nil
	}
	matches := make([]interface{}, len(stringMatches)-1)
	for i := 1; i < len(stringMatches); i++ {
		matches[i-1] = stringMatches[i]
	}
	renamedFilename := fmt.Sprintf(renamedFilenameFormat, matches...) + fileExtension
	renamedPath := filepath.Join(filepath.Dir(path), renamedFilename)

	return os.Rename(path, renamedPath)
}
