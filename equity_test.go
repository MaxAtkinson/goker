package goker

import (
	"math"
	"testing"
)

func TestNewEquityCalculator(t *testing.T) {
	ec := NewEquityCalculator(0)
	if ec.workers <= 0 {
		t.Error("NewEquityCalculator(0) should use default workers")
	}

	ec = NewEquityCalculator(4)
	if ec.workers != 4 {
		t.Errorf("NewEquityCalculator(4) workers = %d, want 4", ec.workers)
	}
}

func TestEquityCalculatorBasic(t *testing.T) {
	ec := NewEquityCalculator(4)

	// AA vs KK preflop - AA should have ~80% equity
	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
	}

	results := ec.Calculate(holeCards, []Card{}, 1000)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// AA should win more often than KK
	if results[0].Equity < results[1].Equity {
		t.Errorf("AA equity (%.2f) should be higher than KK equity (%.2f)",
			results[0].Equity, results[1].Equity)
	}

	// AA typically has 80%+ equity vs KK
	if results[0].Equity < 0.70 || results[0].Equity > 0.90 {
		t.Logf("AA equity = %.2f (expected ~0.80)", results[0].Equity)
	}

	// Total should equal simulations
	if results[0].Total != 1000 {
		t.Errorf("Total = %d, want 1000", results[0].Total)
	}
}

func TestEquityCalculatorWithBoard(t *testing.T) {
	ec := NewEquityCalculator(4)

	// Player 1 has set, Player 2 has overpair
	holeCards := [][]Card{
		{NewCard(Seven, Spades), NewCard(Seven, Hearts)}, // Set of 7s
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},     // Overpair AA
	}

	board := []Card{
		NewCard(Seven, Diamonds),
		NewCard(Two, Clubs),
		NewCard(Five, Hearts),
	}

	results := ec.Calculate(holeCards, board, 1000)

	// Set should dominate overpair
	if results[0].Equity < results[1].Equity {
		t.Errorf("Set equity (%.2f) should be higher than overpair equity (%.2f)",
			results[0].Equity, results[1].Equity)
	}
}

func TestEquityCalculatorMultiplePlayers(t *testing.T) {
	ec := NewEquityCalculator(4)

	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
		{NewCard(Queen, Spades), NewCard(Queen, Hearts)},
	}

	results := ec.Calculate(holeCards, []Card{}, 2000)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Equities should roughly sum to 1 (allowing for tie handling)
	totalEquity := results[0].Equity + results[1].Equity + results[2].Equity
	if math.Abs(totalEquity-1.0) > 0.1 {
		t.Errorf("Total equity = %.2f, expected ~1.0", totalEquity)
	}

	// AA should have significantly more equity than KK and QQ
	// KK vs QQ is very close, so we only check AA > others
	if results[0].Equity < results[1].Equity {
		t.Errorf("AA equity (%.2f) should be higher than KK equity (%.2f)",
			results[0].Equity, results[1].Equity)
	}
	if results[0].Equity < results[2].Equity {
		t.Errorf("AA equity (%.2f) should be higher than QQ equity (%.2f)",
			results[0].Equity, results[2].Equity)
	}
}

func TestEquityCalculatorExact(t *testing.T) {
	ec := NewEquityCalculator(4)

	// Test with turn dealt - only 1 card to come (46 possibilities)
	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
	}

	board := []Card{
		NewCard(Two, Clubs),
		NewCard(Five, Diamonds),
		NewCard(Nine, Hearts),
		NewCard(Jack, Spades),
	}

	results := ec.CalculateExact(holeCards, board, 1000)

	if results == nil {
		t.Fatal("CalculateExact returned nil")
	}

	// Should have exact 44 total (52 - 4 hole - 4 board)
	if results[0].Total != 44 {
		t.Errorf("Total = %d, want 44", results[0].Total)
	}

	// AA should still be ahead
	if results[0].Equity < results[1].Equity {
		t.Errorf("AA equity (%.2f) should be higher than KK equity (%.2f)",
			results[0].Equity, results[1].Equity)
	}
}

func TestEquityCalculatorExactTooMany(t *testing.T) {
	ec := NewEquityCalculator(4)

	// Preflop - way too many combinations
	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
	}

	results := ec.CalculateExact(holeCards, []Card{}, 1000)

	if results != nil {
		t.Error("CalculateExact should return nil when too many combinations")
	}
}

