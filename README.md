# Goker

[![Tests](https://github.com/MaxAtkinson/goker/actions/workflows/test.yml/badge.svg)](https://github.com/MaxAtkinson/goker/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/MaxAtkinson/goker/branch/main/graph/badge.svg)](https://codecov.io/gh/MaxAtkinson/goker)
[![Go Reference](https://pkg.go.dev/badge/github.com/MaxAtkinson/goker.svg)](https://pkg.go.dev/github.com/MaxAtkinson/goker)

A lightweight Go package for Texas Hold'em poker hand evaluation and comparison.

## Installation

```bash
go get github.com/MaxAtkinson/goker
```

## Features

- Card and deck representation
- Hand evaluation (all 10 poker hand rankings)
- Hand comparison with tiebreakers
- Full Texas Hold'em game simulation
- Efficient binary arithmetic evaluation

## Usage

### Basic Card Operations

```go
package main

import (
    "fmt"
    "github.com/MaxAtkinson/goker"
)

func main() {
    // Create cards
    ace := goker.NewCard(goker.Ace, goker.Spades)
    king := goker.NewCard(goker.King, goker.Hearts)

    fmt.Println(ace) // A♠
    fmt.Println(king) // K♥
}
```

### Hand Evaluation

```go
// Create a hand
hand, _ := goker.NewHand([]goker.Card{
    goker.NewCard(goker.Ace, goker.Spades),
    goker.NewCard(goker.King, goker.Spades),
    goker.NewCard(goker.Queen, goker.Spades),
    goker.NewCard(goker.Jack, goker.Spades),
    goker.NewCard(goker.Ten, goker.Spades),
})

fmt.Println(hand.Rank()) // Royal Flush
```

### Hand Comparison

```go
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

fmt.Println(pair.Beats(highCard)) // true
```

### Full Game Simulation

```go
// Create a 4-player game (automatically shuffles and deals hole cards)
game := goker.NewGame(4)

// Deal community cards
game.DealFlop()
game.DealTurn()
game.DealRiver()

// Find winner(s)
winners, hands, _ := game.GetWinners()
for i, winner := range winners {
    fmt.Printf("%s wins with %v (%s)\n", winner.Name, hands[i], hands[i].Rank())
}
```

## Hand Rankings

From lowest to highest:
1. High Card
2. Pair
3. Two Pair
4. Three of a Kind
5. Straight
6. Flush
7. Full House
8. Four of a Kind
9. Straight Flush
10. Royal Flush

## API Reference

### Types

- `Card` - Single playing card
- `CardRank` - Card rank (Two through Ace)
- `CardSuit` - Card suit (Clubs, Diamonds, Hearts, Spades)
- `Deck` - 52-card deck with shuffle/draw operations
- `Hand` - 5-card poker hand with evaluation
- `HandRank` - Poker hand ranking
- `Player` - Player with hole cards
- `Board` - Community cards
- `BoardState` - Preflop, Flop, Turn, River
- `Game` - Complete Texas Hold'em game

### Key Functions

- `NewCard(rank, suit)` - Create a card
- `NewDeck()` - Create shuffled deck
- `NewHand(cards)` - Create and evaluate a 5-card hand
- `NewGame(numPlayers)` - Create a new game

## License

Apache 2.0
