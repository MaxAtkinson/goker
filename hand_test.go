package goker

import "testing"

// Helper to create hands easily
func makeHand(t *testing.T, cards ...Card) *Hand {
	t.Helper()
	h, err := NewHand(cards)
	if err != nil {
		t.Fatalf("Failed to create hand: %v", err)
	}
	return h
}

func TestNewHandInvalidSize(t *testing.T) {
	cards := []Card{
		NewCard(Ace, Spades),
		NewCard(King, Spades),
	}
	_, err := NewHand(cards)
	if err != ErrInvalidHandSize {
		t.Errorf("NewHand() error = %v, want ErrInvalidHandSize", err)
	}
}

func TestNewHandDuplicates(t *testing.T) {
	cards := []Card{
		NewCard(Ace, Spades),
		NewCard(Ace, Spades),
		NewCard(King, Spades),
		NewCard(Queen, Spades),
		NewCard(Jack, Spades),
	}
	_, err := NewHand(cards)
	if err != ErrDuplicateCards {
		t.Errorf("NewHand() error = %v, want ErrDuplicateCards", err)
	}
}

func TestHighCard(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Nine, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Spades),
	)

	if hand.Rank() != HighCard {
		t.Errorf("Hand rank = %v, want HighCard", hand.Rank())
	}
}

func TestPair(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Ace, Hearts),
		NewCard(King, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Spades),
	)

	if hand.Rank() != Pair {
		t.Errorf("Hand rank = %v, want Pair", hand.Rank())
	}
}

func TestTwoPair(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Ace, Hearts),
		NewCard(King, Diamonds),
		NewCard(King, Clubs),
		NewCard(Two, Spades),
	)

	if hand.Rank() != TwoPair {
		t.Errorf("Hand rank = %v, want TwoPair", hand.Rank())
	}
}

func TestThreeOfAKind(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Ace, Hearts),
		NewCard(Ace, Diamonds),
		NewCard(King, Clubs),
		NewCard(Two, Spades),
	)

	if hand.Rank() != ThreeOfAKind {
		t.Errorf("Hand rank = %v, want ThreeOfAKind", hand.Rank())
	}
}

func TestStraight(t *testing.T) {
	hand := makeHand(t,
		NewCard(Nine, Spades),
		NewCard(Eight, Hearts),
		NewCard(Seven, Diamonds),
		NewCard(Six, Clubs),
		NewCard(Five, Spades),
	)

	if hand.Rank() != Straight {
		t.Errorf("Hand rank = %v, want Straight", hand.Rank())
	}
	if !hand.IsStraight() {
		t.Error("IsStraight() = false, want true")
	}
}

func TestWheel(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Two, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Spades),
	)

	if hand.Rank() != Straight {
		t.Errorf("Hand rank = %v, want Straight", hand.Rank())
	}
	if !hand.IsWheel() {
		t.Error("IsWheel() = false, want true")
	}
}

func TestFlush(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Spades),
		NewCard(Nine, Spades),
		NewCard(Five, Spades),
		NewCard(Two, Spades),
	)

	if hand.Rank() != Flush {
		t.Errorf("Hand rank = %v, want Flush", hand.Rank())
	}
	if !hand.IsFlush() {
		t.Error("IsFlush() = false, want true")
	}
}

func TestFullHouse(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Ace, Hearts),
		NewCard(Ace, Diamonds),
		NewCard(King, Clubs),
		NewCard(King, Spades),
	)

	if hand.Rank() != FullHouse {
		t.Errorf("Hand rank = %v, want FullHouse", hand.Rank())
	}
}

func TestFourOfAKind(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Ace, Hearts),
		NewCard(Ace, Diamonds),
		NewCard(Ace, Clubs),
		NewCard(King, Spades),
	)

	if hand.Rank() != FourOfAKind {
		t.Errorf("Hand rank = %v, want FourOfAKind", hand.Rank())
	}
}

func TestStraightFlush(t *testing.T) {
	hand := makeHand(t,
		NewCard(Nine, Spades),
		NewCard(Eight, Spades),
		NewCard(Seven, Spades),
		NewCard(Six, Spades),
		NewCard(Five, Spades),
	)

	if hand.Rank() != StraightFlush {
		t.Errorf("Hand rank = %v, want StraightFlush", hand.Rank())
	}
}

func TestRoyalFlush(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Spades),
		NewCard(Queen, Spades),
		NewCard(Jack, Spades),
		NewCard(Ten, Spades),
	)

	if hand.Rank() != RoyalFlush {
		t.Errorf("Hand rank = %v, want RoyalFlush", hand.Rank())
	}
}

