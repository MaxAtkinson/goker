package goker

import "testing"

func TestCombinations(t *testing.T) {
	items := []int{1, 2, 3, 4}
	combos := Combinations(items, 2)

	// 4 choose 2 = 6
	if len(combos) != 6 {
		t.Errorf("Combinations(4, 2) = %d combinations, want 6", len(combos))
	}

	expected := [][]int{
		{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4},
	}

	for i, combo := range combos {
		if combo[0] != expected[i][0] || combo[1] != expected[i][1] {
			t.Errorf("Combination %d = %v, want %v", i, combo, expected[i])
		}
	}
}

func TestCombinationsEdgeCases(t *testing.T) {
	items := []int{1, 2, 3}

	// n = 0
	combos := Combinations(items, 0)
	if combos != nil {
		t.Errorf("Combinations(3, 0) should be nil")
	}

	// n > len
	combos = Combinations(items, 5)
	if combos != nil {
		t.Errorf("Combinations(3, 5) should be nil")
	}
}

func TestCardCombinations(t *testing.T) {
	cards := []Card{
		NewCard(Ace, Spades),
		NewCard(King, Hearts),
		NewCard(Queen, Diamonds),
		NewCard(Jack, Clubs),
		NewCard(Ten, Spades),
		NewCard(Nine, Hearts),
		NewCard(Eight, Diamonds),
	}

	// 7 choose 5 = 21
	combos := CardCombinations(cards, 5)
	if len(combos) != 21 {
		t.Errorf("CardCombinations(7, 5) = %d combinations, want 21", len(combos))
	}

	// Each combination should have 5 cards
	for i, combo := range combos {
		if len(combo) != 5 {
			t.Errorf("Combination %d has %d cards, want 5", i, len(combo))
		}
	}
}

func TestBitSequenceToInt(t *testing.T) {
	tests := []struct {
		bits     []int
		expected int
	}{
		{[]int{1, 0, 1}, 5},      // 101 = 5
		{[]int{1, 1, 1, 1}, 15},  // 1111 = 15
		{[]int{1, 0, 0, 0}, 8},   // 1000 = 8
		{[]int{0, 0, 0, 1}, 1},   // 0001 = 1
	}

	for _, tt := range tests {
		if got := BitSequenceToInt(tt.bits); got != tt.expected {
			t.Errorf("BitSequenceToInt(%v) = %d, want %d", tt.bits, got, tt.expected)
		}
	}
}

func TestGetBinaryIndexFromCardRank(t *testing.T) {
	// Ace should be at index 0 (highest)
	aceIdx := GetBinaryIndexFromCardRank(Ace)
	// Two should be at index 12 (lowest)
	twoIdx := GetBinaryIndexFromCardRank(Two)

	if aceIdx >= twoIdx {
		t.Errorf("Ace index (%d) should be less than Two index (%d)", aceIdx, twoIdx)
	}
}
