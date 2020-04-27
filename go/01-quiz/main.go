package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

const defaultProblemsFilename = "problems.csv"

func main() {
	f, err := os.Open(defaultProblemsFilename)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer f.Close()

	// read csv file (problems.csv) using csv package
	r := csv.NewReader(f)
	questions, err := r.ReadAll()
	if err != nil {
		fmt.Printf("failed to read csv file: %v\n", err)
		return
	}

	// fmt.Println(r)
	var correctAnswers int
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
		if answer == correctAnswer {
			correctAnswers++
		}
	}

	// output number of questions (total + corrent)
	fmt.Printf("Result: %d/%d\n", correctAnswers, len(questions))
}