func TestHandComparison(t *testing.T) {
	highCard := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Nine, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Spades),
	)

	pair := makeHand(t,
		NewCard(Two, Spades),
		NewCard(Two, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Hearts),
	)

	if !pair.Beats(highCard) {
		t.Error("Pair should beat high card")
	}
	if highCard.Beats(pair) {
		t.Error("High card should not beat pair")
	}
}

func TestHandTiebreaker(t *testing.T) {
	// Higher high card wins
	hand1 := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Nine, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Spades),
	)

	hand2 := makeHand(t,
		NewCard(King, Spades),
		NewCard(Queen, Hearts),
		NewCard(Nine, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Hearts),
	)

	if !hand1.Beats(hand2) {
		t.Error("Higher high card should win")
	}
}

func TestHandTiebreakerPair(t *testing.T) {
	// Higher pair wins
	hand1 := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Ace, Hearts),
		NewCard(Two, Diamonds),
		NewCard(Three, Clubs),
		NewCard(Four, Spades),
	)

	hand2 := makeHand(t,
		NewCard(King, Spades),
		NewCard(King, Hearts),
		NewCard(Ace, Diamonds),
		NewCard(Queen, Clubs),
		NewCard(Jack, Hearts),
	)

	if !hand1.Beats(hand2) {
		t.Error("Higher pair should win")
	}
}

func TestWheelVsHigherStraight(t *testing.T) {
	wheel := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Two, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Spades),
	)

	sixHigh := makeHand(t,
		NewCard(Six, Spades),
		NewCard(Five, Hearts),
		NewCard(Four, Diamonds),
		NewCard(Three, Clubs),
		NewCard(Two, Hearts),
	)

	if wheel.Beats(sixHigh) {
		t.Error("Wheel should not beat 6-high straight")
	}
	if !sixHigh.Beats(wheel) {
		t.Error("6-high straight should beat wheel")
	}
}

func TestHandTies(t *testing.T) {
	hand1 := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Nine, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Spades),
	)

	hand2 := makeHand(t,
		NewCard(Ace, Hearts),
		NewCard(King, Diamonds),
		NewCard(Nine, Clubs),
		NewCard(Five, Spades),
		NewCard(Two, Hearts),
	)

	if !hand1.Ties(hand2) {
		t.Error("Identical ranks should tie")
	}
}

func TestHandContains(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Nine, Diamonds),
		NewCard(Five, Clubs),
		NewCard(Two, Spades),
	)

	if !hand.Contains(Ace) {
		t.Error("Hand should contain Ace")
	}
	if hand.Contains(Queen) {
		t.Error("Hand should not contain Queen")
	}
}

func TestHandString(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Spades),
		NewCard(Queen, Spades),
		NewCard(Jack, Spades),
		NewCard(Ten, Spades),
	)

	str := hand.String()
	if str == "" {
		t.Error("Hand.String() should not be empty")
	}
}

func TestHandRankHierarchy(t *testing.T) {
	// Create one of each hand type, verify ranking order
	hands := []*Hand{
		makeHand(t, NewCard(Ace, Spades), NewCard(King, Hearts), NewCard(Nine, Diamonds), NewCard(Five, Clubs), NewCard(Two, Spades)),         // High Card
		makeHand(t, NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(Nine, Diamonds), NewCard(Five, Clubs), NewCard(Two, Spades)),          // Pair
		makeHand(t, NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(King, Diamonds), NewCard(King, Clubs), NewCard(Two, Spades)),          // Two Pair
		makeHand(t, NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(Ace, Diamonds), NewCard(Five, Clubs), NewCard(Two, Spades)),           // Three of a Kind
		makeHand(t, NewCard(Nine, Spades), NewCard(Eight, Hearts), NewCard(Seven, Diamonds), NewCard(Six, Clubs), NewCard(Five, Spades)),     // Straight
		makeHand(t, NewCard(Ace, Spades), NewCard(King, Spades), NewCard(Nine, Spades), NewCard(Five, Spades), NewCard(Two, Spades)),          // Flush
		makeHand(t, NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(Ace, Diamonds), NewCard(King, Clubs), NewCard(King, Spades)),          // Full House
		makeHand(t, NewCard(Ace, Spades), NewCard(Ace, Hearts), NewCard(Ace, Diamonds), NewCard(Ace, Clubs), NewCard(Two, Spades)),            // Four of a Kind
		makeHand(t, NewCard(Nine, Spades), NewCard(Eight, Spades), NewCard(Seven, Spades), NewCard(Six, Spades), NewCard(Five, Spades)),       // Straight Flush
		makeHand(t, NewCard(Ace, Spades), NewCard(King, Spades), NewCard(Queen, Spades), NewCard(Jack, Spades), NewCard(Ten, Spades)),         // Royal Flush
	}

	for i := 0; i < len(hands)-1; i++ {
		if !hands[i+1].Beats(hands[i]) {
			t.Errorf("Hand %d (%v) should beat hand %d (%v)", i+1, hands[i+1].Rank(), i, hands[i].Rank())
		}
	}
}

