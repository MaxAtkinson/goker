package goker

import "testing"

func TestHandRankString(t *testing.T) {
	tests := []struct {
		rank     HandRank
		expected string
	}{
		{HighCard, "High Card"},
		{Pair, "Pair"},
		{TwoPair, "Two Pair"},
		{ThreeOfAKind, "Three of a Kind"},
		{Straight, "Straight"},
		{Flush, "Flush"},
		{FullHouse, "Full House"},
		{FourOfAKind, "Four of a Kind"},
		{StraightFlush, "Straight Flush"},
		{RoyalFlush, "Royal Flush"},
		{HandRank(99), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.rank.String(); got != tt.expected {
			t.Errorf("HandRank(%d).String() = %s, want %s", tt.rank, got, tt.expected)
		}
	}
}
