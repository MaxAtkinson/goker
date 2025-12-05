package goker

import "testing"

func TestCardRankString(t *testing.T) {
	tests := []struct {
		rank     CardRank
		expected string
	}{
		{Two, "2"},
		{Three, "3"},
		{Nine, "9"},
		{Ten, "T"},
		{Jack, "J"},
		{Queen, "Q"},
		{King, "K"},
		{Ace, "A"},
	}

	for _, tt := range tests {
		if got := tt.rank.String(); got != tt.expected {
			t.Errorf("CardRank(%d).String() = %s, want %s", tt.rank, got, tt.expected)
		}
	}
}

func TestCardSuitString(t *testing.T) {
	tests := []struct {
		suit     CardSuit
		expected string
	}{
		{Clubs, "♣"},
		{Diamonds, "♦"},
		{Hearts, "♥"},
		{Spades, "♠"},
		{CardSuit(99), "?"},
	}

	for _, tt := range tests {
		if got := tt.suit.String(); got != tt.expected {
			t.Errorf("CardSuit(%d).String() = %s, want %s", tt.suit, got, tt.expected)
		}
	}
}

func TestNewCard(t *testing.T) {
	card := NewCard(Ace, Spades)
	if card.Rank != Ace {
		t.Errorf("NewCard rank = %v, want %v", card.Rank, Ace)
	}
	if card.Suit != Spades {
		t.Errorf("NewCard suit = %v, want %v", card.Suit, Spades)
	}
}

func TestCardValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{NewCard(Two, Spades), 2},
		{NewCard(Ten, Hearts), 10},
		{NewCard(Ace, Clubs), 14},
	}

	for _, tt := range tests {
		if got := tt.card.Value(); got != tt.expected {
			t.Errorf("Card(%v).Value() = %d, want %d", tt.card, got, tt.expected)
		}
	}
}

func TestCardString(t *testing.T) {
	card := NewCard(Ace, Spades)
	expected := "A♠"
	if got := card.String(); got != expected {
		t.Errorf("Card.String() = %s, want %s", got, expected)
	}
}

func TestCardEqual(t *testing.T) {
	card1 := NewCard(Ace, Spades)
	card2 := NewCard(Ace, Hearts)
	card3 := NewCard(King, Spades)

	if !card1.Equal(card2) {
		t.Error("Cards with same rank should be equal")
	}
	if card1.Equal(card3) {
		t.Error("Cards with different ranks should not be equal")
	}
}

func TestCardLess(t *testing.T) {
	card1 := NewCard(King, Spades)
	card2 := NewCard(Ace, Spades)

	if !card1.Less(card2) {
		t.Error("King should be less than Ace")
	}
	if card2.Less(card1) {
		t.Error("Ace should not be less than King")
	}
}

func TestAllSuits(t *testing.T) {
	suits := AllSuits()
	if len(suits) != 4 {
		t.Errorf("AllSuits() returned %d suits, want 4", len(suits))
	}
}

func TestAllRanks(t *testing.T) {
	ranks := AllRanks()
	if len(ranks) != 13 {
		t.Errorf("AllRanks() returned %d ranks, want 13", len(ranks))
	}
}
