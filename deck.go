package goker

import (
	"math/rand"
)

// Deck represents a deck of playing cards.
type Deck struct {
	cards []Card
}

// NewDeck creates a new shuffled 52-card deck.
func NewDeck() *Deck {
	d := &Deck{
		cards: make([]Card, 0, 52),
	}
	for _, suit := range AllSuits() {
		for _, rank := range AllRanks() {
			d.cards = append(d.cards, NewCard(rank, suit))
		}
	}
	d.Shuffle()
	return d
}

// Shuffle randomizes the order of cards in the deck.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Len returns the number of cards remaining in the deck.
func (d *Deck) Len() int {
	return len(d.cards)
}

// Draw removes and returns the top card from the deck.
func (d *Deck) Draw() (Card, error) {
	if len(d.cards) == 0 {
		return Card{}, ErrEmptyDeck
	}
	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	return card, nil
}

// DrawMany removes and returns multiple cards from the deck.
func (d *Deck) DrawMany(n int) ([]Card, error) {
	cards := make([]Card, n)
	for i := 0; i < n; i++ {
		card, err := d.Draw()
		if err != nil {
			return nil, err
		}
		cards[i] = card
	}
	return cards, nil
}

// Burn discards the top card (used in Texas Hold'em dealing).
func (d *Deck) Burn() error {
	_, err := d.Draw()
	return err
}

// Remaining returns a copy of the remaining cards.
func (d *Deck) Remaining() []Card {
	result := make([]Card, len(d.cards))
	copy(result, d.cards)
	return result
}
