package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ramin0/live/go/twitter/twitter"
)

var (
	consumerKey    = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
)

func main() {
	var (
		flagStatusID  = flag.String("id", "", "The ID of the status to get the retweeting users for.")
		flagUsersFile = flag.String("users", "users.txt", "The name of the file to store the retweeting users in.")
		// flagPickWinner = flag.Bool("pick", false, "Pick a winner.")
		flagPoll     = flag.Duration("poll", 5*time.Minute, "Duration to wait between polls of users.")
		flagDeadline = flag.Duration("deadline", 30*time.Minute, "Deadline for polling users and picking a winner.")
	)
	flag.Parse()

	// emits an event in the channel every x duration
	pollTicker := time.NewTicker(*flagPoll).C
	// emits only one event in the channel after x duration
	deadlineTimer := time.NewTimer(*flagDeadline).C

	// indefinitely ...
	for {
		// .. wait for ...
		select {
		// ... the ticker to tick ...
		case <-pollTicker:
			if *flagStatusID == "" {
				log.Fatalf("Missing status ID")
			}

			users, err := fetchUsers(*flagStatusID)
			if err != nil {
				log.Fatalf("Failed to fetch users: %v", err)
			}

			n, err := persistUsers(users, *flagUsersFile)
			if err != nil {
				log.Fatalf("Failed to persist users: %v", err)
			}
			log.Printf("Added %d users to %s", n, *flagUsersFile)

		// ... or timer to end
		case <-deadlineTimer:
			winner, err := pickWinner(*flagUsersFile)
			if err != nil {
				log.Fatalf("Failed to pick a winner from %s: %v", *flagUsersFile, err)
			}
			log.Printf("The winner is %s", winner)
			return
		}
	}
}

func fetchUsers(statusID string) ([]string, error) {
	c := twitter.NewClient(consumerKey, consumerSecret)

	statuses, err := c.StatusRetweets(statusID)
	if err != nil {
		return nil, err
	}

	users := make([]string, len(statuses))
	for i, s := range statuses {
		users[i] = s.User.ScreenName
	}
	return users, nil
}

func persistUsers(users []string, usersFile string) (int, error) {
	// open or create (O_CREATE) the file if it doesn't exist
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// use a map to ensure uniqueness of users
	usersMap := map[string]bool{}

	// read the users from the file
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		user := sc.Text()
		// add the users to the map
		usersMap[user] = true
	}
	if sc.Err(); err != nil {
		return 0, err
	}

	// add new users from the api passed as an argument to the function
	var usersNum int
	for _, user := range users {
		// if the user already exists in the file, skip it
		if usersMap[user] {
			continue
		}
		// otherwise, add the user to the file
		fmt.Fprintln(f, user)
		usersNum++
	}
	return usersNum, nil
}

func pickWinner(usersFile string) (string, error) {
	f, err := os.Open(usersFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var users []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		user := sc.Text()
		users = append(users, user)
	}
	if sc.Err(); err != nil {
		return "", err
	}

	// ensure randomness on each run of the program
	rand.Seed(int64(time.Now().Nanosecond()))
	// generate a random number between 0 and len(users)-1
	n := rand.Intn(len(users))
	return users[n], nil
}
