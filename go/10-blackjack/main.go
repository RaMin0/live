package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ramin0/live/go/blackjack/deck"
)

type playType int

const (
	playTypeHit playType = iota
	playTypeStand
)

func main() {
	cards := deck.New(deck.OptionShuffle())

	players := []player{
		&humanPlayer{},
		&dealerPlayer{},
	}

	deal := func(p player) {
		var c deck.Card
		c, cards = cards[0], cards[1:]
		p.Deal(c)
	}

Game:
	for {
		// 1. Every player is dealt 2 cards
		for i := 1; i <= 2; i++ {
			for _, p := range players {
				deal(p)
			}
		}

		// 2. The player's turn
		// 3. The dealer's turn
	Round:
		for g := 1; ; g++ {
			for _, p := range players {
				fmt.Printf("%s: %v, Score: %d\n",
					p.Name(), p.Hand(), calcScore(p.Hand()))
			}

			// 4. Determining the winner
			// TODO: Fix for handling more than 2 players
			for i := 0; i < len(players); i++ {
				p := players[i]
				s := calcScore(p.Hand())
				if s == 21 {
					fmt.Printf("%s won!\n", p.Name())
					break Game
				}
				if s > 21 {
					otherPlayer := players[(i+1)%len(players)]
					fmt.Printf("%s won!\n", otherPlayer.Name())
					break Game
				}
			}

			for i := 0; i < len(players); {
				p := players[i]

				fmt.Printf("%s's Turn\n", p.Name())
				if _, ok := p.(*dealerPlayer); ok {
					if g == 1 {
						continue
					}
				}

				input, err := p.Play()
				if err != nil {
					log.Fatal(err)
				}

				switch input {
				case playTypeHit:
					deal(p)
					continue Round
				case playTypeStand:
					i++
				}
			}
		}
	}
}

func calcScore(cards []deck.Card) (score int) {
	var aces int
	for _, c := range cards {
		switch s := int(c.Value); {
		case deck.ValueTwo <= c.Value && c.Value <= deck.ValueTen:
			score += s
		case c.Value == deck.ValueAce:
			aces++
			score++
			fallthrough // TODO: Remove this
		case deck.ValueJack <= c.Value && c.Value <= deck.ValueKing:
			score += 10
		}
	}

	// TODO: Decide how to calculate the aces (0 or 10)

	return score
}

type player interface {
	Name() string
	Hand() []deck.Card
	Deal(deck.Card)
	Play() (playType, error)
}

type basePlayer struct {
	hand []deck.Card
}

func (p basePlayer) Hand() []deck.Card {
	return p.hand
}
func (p *basePlayer) Deal(c deck.Card) {
	p.hand = append(p.hand, c)
}

type humanPlayer struct {
	basePlayer
}

func (humanPlayer) Name() string { return "Player" }
func (p humanPlayer) Play() (playType, error) {
	for {
		fmt.Print("(H)it or (S)tand? ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			return 0, err
		}

		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			continue
		}
		switch input[:1] {
		case "h":
			return playTypeHit, nil
		case "s":
			return playTypeStand, nil
		}
	}
}

type dealerPlayer struct {
	basePlayer
}

func (dealerPlayer) Name() string { return "Dealer" }
func (p dealerPlayer) Play() (playType, error) {
	if calcScore(p.hand) <= 16 {
		return playTypeHit, nil
	}
	return playTypeStand, nil
}