func TestNewHandWithPlayer(t *testing.T) {
	player := NewPlayer("Test")
	cards := []Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
		NewCard(Jack, Clubs),
		NewCard(Ten, Spades),
	}

	hand, err := NewHandWithPlayer(cards, player)
	if err != nil {
		t.Errorf("NewHandWithPlayer() error = %v", err)
	}
	if hand.Player != player {
		t.Error("Hand.Player should be set")
	}
}

func TestNewHandWithPlayerInvalid(t *testing.T) {
	player := NewPlayer("Test")
	cards := []Card{NewCard(Ace, Spades)} // Invalid - not 5 cards

	_, err := NewHandWithPlayer(cards, player)
	if err != ErrInvalidHandSize {
		t.Errorf("NewHandWithPlayer() error = %v, want ErrInvalidHandSize", err)
	}
}

func TestHandCompareEqual(t *testing.T) {
	hand1 := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
		NewCard(Jack, Clubs),
		NewCard(Nine, Spades),
	)
	hand2 := makeHand(t,
		NewCard(Ace, Hearts),
		NewCard(King, Diamonds),
		NewCard(Queen, Clubs),
		NewCard(Jack, Spades),
		NewCard(Nine, Hearts),
	)

	if hand1.Compare(hand2) != 0 {
		t.Error("Equal hands should return 0")
	}
}

func TestHandCompareLoss(t *testing.T) {
	pair := makeHand(t,
		NewCard(Two, Spades),
		NewCard(Two, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Spades),
	)
	trips := makeHand(t,
		NewCard(Two, Spades),
		NewCard(Two, Hearts),
		NewCard(Two, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Hearts),
	)

	if pair.Compare(trips) != -1 {
		t.Error("Pair should lose to trips (return -1)")
	}
}

func TestWheelTiesWheel(t *testing.T) {
	wheel1 := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(Two, Hearts),
		NewCard(Three, Diamonds),
		NewCard(Four, Clubs),
		NewCard(Five, Spades),
	)
	wheel2 := makeHand(t,
		NewCard(Ace, Hearts),
		NewCard(Two, Diamonds),
		NewCard(Three, Clubs),
		NewCard(Four, Spades),
		NewCard(Five, Hearts),
	)

	if !wheel1.Ties(wheel2) {
		t.Error("Two wheels should tie")
	}
}

func TestHandEvaluateIdempotent(t *testing.T) {
	hand := makeHand(t,
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
		NewCard(Jack, Clubs),
		NewCard(Ten, Spades),
	)

	// Call evaluate multiple times - should be idempotent
	rank1 := hand.Rank()
	// Manually trigger evaluate again (it's already called in NewHand)
	// The evaluate function has early return if already evaluated
	rank2 := hand.Rank()

	if rank1 != rank2 {
		t.Error("Hand rank should be consistent after multiple evaluations")
	}
}

func TestHandEmptyCardsFlush(t *testing.T) {
	// Create a hand struct directly with empty cards to test edge case
	h := &Hand{Cards: []Card{}}
	if h.computeIsFlush() {
		t.Error("Empty hand should not be a flush")
	}
}

func TestHandZeroCardValues(t *testing.T) {
	// Create a hand struct directly with empty cards to test edge case
	h := &Hand{Cards: []Card{}}
	isStraight, isWheel := h.computeIsStraight()
	if isStraight || isWheel {
		t.Error("Empty hand should not be a straight")
	}
}

func TestHandEvaluateAlreadyEvaluated(t *testing.T) {
	// Create a hand and manually set evaluated flag
	h := &Hand{
		Cards: []Card{
			NewCard(Ace, Spades),
			NewCard(King, Hearts),
			NewCard(Queen, Diamonds),
			NewCard(Jack, Clubs),
			NewCard(Ten, Spades),
		},
		evaluated: true,
		handRank:  HighCard, // Wrong rank, but evaluate should skip
	}

	// Call evaluate - should return early
	h.evaluate()

	// Should still have wrong rank since evaluate returned early
	if h.handRank != HighCard {
		t.Error("evaluate() should return early when already evaluated")
	}
}
