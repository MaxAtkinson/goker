package goker

import "testing"

func TestNewBoard(t *testing.T) {
	board := NewBoard()
	if board.State() != Preflop {
		t.Errorf("NewBoard state = %v, want Preflop", board.State())
	}
	if len(board.Cards) != 0 {
		t.Errorf("NewBoard cards = %d, want 0", len(board.Cards))
	}
}

func TestBoardSetFlop(t *testing.T) {
	board := NewBoard()
	flop := []Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	}

	err := board.SetFlop(flop)
	if err != nil {
		t.Errorf("SetFlop() error = %v", err)
	}
	if board.State() != Flop {
		t.Errorf("After SetFlop, state = %v, want Flop", board.State())
	}
	if len(board.Flop()) != 3 {
		t.Errorf("Flop() = %d cards, want 3", len(board.Flop()))
	}
}

func TestBoardSetFlopInvalid(t *testing.T) {
	board := NewBoard()

	// Wrong number of cards
	err := board.SetFlop([]Card{NewCard(Ace, Spades)})
	if err != ErrInvalidBoardState {
		t.Errorf("SetFlop() error = %v, want ErrInvalidBoardState", err)
	}

	// Set flop twice
	board.SetFlop([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	})
	err = board.SetFlop([]Card{
		NewCard(Jack, Clubs),
		NewCard(Ten, Spades),
		NewCard(Nine, Hearts),
	})
	if err != ErrInvalidBoardState {
		t.Errorf("SetFlop() twice error = %v, want ErrInvalidBoardState", err)
	}
}

func TestBoardSetTurn(t *testing.T) {
	board := NewBoard()
	board.SetFlop([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	})

	err := board.SetTurn(NewCard(Jack, Clubs))
	if err != nil {
		t.Errorf("SetTurn() error = %v", err)
	}
	if board.State() != Turn {
		t.Errorf("After SetTurn, state = %v, want Turn", board.State())
	}
	if board.TurnCard() == nil {
		t.Error("TurnCard() should not be nil")
	}
}

func TestBoardSetTurnInvalid(t *testing.T) {
	board := NewBoard()

	// Set turn without flop
	err := board.SetTurn(NewCard(Jack, Clubs))
	if err != ErrInvalidBoardState {
		t.Errorf("SetTurn() without flop error = %v, want ErrInvalidBoardState", err)
	}
}

func TestBoardSetRiver(t *testing.T) {
	board := NewBoard()
	board.SetFlop([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	})
	board.SetTurn(NewCard(Jack, Clubs))

	err := board.SetRiver(NewCard(Ten, Spades))
	if err != nil {
		t.Errorf("SetRiver() error = %v", err)
	}
	if board.State() != River {
		t.Errorf("After SetRiver, state = %v, want River", board.State())
	}
	if board.RiverCard() == nil {
		t.Error("RiverCard() should not be nil")
	}
}

func TestBoardSetRiverInvalid(t *testing.T) {
	board := NewBoard()
	board.SetFlop([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	})

	// Set river without turn
	err := board.SetRiver(NewCard(Ten, Spades))
	if err != ErrInvalidBoardState {
		t.Errorf("SetRiver() without turn error = %v, want ErrInvalidBoardState", err)
	}
}

func TestBoardString(t *testing.T) {
	board := NewBoard()
	board.SetFlop([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
	})

	str := board.String()
	if str == "" {
		t.Error("Board.String() should not be empty")
	}
}

func TestBoardStateString(t *testing.T) {
	tests := []struct {
		state    BoardState
		expected string
	}{
		{Preflop, "Preflop"},
		{Flop, "Flop"},
		{Turn, "Turn"},
		{River, "River"},
		{BoardState(99), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("BoardState(%d).String() = %s, want %s", tt.state, got, tt.expected)
		}
	}
}

func TestBoardFlopBeforeSet(t *testing.T) {
	board := NewBoard()
	if board.Flop() != nil {
		t.Error("Flop() should be nil before flop is set")
	}
}

func TestBoardTurnBeforeSet(t *testing.T) {
	board := NewBoard()
	if board.TurnCard() != nil {
		t.Error("TurnCard() should be nil before turn is set")
	}
}

func TestBoardRiverBeforeSet(t *testing.T) {
	board := NewBoard()
	if board.RiverCard() != nil {
		t.Error("RiverCard() should be nil before river is set")
	}
}

func TestBoardStateWithInvalidCardCount(t *testing.T) {
	board := NewBoard()
	// Manually add 1 card (invalid state)
	board.Cards = append(board.Cards, NewCard(Ace, Spades))
	if board.State() != Preflop {
		t.Errorf("Board with 1 card should return Preflop, got %v", board.State())
	}
}
