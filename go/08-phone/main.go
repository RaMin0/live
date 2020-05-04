package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

var (
	originalPhoneNumbers = []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
)

func main() {
	db, err := sql.Open("sqlite3", "phone_numbers.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS phone_numbers (phone_number TEXT)`,
	); err != nil {
		log.Fatalf("Failed to exec statement: %v", err)
	}

	if _, err := db.Exec(
		`DELETE FROM phone_numbers`,
	); err != nil {
		log.Fatalf("Failed to exec statement: %v", err)
	}

	stmt, err := db.Prepare(
		`INSERT INTO phone_numbers (phone_number) VALUES (?)`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	for _, phoneNumber := range originalPhoneNumbers {
		if _, err := stmt.Exec(phoneNumber); err != nil {
			log.Fatalf("Failed to execute statement: %v", err)
		}
	}

	fmt.Println("Original:")
	rows, err := db.Query("SELECT phone_number FROM phone_numbers")
	if err != nil {
		log.Fatalf("Failed to query statement: %v", err)
	}
	defer rows.Close()
	var phoneNumbers []string
	var formattedPhoneNumbers []string
	for rows.Next() {
		var phoneNumber string
		if err := rows.Scan(&phoneNumber); err != nil {
			log.Fatalf("Failed to scan rows: %v", err)
		}
		formattedPhoneNumber := formatPhoneNumber(phoneNumber)
		phoneNumbers = append(phoneNumbers, phoneNumber)
		formattedPhoneNumbers = append(formattedPhoneNumbers, formattedPhoneNumber)
		fmt.Printf("- %s\n", phoneNumber)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to scan statement: %v", err)
	}

	stmtQuery, err := db.Prepare(
		`SELECT 1 FROM phone_numbers WHERE phone_number = ? LIMIT 1`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	stmtUpdate, err := db.Prepare(
		`UPDATE phone_numbers SET phone_number = ? WHERE phone_number = ?`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	stmtDelete, err := db.Prepare(
		`DELETE FROM phone_numbers WHERE phone_number = ?`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	for i, phoneNumber := range phoneNumbers {
		formattedPhoneNumber := formattedPhoneNumbers[i]

		res, err := stmtQuery.Query(formattedPhoneNumber)
		if err != nil {
			log.Fatalf("Failed to query statement: %v", err)
		}
		duplicateFound := res.Next()
		res.Close()

		if duplicateFound {
			// log.Printf("Duplicate found for %s, delete %s", formattedPhoneNumber, phoneNumber)
			if _, err := stmtDelete.Exec(phoneNumber); err != nil {
				log.Fatalf("Failed to execute statement: %v", err)
			}
			continue
		}

		// log.Printf("Update %s with %s", phoneNumber, formattedPhoneNumber)
		if _, err := stmtUpdate.Exec(formattedPhoneNumber, phoneNumber); err != nil {
			log.Fatalf("Failed to exec statement: %v", err)
		}
	}

	fmt.Println("Formatted:")
	rows, err = db.Query("SELECT phone_number FROM phone_numbers ORDER BY phone_number")
	if err != nil {
		log.Fatalf("Failed to query statement: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var phoneNumber string
		if err := rows.Scan(&phoneNumber); err != nil {
			log.Fatalf("Failed to scan rows: %v", err)
		}
		fmt.Printf("- %s\n", phoneNumber)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to scan statement: %v", err)
	}
}

func formatPhoneNumber(phoneNumber string) string {
	// var formattedPhoneNumber string
	// for _, c := range phoneNumber {
	// 	if !('0' <= c && c <= '9') {
	// 		continue
	// 	}
	// 	formattedPhoneNumber += string(c)
	// }
	// return formattedPhoneNumber

	return regexp.
		MustCompile("[^\\d]").
		ReplaceAllString(phoneNumber, "")
}