func TestEquityResultFields(t *testing.T) {
	ec := NewEquityCalculator(2)

	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(Two, Spades), NewCard(Seven, Hearts)},
	}

	results := ec.Calculate(holeCards, []Card{}, 100)

	for i, r := range results {
		// Wins + Ties + Losses should equal Total
		if r.Wins+r.Ties+r.Losses != r.Total {
			t.Errorf("Player %d: Wins(%d) + Ties(%d) + Losses(%d) != Total(%d)",
				i, r.Wins, r.Ties, r.Losses, r.Total)
		}

		// Equity should be between 0 and 1
		if r.Equity < 0 || r.Equity > 1 {
			t.Errorf("Player %d: Equity = %.2f, should be between 0 and 1", i, r.Equity)
		}
	}
}

func TestCardKey(t *testing.T) {
	c1 := NewCard(Ace, Spades)
	c2 := NewCard(Ace, Hearts)
	c3 := NewCard(Ace, Spades)

	if cardKey(c1) == cardKey(c2) {
		t.Error("Different cards should have different keys")
	}
	if cardKey(c1) != cardKey(c3) {
		t.Error("Same cards should have same keys")
	}
}

func TestBuildRemainingDeck(t *testing.T) {
	used := map[string]bool{
		cardKey(NewCard(Ace, Spades)):   true,
		cardKey(NewCard(King, Hearts)):  true,
		cardKey(NewCard(Queen, Diamonds)): true,
	}

	deck := buildRemainingDeck(used)

	if len(deck) != 49 {
		t.Errorf("Remaining deck has %d cards, want 49", len(deck))
	}

	// Check used cards are not in deck
	for _, c := range deck {
		if used[cardKey(c)] {
			t.Errorf("Used card %v found in remaining deck", c)
		}
	}
}

func TestFindBestHand(t *testing.T) {
	cards := []Card{
		NewCard(Ace, Spades),
		NewCard(King, Spades),
		NewCard(Queen, Spades),
		NewCard(Jack, Spades),
		NewCard(Ten, Spades),
		NewCard(Two, Hearts),
		NewCard(Three, Hearts),
	}

	best := findBestHand(cards)

	if best == nil {
		t.Fatal("findBestHand returned nil")
	}
	if best.Rank() != RoyalFlush {
		t.Errorf("Best hand = %v, want RoyalFlush", best.Rank())
	}
}

func TestFindWinnerIndices(t *testing.T) {
	hands := []*Hand{
		makeTestHand(Pair),
		makeTestHand(RoyalFlush),
		makeTestHand(HighCard),
	}

	winners := findWinnerIndices(hands)

	if len(winners) != 1 || winners[0] != 1 {
		t.Errorf("Winners = %v, want [1]", winners)
	}
}

func TestFindWinnerIndicesTie(t *testing.T) {
	// Create two identical high card hands
	hand1, _ := NewHand([]Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
		NewCard(Jack, Clubs),
		NewCard(Nine, Spades),
	})
	hand2, _ := NewHand([]Card{
		NewCard(Ace, Hearts),
		NewCard(King, Diamonds),
		NewCard(Queen, Clubs),
		NewCard(Jack, Spades),
		NewCard(Nine, Hearts),
	})

	winners := findWinnerIndices([]*Hand{hand1, hand2})

	if len(winners) != 2 {
		t.Errorf("Should have 2 winners (tie), got %d", len(winners))
	}
}

func makeTestHand(rank HandRank) *Hand {
	var cards []Card
	switch rank {
	case RoyalFlush:
		cards = []Card{
			NewCard(Ace, Spades),
			NewCard(King, Spades),
			NewCard(Queen, Spades),
			NewCard(Jack, Spades),
			NewCard(Ten, Spades),
		}
	case Pair:
		cards = []Card{
			NewCard(Ace, Spades),
			NewCard(Ace, Hearts),
			NewCard(King, Diamonds),
			NewCard(Queen, Clubs),
			NewCard(Jack, Spades),
		}
	default: // HighCard
		cards = []Card{
			NewCard(Ace, Spades),
			NewCard(King, Hearts),
			NewCard(Queen, Diamonds),
			NewCard(Jack, Clubs),
			NewCard(Nine, Spades),
		}
	}
	hand, _ := NewHand(cards)
	return hand
}

func BenchmarkEquityCalculate1000(b *testing.B) {
	ec := NewEquityCalculator(0)

	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ec.Calculate(holeCards, []Card{}, 1000)
	}
}

func BenchmarkEquityCalculate10000(b *testing.B) {
	ec := NewEquityCalculator(0)

	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ec.Calculate(holeCards, []Card{}, 10000)
	}
}

func BenchmarkEquityExactRiver(b *testing.B) {
	ec := NewEquityCalculator(0)

	holeCards := [][]Card{
		{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
		{NewCard(King, Spades), NewCard(King, Hearts)},
	}

	board := []Card{
		NewCard(Two, Clubs),
		NewCard(Five, Diamonds),
		NewCard(Nine, Hearts),
		NewCard(Jack, Spades),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ec.CalculateExact(holeCards, board, 1000)
	}
}
