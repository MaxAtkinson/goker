package goker

// Combinations generates all n-element combinations from the given items.
func Combinations[T any](items []T, n int) [][]T {
	if n <= 0 || n > len(items) {
		return nil
	}

	result := make([][]T, 0)
	combination := make([]T, n)
	var generate func(start, depth int)

	generate = func(start, depth int) {
		if depth == n {
			combo := make([]T, n)
			copy(combo, combination)
			result = append(result, combo)
			return
		}

		for i := start; i <= len(items)-(n-depth); i++ {
			combination[depth] = items[i]
			generate(i+1, depth+1)
		}
	}

	generate(0, 0)
	return result
}

// CardCombinations generates all n-card combinations from the given cards.
func CardCombinations(cards []Card, n int) [][]Card {
	return Combinations(cards, n)
}

// BitSequenceToInt converts a slice of bit values (0 or 1) to an integer.
func BitSequenceToInt(bits []int) int {
	result := 0
	for _, bit := range bits {
		result = (result << 1) | (bit & 1)
	}
	return result
}

// GetBinaryIndexFromCardRank returns the binary index position for a card rank.
func GetBinaryIndexFromCardRank(rank CardRank) int {
	return numBits - int(rank) + 1 - unusedBits
}
