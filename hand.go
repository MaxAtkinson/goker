package goker

import (
	"fmt"
	"sort"
)

const handSize = 5

// Hand represents a 5-card poker hand with evaluation capabilities.
type Hand struct {
	Cards  []Card
	Player *Player

	// Cached evaluation results
	cardValues int
	cardCounts int
	isFlush    bool
	isStraight bool
	isWheel    bool
	handRank   HandRank
	evaluated  bool
}

// NewHand creates a new Hand from 5 cards.
func NewHand(cards []Card) (*Hand, error) {
	if len(cards) != handSize {
		return nil, ErrInvalidHandSize
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, c := range cards {
		key := fmt.Sprintf("%d-%d", c.Rank, c.Suit)
		if seen[key] {
			return nil, ErrDuplicateCards
		}
		seen[key] = true
	}

	// Sort cards by rank descending
	sortedCards := make([]Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Rank > sortedCards[j].Rank
	})

	h := &Hand{
		Cards: sortedCards,
	}
	h.evaluate()
	return h, nil
}

// NewHandWithPlayer creates a hand associated with a player.
func NewHandWithPlayer(cards []Card, player *Player) (*Hand, error) {
	h, err := NewHand(cards)
	if err != nil {
		return nil, err
	}
	h.Player = player
	return h, nil
}

// evaluate computes and caches all hand properties.
func (h *Hand) evaluate() {
	if h.evaluated {
		return
	}

	h.cardValues = h.computeCardValues()
	h.cardCounts = h.computeCardCounts()
	h.isFlush = h.computeIsFlush()
	h.isStraight, h.isWheel = h.computeIsStraight()
	h.handRank = h.computeHandRank()
	h.evaluated = true
}

// computeCardValues creates a bitmap of which ranks are present.
// Bit positions correspond to card ranks: Ace=14, King=13, ..., Two=2
func (h *Hand) computeCardValues() int {
	values := 0
	for _, card := range h.Cards {
		values |= (1 << int(card.Rank))
	}
	return values
}

// computeCardCounts returns counts of each rank occurrence.
func (h *Hand) computeCardCounts() int {
	// Not used with new simplified approach, kept for interface compatibility
	return 0
}

// getRankCounts returns a map of rank to occurrence count.
func (h *Hand) getRankCounts() map[CardRank]int {
	counts := make(map[CardRank]int)
	for _, card := range h.Cards {
		counts[card.Rank]++
	}
	return counts
}

// computeIsFlush checks if all cards have the same suit.
func (h *Hand) computeIsFlush() bool {
	if len(h.Cards) == 0 {
		return false
	}
	suit := h.Cards[0].Suit
	for _, card := range h.Cards[1:] {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

// computeIsStraight checks for straight and wheel (A2345) patterns.
func (h *Hand) computeIsStraight() (isStraight, isWheel bool) {
	// Check for wheel (A2345)
	if h.cardValues == wheelStraightValue {
		return true, true
	}

	// Check for regular straight using LSB technique
	if h.cardValues == 0 {
		return false, false
	}

	// Find LSB (lowest set bit) - always non-zero if cardValues != 0
	lsb := h.cardValues & (-h.cardValues)

	// A straight has 5 consecutive bits, which means cardValues/LSB == 31
	quotient := h.cardValues / lsb
	return quotient == straightValue, false
}

// computeHandRank determines the hand ranking.
func (h *Hand) computeHandRank() HandRank {
	// Check for flush-based hands first
	if h.isFlush {
		if h.cardValues == royalFlushValue {
			return RoyalFlush
		}
		if h.isStraight {
			return StraightFlush
		}
		return Flush
	}

	// Check for straight
	if h.isStraight {
		return Straight
	}

	// Count pairs, trips, and quads
	counts := h.getRankCounts()
	pairs := 0
	trips := 0
	quads := 0

	for _, count := range counts {
		switch count {
		case 2:
			pairs++
		case 3:
			trips++
		case 4:
			quads++
		}
	}

	// Determine hand rank based on counts
	if quads == 1 {
		return FourOfAKind
	}
	if trips == 1 && pairs == 1 {
		return FullHouse
	}
	if trips == 1 {
		return ThreeOfAKind
	}
	if pairs == 2 {
		return TwoPair
	}
	if pairs == 1 {
		return Pair
	}

	return HighCard
}

// Rank returns the evaluated hand rank.
func (h *Hand) Rank() HandRank {
	return h.handRank
}

// IsFlush returns true if the hand is a flush.
func (h *Hand) IsFlush() bool {
	return h.isFlush
}

// IsStraight returns true if the hand is a straight.
func (h *Hand) IsStraight() bool {
	return h.isStraight
}

// IsWheel returns true if the hand is the wheel straight (A2345).
func (h *Hand) IsWheel() bool {
	return h.isWheel
}

// TiebreakScore returns a numeric score for breaking ties between hands of equal rank.
func (h *Hand) TiebreakScore() int {
	// Sort cards by frequency then by value
	type cardWithCount struct {
		card  Card
		count int
	}

	counts := make(map[CardRank]int)
	for _, card := range h.Cards {
		counts[card.Rank]++
	}

	cardsWithCounts := make([]cardWithCount, len(h.Cards))
	for i, card := range h.Cards {
		cardsWithCounts[i] = cardWithCount{card, counts[card.Rank]}
	}

	// Sort by count descending, then by rank descending
	sort.Slice(cardsWithCounts, func(i, j int) bool {
		if cardsWithCounts[i].count != cardsWithCounts[j].count {
			return cardsWithCounts[i].count > cardsWithCounts[j].count
		}
		return cardsWithCounts[i].card.Rank > cardsWithCounts[j].card.Rank
	})

	// Build tiebreak score using bit shifts
	score := 0
	for i, cwc := range cardsWithCounts {
		if i < len(tiebreakerShifts) {
			score |= int(cwc.card.Rank) << tiebreakerShifts[i]
		}
	}

	return score
}

// Compare compares two hands. Returns:
// -1 if h loses to other
//
//	0 if tie
//	1 if h beats other
func (h *Hand) Compare(other *Hand) int {
	if h.handRank > other.handRank {
		return 1
	}
	if h.handRank < other.handRank {
		return -1
	}

	// Same hand rank - need tiebreaker
	// Special case: wheel vs non-wheel straight
	if h.isStraight && other.isStraight {
		if h.isWheel && !other.isWheel {
			return -1
		}
		if !h.isWheel && other.isWheel {
			return 1
		}
	}

	myScore := h.TiebreakScore()
	otherScore := other.TiebreakScore()

	if myScore > otherScore {
		return 1
	}
	if myScore < otherScore {
		return -1
	}
	return 0
}

// Beats returns true if this hand beats the other.
func (h *Hand) Beats(other *Hand) bool {
	return h.Compare(other) > 0
}

// Ties returns true if this hand ties with the other.
func (h *Hand) Ties(other *Hand) bool {
	return h.Compare(other) == 0
}

// Contains checks if the hand contains a card of the given rank.
func (h *Hand) Contains(rank CardRank) bool {
	for _, card := range h.Cards {
		if card.Rank == rank {
			return true
		}
	}
	return false
}

// String returns a string representation of the hand.
func (h *Hand) String() string {
	s := "<Hand: "
	for i, card := range h.Cards {
		if i > 0 {
			s += " "
		}
		s += card.String()
	}
	s += ">"
	return s
}
