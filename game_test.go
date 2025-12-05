package goker

import "testing"

func TestNewGame(t *testing.T) {
	game := NewGame(4)

	if len(game.Players) != 4 {
		t.Errorf("NewGame(4) created %d players, want 4", len(game.Players))
	}

	for _, player := range game.Players {
		if len(player.HoleCards) != 2 {
			t.Errorf("Player has %d hole cards, want 2", len(player.HoleCards))
		}
	}

	// 4 players * 2 cards = 8 cards dealt
	expectedRemaining := 52 - 8
	if game.Deck.Len() != expectedRemaining {
		t.Errorf("Deck has %d cards, want %d", game.Deck.Len(), expectedRemaining)
	}
}

func TestGameDealFlop(t *testing.T) {
	game := NewGame(2)

	err := game.DealFlop()
	if err != nil {
		t.Errorf("DealFlop() error = %v", err)
	}

	if game.Board.State() != Flop {
		t.Errorf("After DealFlop, state = %v, want Flop", game.Board.State())
	}
	if len(game.Board.Cards) != 3 {
		t.Errorf("Board has %d cards, want 3", len(game.Board.Cards))
	}
}

func TestGameDealTurn(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()

	err := game.DealTurn()
	if err != nil {
		t.Errorf("DealTurn() error = %v", err)
	}

	if game.Board.State() != Turn {
		t.Errorf("After DealTurn, state = %v, want Turn", game.Board.State())
	}
	if len(game.Board.Cards) != 4 {
		t.Errorf("Board has %d cards, want 4", len(game.Board.Cards))
	}
}

func TestGameDealRiver(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()
	game.DealTurn()

	err := game.DealRiver()
	if err != nil {
		t.Errorf("DealRiver() error = %v", err)
	}

	if game.Board.State() != River {
		t.Errorf("After DealRiver, state = %v, want River", game.Board.State())
	}
	if len(game.Board.Cards) != 5 {
		t.Errorf("Board has %d cards, want 5", len(game.Board.Cards))
	}
}

func TestGameDealNextStreet(t *testing.T) {
	game := NewGame(2)

	// Should deal flop
	err := game.DealNextStreet()
	if err != nil {
		t.Errorf("DealNextStreet() [flop] error = %v", err)
	}
	if game.Board.State() != Flop {
		t.Errorf("State = %v, want Flop", game.Board.State())
	}

	// Should deal turn
	err = game.DealNextStreet()
	if err != nil {
		t.Errorf("DealNextStreet() [turn] error = %v", err)
	}
	if game.Board.State() != Turn {
		t.Errorf("State = %v, want Turn", game.Board.State())
	}

	// Should deal river
	err = game.DealNextStreet()
	if err != nil {
		t.Errorf("DealNextStreet() [river] error = %v", err)
	}
	if game.Board.State() != River {
		t.Errorf("State = %v, want River", game.Board.State())
	}

	// Should error on river
	err = game.DealNextStreet()
	if err != ErrInvalidBoardState {
		t.Errorf("DealNextStreet() after river error = %v, want ErrInvalidBoardState", err)
	}
}

func TestGameDealInvalidState(t *testing.T) {
	game := NewGame(2)

	// Try to deal turn before flop
	err := game.DealTurn()
	if err != ErrInvalidBoardState {
		t.Errorf("DealTurn() before flop error = %v, want ErrInvalidBoardState", err)
	}

	// Try to deal river before flop
	err = game.DealRiver()
	if err != ErrInvalidBoardState {
		t.Errorf("DealRiver() before flop error = %v, want ErrInvalidBoardState", err)
	}

	// Try to deal flop twice
	game.DealFlop()
	err = game.DealFlop()
	if err != ErrInvalidBoardState {
		t.Errorf("DealFlop() twice error = %v, want ErrInvalidBoardState", err)
	}
}

func TestGameGetCandidateHands(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	hands, err := game.GetCandidateHands(game.Players[0])
	if err != nil {
		t.Errorf("GetCandidateHands() error = %v", err)
	}

	// 7 cards choose 5 = 21 combinations
	if len(hands) != 21 {
		t.Errorf("GetCandidateHands() returned %d hands, want 21", len(hands))
	}
}

