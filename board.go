package goker

// BoardState represents the current state of the community cards.
type BoardState int

const (
	Preflop BoardState = iota + 1
	Flop
	Turn
	River
)

func (s BoardState) String() string {
	switch s {
	case Preflop:
		return "Preflop"
	case Flop:
		return "Flop"
	case Turn:
		return "Turn"
	case River:
		return "River"
	default:
		return "Unknown"
	}
}

// Board represents the community cards in a Texas Hold'em game.
type Board struct {
	Cards []Card
}

// NewBoard creates a new empty board.
func NewBoard() *Board {
	return &Board{
		Cards: make([]Card, 0, 5),
	}
}

// State returns the current state of the board based on card count.
func (b *Board) State() BoardState {
	switch len(b.Cards) {
	case 0:
		return Preflop
	case 3:
		return Flop
	case 4:
		return Turn
	case 5:
		return River
	default:
		return Preflop
	}
}

// Flop returns the first 3 community cards.
func (b *Board) Flop() []Card {
	if len(b.Cards) >= 3 {
		return b.Cards[:3]
	}
	return nil
}

// SetFlop sets the first 3 community cards.
func (b *Board) SetFlop(cards []Card) error {
	if len(cards) != 3 {
		return ErrInvalidBoardState
	}
	if len(b.Cards) != 0 {
		return ErrInvalidBoardState
	}
	b.Cards = append(b.Cards, cards...)
	return nil
}

// TurnCard returns the 4th community card.
func (b *Board) TurnCard() *Card {
	if len(b.Cards) >= 4 {
		return &b.Cards[3]
	}
	return nil
}

// SetTurn adds the turn card.
func (b *Board) SetTurn(card Card) error {
	if len(b.Cards) != 3 {
		return ErrInvalidBoardState
	}
	b.Cards = append(b.Cards, card)
	return nil
}

// RiverCard returns the 5th community card.
func (b *Board) RiverCard() *Card {
	if len(b.Cards) >= 5 {
		return &b.Cards[4]
	}
	return nil
}

// SetRiver adds the river card.
func (b *Board) SetRiver(card Card) error {
	if len(b.Cards) != 4 {
		return ErrInvalidBoardState
	}
	b.Cards = append(b.Cards, card)
	return nil
}

// String returns a string representation of the board.
func (b *Board) String() string {
	s := "<Board: "
	for i, card := range b.Cards {
		if i > 0 {
			s += " "
		}
		s += card.String()
	}
	s += ">"
	return s
}
