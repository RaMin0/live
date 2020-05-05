// Package deck can be used to create decks of playing cards
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

// Suit is used to define the suit of a Card
type Suit int

// The list of Suits that can be assigned to a Card
const (
	SuitSpades Suit = iota
	SuitHearts
	SuitDiamonds
	SuitClubs
	SuitJoker
)

func (s Suit) String() string {
	switch s {
	case SuitSpades:
		return "♠"
	case SuitHearts:
		return "♥"
	case SuitDiamonds:
		return "♦"
	case SuitClubs:
		return "♣"
	case SuitJoker:
		return "J"
	}
	return ""
}

// Value is used to define the value of a Card
type Value int

// The list of Values that can be assigned to a Card
const (
	_ Value = iota
	_
	ValueTwo
	ValueThree
	ValueFour
	ValueFive
	ValueSix
	ValueSeven
	ValueEight
	ValueNine
	ValueTen
	ValueJack
	ValueQueen
	ValueKing
	ValueAce
)

func (v Value) String() string {
	switch v {
	case ValueJack:
		return "J"
	case ValueQueen:
		return "Q"
	case ValueKing:
		return "K"
	case ValueAce:
		return "A"
	default:
		return strconv.Itoa(int(v))
	}
}

// Card holds a combination of a Suit and a Value
type Card struct {
	Suit  Suit
	Value Value
}

func (c Card) String() string {
	if c.Suit == SuitJoker {
		return "[JOKER]"
	}
	return fmt.Sprintf("[%v  %-2v]", c.Suit, c.Value)
}

// Option defines a way to manipulate a deck of Cards
type Option func([]Card) []Card

// New creates a new deck with the specified Options
func New(opts ...Option) []Card {
	var deck []Card
	for suit := SuitSpades; suit <= SuitClubs; suit++ {
		for value := ValueTwo; value <= ValueAce; value++ {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	for _, opt := range opts {
		deck = opt(deck)
	}
	return deck
}

// SortDefault provides the default sorting logic for a deck
func SortDefault(i, j Card) bool {
	return i.Suit < j.Suit ||
		i.Suit == j.Suit && i.Value < j.Value
}

// OptionSort can sort a deck based on the sorting function fn
func OptionSort(fn func(Card, Card) bool) Option {
	return func(deck []Card) []Card {
		sort.Slice(deck, func(i, j int) bool {
			return fn(deck[i], deck[j])
		})
		return deck
	}
}

// OptionShuffle shuffles a deck
func OptionShuffle() Option {
	return func(deck []Card) []Card {
		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})
		return deck
	}
}

// OptionAddJokers adds n arbitary Jokers to the end of a deck
func OptionAddJokers(n int) Option {
	return func(deck []Card) []Card {
		for i := 1; i <= n; i++ {
			deck = append(deck, Card{Suit: SuitJoker})
		}
		return deck
	}
}

// OptionExclude uses fn to know which cards to exludes from a deck
func OptionExclude(fn func(Card) bool) Option {
	return func(deck []Card) []Card {
		var newDeck []Card
		for _, c := range deck {
			if fn(c) {
				continue
			}
			newDeck = append(newDeck, c)
		}
		return newDeck
	}
}

// OptionCompose composes a bigger deck by adding other decks to a deck
func OptionCompose(decks ...[]Card) Option {
	return func(deck []Card) []Card {
		for _, d := range decks {
			deck = append(deck, d...)
		}
		return deck
	}
}
