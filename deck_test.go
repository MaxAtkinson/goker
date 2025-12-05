package goker

import "testing"

func TestNewDeck(t *testing.T) {
	deck := NewDeck()
	if deck.Len() != 52 {
		t.Errorf("NewDeck() has %d cards, want 52", deck.Len())
	}
}

func TestDeckDraw(t *testing.T) {
	deck := NewDeck()
	card, err := deck.Draw()
	if err != nil {
		t.Errorf("Draw() error = %v", err)
	}
	if deck.Len() != 51 {
		t.Errorf("After Draw(), deck has %d cards, want 51", deck.Len())
	}
	if card.Rank < Two || card.Rank > Ace {
		t.Errorf("Draw() returned invalid card rank: %v", card.Rank)
	}
}

func TestDeckDrawMany(t *testing.T) {
	deck := NewDeck()
	cards, err := deck.DrawMany(5)
	if err != nil {
		t.Errorf("DrawMany() error = %v", err)
	}
	if len(cards) != 5 {
		t.Errorf("DrawMany(5) returned %d cards, want 5", len(cards))
	}
	if deck.Len() != 47 {
		t.Errorf("After DrawMany(5), deck has %d cards, want 47", deck.Len())
	}
}

func TestDeckBurn(t *testing.T) {
	deck := NewDeck()
	err := deck.Burn()
	if err != nil {
		t.Errorf("Burn() error = %v", err)
	}
	if deck.Len() != 51 {
		t.Errorf("After Burn(), deck has %d cards, want 51", deck.Len())
	}
}

func TestDeckEmptyDraw(t *testing.T) {
	deck := NewDeck()
	// Draw all cards
	for i := 0; i < 52; i++ {
		_, _ = deck.Draw()
	}

	_, err := deck.Draw()
	if err != ErrEmptyDeck {
		t.Errorf("Draw() from empty deck error = %v, want ErrEmptyDeck", err)
	}
}

func TestDeckRemaining(t *testing.T) {
	deck := NewDeck()
	remaining := deck.Remaining()
	if len(remaining) != 52 {
		t.Errorf("Remaining() returned %d cards, want 52", len(remaining))
	}

	// Store original first card
	originalFirst := remaining[0]

	// Modify remaining - pick a card that's definitely different
	var differentCard Card
	if originalFirst.Rank == Ace && originalFirst.Suit == Spades {
		differentCard = NewCard(Two, Clubs)
	} else {
		differentCard = NewCard(Ace, Spades)
	}
	remaining[0] = differentCard

	// Get remaining again - should still have original card
	newRemaining := deck.Remaining()
	if newRemaining[0] != originalFirst {
		t.Error("Remaining() should return a copy, not the original slice")
	}
}

func TestDeckShuffle(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()

	// Very unlikely to be in same order after shuffle
	same := true
	for i := 0; i < 52; i++ {
		c1, _ := deck1.Draw()
		c2, _ := deck2.Draw()
		if c1 != c2 {
			same = false
			break
		}
	}

	// This test has a tiny probability of false positive
	// but practically never happens
	if same {
		t.Error("Two shuffled decks should not be identical")
	}
}
