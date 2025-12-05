package goker

import (
	"runtime"
	"sync"
)

// GetWinnersParallel returns the winning player(s) using concurrent evaluation.
// More efficient than GetWinners when there are many players.
func (g *Game) GetWinnersParallel() ([]*Player, []*Hand, error) {
	if g.Board.State() != River {
		return nil, nil, ErrInvalidBoardState
	}

	type playerHand struct {
		player *Player
		hand   *Hand
		err    error
	}

	results := make([]playerHand, len(g.Players))
	var wg sync.WaitGroup

	for i, player := range g.Players {
		wg.Add(1)
		go func(idx int, p *Player) {
			defer wg.Done()
			best, err := g.GetBestHand(p)
			results[idx] = playerHand{p, best, err}
		}(i, player)
	}

	wg.Wait()

	// Check for errors
	for _, r := range results {
		if r.err != nil {
			return nil, nil, r.err
		}
	}

	// Find winners
	var winners []playerHand
	for _, ph := range results {
		if len(winners) == 0 {
			winners = []playerHand{ph}
			continue
		}

		cmp := ph.hand.Compare(winners[0].hand)
		if cmp > 0 {
			winners = []playerHand{ph}
		} else if cmp == 0 {
			winners = append(winners, ph)
		}
	}

	players := make([]*Player, len(winners))
	hands := make([]*Hand, len(winners))
	for i, w := range winners {
		players[i] = w.player
		hands[i] = w.hand
	}

	return players, hands, nil
}

// EvaluateHandsBatch evaluates multiple hands concurrently using a worker pool.
// Returns hands in the same order as input. Invalid hands return nil.
func EvaluateHandsBatch(cardSets [][]Card, workers int) []*Hand {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	n := len(cardSets)
	results := make([]*Hand, n)

	// Use semaphore pattern for worker pool
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for i, cards := range cardSets {
		wg.Add(1)
		go func(idx int, c []Card) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			hand, err := NewHand(c)
			if err == nil {
				results[idx] = hand
			}
		}(i, cards)
	}

	wg.Wait()
	return results
}

// EvaluateAndCompareBatch evaluates hands and returns the best one(s).
// Returns indices of winning hands and the winning hand.
func EvaluateAndCompareBatch(cardSets [][]Card, workers int) ([]int, *Hand) {
	hands := EvaluateHandsBatch(cardSets, workers)

	var bestIndices []int
	var bestHand *Hand

	for i, hand := range hands {
		if hand == nil {
			continue
		}

		if bestHand == nil {
			bestIndices = []int{i}
			bestHand = hand
			continue
		}

		cmp := hand.Compare(bestHand)
		if cmp > 0 {
			bestIndices = []int{i}
			bestHand = hand
		} else if cmp == 0 {
			bestIndices = append(bestIndices, i)
		}
	}

	return bestIndices, bestHand
}
