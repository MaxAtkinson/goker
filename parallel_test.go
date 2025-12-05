package goker

import (
	"testing"
)

func TestGetWinnersParallel(t *testing.T) {
	game := NewGame(4)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	winners, hands, err := game.GetWinnersParallel()
	if err != nil {
		t.Errorf("GetWinnersParallel() error = %v", err)
	}
	if len(winners) == 0 {
		t.Error("GetWinnersParallel() returned no winners")
	}
	if len(winners) != len(hands) {
		t.Errorf("Winners count %d != hands count %d", len(winners), len(hands))
	}
}

func TestGetWinnersParallelBeforeRiver(t *testing.T) {
	game := NewGame(2)
	game.DealFlop()

	_, _, err := game.GetWinnersParallel()
	if err != ErrInvalidBoardState {
		t.Errorf("GetWinnersParallel() before river error = %v, want ErrInvalidBoardState", err)
	}
}

func TestGetWinnersParallelMatchesSequential(t *testing.T) {
	// Run multiple times to catch race conditions
	for i := 0; i < 10; i++ {
		game := NewGame(6)
		game.DealFlop()
		game.DealTurn()
		game.DealRiver()

		seqWinners, seqHands, _ := game.GetWinners()
		parWinners, parHands, _ := game.GetWinnersParallel()

		if len(seqWinners) != len(parWinners) {
			t.Errorf("Winner count mismatch: sequential=%d, parallel=%d", len(seqWinners), len(parWinners))
		}

		// Check same players won
		seqNames := make(map[string]bool)
		for _, p := range seqWinners {
			seqNames[p.Name] = true
		}
		for _, p := range parWinners {
			if !seqNames[p.Name] {
				t.Errorf("Parallel winner %s not in sequential winners", p.Name)
			}
		}

		// Check same hand ranks
		if len(seqHands) > 0 && len(parHands) > 0 {
			if seqHands[0].Rank() != parHands[0].Rank() {
				t.Errorf("Hand rank mismatch: sequential=%v, parallel=%v", seqHands[0].Rank(), parHands[0].Rank())
			}
		}
	}
}

func TestEvaluateHandsBatch(t *testing.T) {
	cardSets := [][]Card{
		// Royal flush
		{NewCard(Ace, Spades), NewCard(King, Spades), NewCard(Queen, Spades), NewCard(Jack, Spades), NewCard(Ten, Spades)},
		// Pair
		{NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(King, Diamonds), NewCard(Queen, Clubs), NewCard(Jack, Spades)},
		// High card
		{NewCard(Ace, Spades), NewCard(King, Hearts), NewCard(Queen, Diamonds), NewCard(Jack, Clubs), NewCard(Nine, Spades)},
	}

	hands := EvaluateHandsBatch(cardSets, 2)

	if len(hands) != 3 {
		t.Errorf("EvaluateHandsBatch() returned %d hands, want 3", len(hands))
	}

	if hands[0].Rank() != RoyalFlush {
		t.Errorf("Hand 0 rank = %v, want RoyalFlush", hands[0].Rank())
	}
	if hands[1].Rank() != Pair {
		t.Errorf("Hand 1 rank = %v, want Pair", hands[1].Rank())
	}
	if hands[2].Rank() != HighCard {
		t.Errorf("Hand 2 rank = %v, want HighCard", hands[2].Rank())
	}
}

func TestEvaluateHandsBatchWithInvalid(t *testing.T) {
	cardSets := [][]Card{
		// Valid
		{NewCard(Ace, Spades), NewCard(King, Spades), NewCard(Queen, Spades), NewCard(Jack, Spades), NewCard(Ten, Spades)},
		// Invalid - only 3 cards
		{NewCard(Ace, Spades), NewCard(King, Hearts), NewCard(Queen, Diamonds)},
		// Valid
		{NewCard(Two, Spades), NewCard(Two, Hearts), NewCard(Three, Diamonds), NewCard(Four, Clubs), NewCard(Five, Spades)},
	}

	hands := EvaluateHandsBatch(cardSets, 0) // Use default workers

	if hands[0] == nil {
		t.Error("Hand 0 should be valid")
	}
	if hands[1] != nil {
		t.Error("Hand 1 should be nil (invalid)")
	}
	if hands[2] == nil {
		t.Error("Hand 2 should be valid")
	}
}

func TestEvaluateAndCompareBatch(t *testing.T) {
	cardSets := [][]Card{
		// Pair
		{NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(King, Diamonds), NewCard(Queen, Clubs), NewCard(Jack, Spades)},
		// Royal flush - winner
		{NewCard(Ace, Spades), NewCard(King, Spades), NewCard(Queen, Spades), NewCard(Jack, Spades), NewCard(Ten, Spades)},
		// High card
		{NewCard(Ace, Spades), NewCard(King, Hearts), NewCard(Queen, Diamonds), NewCard(Jack, Clubs), NewCard(Nine, Spades)},
	}

	indices, bestHand := EvaluateAndCompareBatch(cardSets, 2)

	if len(indices) != 1 || indices[0] != 1 {
		t.Errorf("Winner indices = %v, want [1]", indices)
	}
	if bestHand.Rank() != RoyalFlush {
		t.Errorf("Best hand rank = %v, want RoyalFlush", bestHand.Rank())
	}
}

func TestEvaluateAndCompareBatchTie(t *testing.T) {
	cardSets := [][]Card{
		// Same high card hand
		{NewCard(Ace, Spades), NewCard(King, Hearts), NewCard(Queen, Diamonds), NewCard(Jack, Clubs), NewCard(Nine, Spades)},
		{NewCard(Ace, Hearts), NewCard(King, Diamonds), NewCard(Queen, Clubs), NewCard(Jack, Spades), NewCard(Nine, Hearts)},
	}

	indices, _ := EvaluateAndCompareBatch(cardSets, 2)

	if len(indices) != 2 {
		t.Errorf("Should have 2 winners (tie), got %d", len(indices))
	}
}

func BenchmarkGetWinners(b *testing.B) {
	game := NewGame(9)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.GetWinners()
	}
}

func BenchmarkGetWinnersParallel(b *testing.B) {
	game := NewGame(9)
	game.DealFlop()
	game.DealTurn()
	game.DealRiver()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.GetWinnersParallel()
	}
}

func BenchmarkEvaluateHandsBatch100(b *testing.B) {
	// Generate 100 random hands
	cardSets := make([][]Card, 100)
	for i := 0; i < 100; i++ {
		deck := NewDeck()
		cards, _ := deck.DrawMany(5)
		cardSets[i] = cards
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateHandsBatch(cardSets, 0)
	}
}

func BenchmarkEvaluateHandsBatch1000(b *testing.B) {
	// Generate 1000 random hands
	cardSets := make([][]Card, 1000)
	for i := 0; i < 1000; i++ {
		deck := NewDeck()
		cards, _ := deck.DrawMany(5)
		cardSets[i] = cards
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateHandsBatch(cardSets, 0)
	}
}
