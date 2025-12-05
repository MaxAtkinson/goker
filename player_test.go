package goker

import "testing"

func TestNewPlayer(t *testing.T) {
	player := NewPlayer("Alice")
	if player.Name != "Alice" {
		t.Errorf("Player name = %s, want Alice", player.Name)
	}
	if len(player.HoleCards) != 0 {
		t.Errorf("New player has %d hole cards, want 0", len(player.HoleCards))
	}
}

func TestPlayerSetHoleCards(t *testing.T) {
	player := NewPlayer("Bob")
	cards := []Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
	}

	err := player.SetHoleCards(cards)
	if err != nil {
		t.Errorf("SetHoleCards() error = %v", err)
	}
	if len(player.HoleCards) != 2 {
		t.Errorf("Player has %d hole cards, want 2", len(player.HoleCards))
	}
}

func TestPlayerSetHoleCardsInvalid(t *testing.T) {
	player := NewPlayer("Charlie")

	// Too few cards
	err := player.SetHoleCards([]Card{NewCard(Ace, Spades)})
	if err != ErrInvalidHoleCards {
		t.Errorf("SetHoleCards() with 1 card error = %v, want ErrInvalidHoleCards", err)
	}

	// Too many cards
	err = player.SetHoleCards([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	})
	if err != ErrInvalidHoleCards {
		t.Errorf("SetHoleCards() with 3 cards error = %v, want ErrInvalidHoleCards", err)
	}
}

func TestPlayerString(t *testing.T) {
	player := NewPlayer("Diana")
	expected := "<Player: Diana>"
	if got := player.String(); got != expected {
		t.Errorf("Player.String() = %s, want %s", got, expected)
	}
}
