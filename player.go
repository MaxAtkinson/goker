package goker

import "fmt"

// Player represents a poker player with hole cards.
type Player struct {
	Name      string
	HoleCards []Card
}

// NewPlayer creates a new player with the given name.
func NewPlayer(name string) *Player {
	return &Player{
		Name:      name,
		HoleCards: make([]Card, 0, 2),
	}
}

// SetHoleCards sets the player's hole cards (must be exactly 2).
func (p *Player) SetHoleCards(cards []Card) error {
	if len(cards) != 2 {
		return ErrInvalidHoleCards
	}
	p.HoleCards = make([]Card, 2)
	copy(p.HoleCards, cards)
	return nil
}

// String returns a string representation of the player.
func (p *Player) String() string {
	return fmt.Sprintf("<Player: %s>", p.Name)
}
