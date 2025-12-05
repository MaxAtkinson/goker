package goker_test

import (
	"fmt"

	"github.com/MaxAtkinson/goker"
)

func ExampleNewCard() {
	ace := goker.NewCard(goker.Ace, goker.Spades)
	king := goker.NewCard(goker.King, goker.Hearts)

	fmt.Println(ace)
	fmt.Println(king)
	// Output:
	// A♠
	// K♥
}

func ExampleNewHand() {
	hand, _ := goker.NewHand([]goker.Card{
		goker.NewCard(goker.Ace, goker.Spades),
		goker.NewCard(goker.King, goker.Spades),
		goker.NewCard(goker.Queen, goker.Spades),
		goker.NewCard(goker.Jack, goker.Spades),
		goker.NewCard(goker.Ten, goker.Spades),
	})

	fmt.Println(hand.Rank())
	// Output: Royal Flush
}

func ExampleHand_Compare() {
	pair, _ := goker.NewHand([]goker.Card{
		goker.NewCard(goker.Ace, goker.Spades),
		goker.NewCard(goker.Ace, goker.Hearts),
		goker.NewCard(goker.King, goker.Diamonds),
		goker.NewCard(goker.Queen, goker.Clubs),
		goker.NewCard(goker.Jack, goker.Spades),
	})

	highCard, _ := goker.NewHand([]goker.Card{
		goker.NewCard(goker.Ace, goker.Spades),
		goker.NewCard(goker.King, goker.Hearts),
		goker.NewCard(goker.Nine, goker.Diamonds),
		goker.NewCard(goker.Five, goker.Clubs),
		goker.NewCard(goker.Two, goker.Spades),
	})

	fmt.Println(pair.Beats(highCard))
	// Output: true
}

func ExampleNewDeck() {
	deck := goker.NewDeck()
	fmt.Println(deck.Len())

	card, _ := deck.Draw()
	fmt.Println(deck.Len())
	_ = card // use card

	// Output:
	// 52
	// 51
}

func ExampleNewGame() {
	game := goker.NewGame(2)

	// Check players have hole cards
	for _, player := range game.Players {
		fmt.Printf("%s has %d hole cards\n", player.Name, len(player.HoleCards))
	}

	// Output:
	// Player 1 has 2 hole cards
	// Player 2 has 2 hole cards
}
