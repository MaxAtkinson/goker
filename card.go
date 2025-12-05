package goker

import "fmt"

// CardSuit represents the suit of a playing card.
type CardSuit int

const (
	Clubs CardSuit = iota
	Diamonds
	Hearts
	Spades
)

func (s CardSuit) String() string {
	switch s {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}

// CardRank represents the rank of a playing card (2-14, where 14 is Ace).
type CardRank int

const (
	Two CardRank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r CardRank) String() string {
	switch r {
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return fmt.Sprintf("%d", r)
	}
}

// Card represents a single playing card with a rank and suit.
type Card struct {
	Rank CardRank
	Suit CardSuit
}

// NewCard creates a new Card with the given rank and suit.
func NewCard(rank CardRank, suit CardSuit) Card {
	return Card{Rank: rank, Suit: suit}
}

// Value returns the numeric value of the card (2-14).
func (c Card) Value() int {
	return int(c.Rank)
}

// String returns a string representation like "A♠".
func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

// Equal checks if two cards have the same rank (ignoring suit).
func (c Card) Equal(other Card) bool {
	return c.Rank == other.Rank
}

// Less returns true if this card's rank is less than the other's.
func (c Card) Less(other Card) bool {
	return c.Rank < other.Rank
}

// AllSuits returns all possible card suits.
func AllSuits() []CardSuit {
	return []CardSuit{Clubs, Diamonds, Hearts, Spades}
}

// AllRanks returns all possible card ranks.
func AllRanks() []CardRank {
	return []CardRank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
}