func TestGameGetBestHand(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	hand, err := game.GetBestHand(game.Players[0])
	if err != nil {
		t.Errorf("GetBestHand() error = %v", err)
	}
	if hand == nil {
		t.Error("GetBestHand() returned nil")
	}
}

func TestGameGetWinners(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	winners, hands, err := game.GetWinners()
	if err != nil {
		t.Errorf("GetWinners() error = %v", err)
	}
	if len(winners) == 0 {
		t.Error("GetWinners() returned no winners")
	}
	if len(winners) != len(hands) {
		t.Errorf("Winners count %d != hands count %d", len(winners), len(hands))
	}
}

func TestGameGetWinnersBeforeRiver(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()

	_, _, err := game.GetWinners()
	if err != ErrInvalidBoardState {
		t.Errorf("GetWinners() before river error = %v, want ErrInvalidBoardState", err)
	}
}

func TestGameDeckCardsRemaining(t *testing.T) {
	game := NewGame(4)
	// 52 - (4 players * 2 cards) = 44

	game.DealFlop() // burn 1, deal 3 = 44 - 4 = 40
	game.DealTurn() // burn 1, deal 1 = 40 - 2 = 38
	game.DealRiver() // burn 1, deal 1 = 38 - 2 = 36

	expected := 52 - 8 - 4 - 2 - 2 // hole cards, flop, turn, river (with burns)
	if game.Deck.Len() != expected {
		t.Errorf("Deck has %d cards remaining, want %d", game.Deck.Len(), expected)
	}
}

func TestGameGetCandidateHandsBeforeFlop(t *testing.T) {
	game := NewGame(2)
	// Don't deal flop

	_, err := game.GetCandidateHands(game.Players[0])
	if err != ErrInvalidBoardState {
		t.Errorf("GetCandidateHands() before flop error = %v, want ErrInvalidBoardState", err)
	}
}

func TestGameGetBestHandBeforeFlop(t *testing.T) {
	game := NewGame(2)
	// Don't deal flop

	_, err := game.GetBestHand(game.Players[0])
	if err != ErrInvalidBoardState {
		t.Errorf("GetBestHand() before flop error = %v, want ErrInvalidBoardState", err)
	}
}

func TestGameMultipleWinners(t *testing.T) {
	// Create a game where both players have same best hand
	game := &Game{
		Deck:  NewDeck(),
		Board: NewBoard(),
		Players: []*Player{
			NewPlayer("Player 1"),
			NewPlayer("Player 2"),
		},
	}

	// Give both players the same hole cards (different suits)
	game.Players[0].SetHoleCards([]Card{
		NewCard(Two, Spades),
		NewCard(Three, Spades),
	})
	game.Players[1].SetHoleCards([]Card{
		NewCard(Two, Hearts),
		NewCard(Three, Hearts),
	})

	// Set up board with a straight that both players play
	game.Board.SetFlop([]Card{
		NewCard(Four, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Six, Diamonds),
	})
	game.Board.SetTurn(NewCard(King, Clubs))
	game.Board.SetRiver(NewCard(Queen, Clubs))

	winners, _, err := game.GetWinners()
	if err != nil {
		t.Errorf("GetWinners() error = %v", err)
	}
	if len(winners) != 2 {
		t.Errorf("Expected 2 winners (split pot), got %d", len(winners))
	}
}

