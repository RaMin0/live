package deck

import (
	"fmt"
)

func ExampleNew() {
	deck := New(
		OptionShuffle(),
	)
	fmt.Println(deck)
}
