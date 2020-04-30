package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const defaultProblemsFilename = "problems.csv"

var (
	correctAnswers int
	totalQuestions int
)

func main() {
	// define flags for problems filename and quiz time
	var (
		flagProblemsFilename = flag.String("p", defaultProblemsFilename, "The path to the problems CSV file")
		flagTimer            = flag.Duration("t", 30*time.Second, "The max time for the quiz")
		flagShuffle          = flag.Bool("s", false, "Shuffle the quiz questions")
	)
	// parse the value passed from the command line to the flags above
	flag.Parse()

	// make sure the flags are not nil
	if flagProblemsFilename == nil ||
		flagTimer == nil ||
		flagShuffle == nil {
		fmt.Println("Missing problems filename and/or timer")
		return
	}

	// print the values of the flags
	fmt.Printf("Hit enter to start quiz from %q in %v?",
		*flagProblemsFilename, *flagTimer)
	// wait for the user to hit enter before starting
	fmt.Scanln()

	// open the problems file
	f, err := os.Open(*flagProblemsFilename)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	// make sure to close it when we're done
	defer f.Close()

	// read csv file (problems.csv) using csv package
	r := csv.NewReader(f)
	questions, err := r.ReadAll()
	if *flagShuffle {
		// shuffle the questions
		fmt.Println("Shuffling...")
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}
	totalQuestions = len(questions)
	if err != nil {
		fmt.Printf("failed to read csv file: %v\n", err)
		return
	}

	// start the quiz
	quizDone := startQuiz(questions)
	// define the timer
	quizTimer := time.NewTimer(*flagTimer).C

	// wait for quiz timer or quiz is done
	select {
	case <-quizDone:
	case <-quizTimer:
	}

	// output number of questions (total + correct)
	fmt.Printf("Result: %d/%d\n", correctAnswers, totalQuestions)
}

func startQuiz(questions [][]string) chan bool {
	done := make(chan bool)
	go func() {
		// print the questions
		for i, record := range questions {
			question, correctAnswer := record[0], record[1]
			// display one question at a time
			fmt.Printf("%d. %s?\n", i+1, question)
			var answer string
			// get answer from user, then proceed to next question
			// immediately
			if _, err := fmt.Scan(&answer); err != nil {
				fmt.Printf("failed to scan: %v\n", err)
				return
			}
			// clean up answer
			// - by removing extra white space
			answer = strings.TrimSpace(answer)
			// - lower casing the answers to avoid capitalization
			answer = strings.ToLower(answer)
			if answer == correctAnswer {
				correctAnswers++
			}
		}
		// notify the main thread that we're done running the quiz
		done <- true
	}()
	return done
}