func TestDeckDrawManyError(t *testing.T) {
	deck := NewDeck()
	// Draw most cards
	for i := 0; i < 50; i++ {
		deck.Draw()
	}

	// Try to draw more than remaining
	_, err := deck.DrawMany(5)
	if err != ErrEmptyDeck {
		t.Errorf("DrawMany() error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealFlopEmptyDeck(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}
	game.Players[0].SetHoleCards([]Card{NewCard(Ace, Spades), NewCard(King, Hearts)})

	// Empty the deck
	for game.Deck.Len() > 0 {
		game.Deck.Draw()
	}

	err := game.DealFlop()
	if err != ErrEmptyDeck {
		t.Errorf("DealFlop() with empty deck error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealTurnEmptyDeck(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}
	game.Players[0].SetHoleCards([]Card{NewCard(Ace, Spades), NewCard(King, Hearts)})
	game.Board.SetFlop([]Card{
		NewCard(Two, Clubs),
		NewCard(Three, Diamonds),
		NewCard(Four, Hearts),
	})

	// Empty the deck
	for game.Deck.Len() > 0 {
		game.Deck.Draw()
	}

	err := game.DealTurn()
	if err != ErrEmptyDeck {
		t.Errorf("DealTurn() with empty deck error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealRiverEmptyDeck(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}
	game.Players[0].SetHoleCards([]Card{NewCard(Ace, Spades), NewCard(King, Hearts)})
	game.Board.SetFlop([]Card{
		NewCard(Two, Clubs),
		NewCard(Three, Diamonds),
		NewCard(Four, Hearts),
	})
	game.Board.SetTurn(NewCard(Five, Spades))

	// Empty the deck
	for game.Deck.Len() > 0 {
		game.Deck.Draw()
	}

	err := game.DealRiver()
	if err != ErrEmptyDeck {
		t.Errorf("DealRiver() with empty deck error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealHoleCardsEmptyDeck(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1"), NewPlayer("P2")},
	}

	// Empty the deck
	for game.Deck.Len() > 0 {
		game.Deck.Draw()
	}

	err := game.DealHoleCards()
	if err != ErrEmptyDeck {
		t.Errorf("DealHoleCards() with empty deck error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealFlopBurnError(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}

	// Leave only 3 cards (not enough for burn + flop)
	for game.Deck.Len() > 0 {
		game.Deck.Draw()
	}

	err := game.DealFlop()
	if err != ErrEmptyDeck {
		t.Errorf("DealFlop() with insufficient cards error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealTurnDrawError(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}
	game.Board.SetFlop([]Card{
		NewCard(Two, Clubs),
		NewCard(Three, Diamonds),
		NewCard(Four, Hearts),
	})

	// Leave only 1 card (enough for burn, not for draw)
	for game.Deck.Len() > 1 {
		game.Deck.Draw()
	}

	err := game.DealTurn()
	if err != ErrEmptyDeck {
		t.Errorf("DealTurn() with 1 card error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealRiverDrawError(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}
	game.Board.SetFlop([]Card{
		NewCard(Two, Clubs),
		NewCard(Three, Diamonds),
		NewCard(Four, Hearts),
	})
	game.Board.SetTurn(NewCard(Five, Spades))

	// Leave only 1 card (enough for burn, not for draw)
	for game.Deck.Len() > 1 {
		game.Deck.Draw()
	}

	err := game.DealRiver()
	if err != ErrEmptyDeck {
		t.Errorf("DealRiver() with 1 card error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameDealFlopDrawManyError(t *testing.T) {
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}

	// Leave only 2 cards (enough for burn, not enough for DrawMany(3))
	for game.Deck.Len() > 2 {
		game.Deck.Draw()
	}

	err := game.DealFlop()
	if err != ErrEmptyDeck {
		t.Errorf("DealFlop() with 2 cards error = %v, want ErrEmptyDeck", err)
	}
}

func TestGameGetBestHandNoHands(t *testing.T) {
	// Create game with player having no hole cards
	game := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: []*Player{NewPlayer("P1")},
	}
	// Don't set hole cards - player.HoleCards will be empty
	game.Board.SetFlop([]Card{
		NewCard(Two, Clubs),
		NewCard(Three, Diamonds),
		NewCard(Four, Hearts),
	})

	// With only 3 board cards and 0 hole cards, can't make valid 5-card hands
	_, err := game.GetBestHand(game.Players[0])
	if err != ErrInvalidBoardState {
		t.Errorf("GetBestHand() with no valid hands error = %v, want ErrInvalidBoardState", err)
	}
}


func TestGameGetCandidateHandsSkipsInvalid(t *testing.T) {
	// Normal game - all combinations should be valid
	game := NewGame(2)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	hands, err := game.GetCandidateHands(game.Players[0])
	if err != nil {
		t.Errorf("GetCandidateHands() error = %v", err)
	}
	// 7 choose 5 = 21 combinations, all should be valid
	if len(hands) != 21 {
		t.Errorf("GetCandidateHands() = %d hands, want 21", len(hands))
	}
}
