package main

import (
	"fmt"

	"github.com/ramin0/live/go/deck/deck"
)

func main() {
	for _, c := range deck.New(
		deck.OptionSort(func(i, j deck.Card) bool {
			return i.Suit > j.Suit ||
				i.Suit == j.Suit && i.Value < j.Value
		}),
		deck.OptionSort(deck.SortDefault),
		deck.OptionShuffle(),
		deck.OptionAddJokers(3),
		deck.OptionExclude(func(c deck.Card) bool {
			return c.Suit != deck.SuitSpades
		}),
		deck.OptionCompose(
			deck.New(
				deck.OptionExclude(func(c deck.Card) bool {
					return true
				}),
				deck.OptionAddJokers(3),
			),
			deck.New(),
		),
	) {
		fmt.Println(c)
	}
}
