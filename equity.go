package goker

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
)

// EquityResult holds the equity calculation results for a player.
type EquityResult struct {
	Wins   int     // Number of wins
	Ties   int     // Number of ties
	Losses int     // Number of losses
	Total  int     // Total simulations
	Equity float64 // Win probability (wins + ties/numPlayers) / total
}

// EquityCalculator calculates hand equity through Monte Carlo simulation.
type EquityCalculator struct {
	workers int
}

// NewEquityCalculator creates a new equity calculator with specified worker count.
// If workers <= 0, uses number of CPUs.
func NewEquityCalculator(workers int) *EquityCalculator {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &EquityCalculator{workers: workers}
}

// Calculate runs Monte Carlo simulation to determine equity for each player's hole cards.
// holeCards: slice of 2-card arrays for each player
// board: current community cards (0-5 cards)
// simulations: number of random board runouts to simulate
func (ec *EquityCalculator) Calculate(holeCards [][]Card, board []Card, simulations int) []EquityResult {
	numPlayers := len(holeCards)
	results := make([]EquityResult, numPlayers)

	// Track wins/ties atomically
	wins := make([]int64, numPlayers)
	ties := make([]int64, numPlayers)

	// Build set of used cards
	usedCards := make(map[string]bool)
	for _, hole := range holeCards {
		for _, c := range hole {
			usedCards[cardKey(c)] = true
		}
	}
	for _, c := range board {
		usedCards[cardKey(c)] = true
	}

	// Build remaining deck
	remainingDeck := buildRemainingDeck(usedCards)
	cardsNeeded := 5 - len(board)

	// Distribute simulations across workers
	simsPerWorker := simulations / ec.workers
	remainder := simulations % ec.workers

	var wg sync.WaitGroup
	for w := 0; w < ec.workers; w++ {
		workerSims := simsPerWorker
		if w < remainder {
			workerSims++
		}

		wg.Add(1)
		go func(numSims int) {
			defer wg.Done()
			ec.runSimulations(holeCards, board, remainingDeck, cardsNeeded, numSims, wins, ties)
		}(workerSims)
	}

	wg.Wait()

	// Calculate final results
	for i := 0; i < numPlayers; i++ {
		w := int(atomic.LoadInt64(&wins[i]))
		t := int(atomic.LoadInt64(&ties[i]))
		l := simulations - w - t

		results[i] = EquityResult{
			Wins:   w,
			Ties:   t,
			Losses: l,
			Total:  simulations,
			Equity: (float64(w) + float64(t)/float64(numPlayers)) / float64(simulations),
		}
	}

	return results
}

func (ec *EquityCalculator) runSimulations(
	holeCards [][]Card,
	board []Card,
	deck []Card,
	cardsNeeded int,
	numSims int,
	wins []int64,
	ties []int64,
) {
	numPlayers := len(holeCards)
	localDeck := make([]Card, len(deck))

	for sim := 0; sim < numSims; sim++ {
		// Copy and shuffle deck
		copy(localDeck, deck)
		shuffleDeck(localDeck)

		// Complete the board
		fullBoard := make([]Card, len(board), 5)
		copy(fullBoard, board)
		fullBoard = append(fullBoard, localDeck[:cardsNeeded]...)

		// Evaluate each player's best hand
		bestHands := make([]*Hand, numPlayers)
		for i, hole := range holeCards {
			allCards := make([]Card, 0, 7)
			allCards = append(allCards, hole...)
			allCards = append(allCards, fullBoard...)

			best := findBestHand(allCards)
			bestHands[i] = best
		}

		// Determine winners
		winnerIndices := findWinnerIndices(bestHands)

		if len(winnerIndices) == 1 {
			atomic.AddInt64(&wins[winnerIndices[0]], 1)
		} else {
			for _, idx := range winnerIndices {
				atomic.AddInt64(&ties[idx], 1)
			}
		}
	}
}

// CalculateExact calculates exact equity by enumerating all possible board runouts.
// Only practical when few cards remain to be dealt (e.g., river only).
// Returns nil if too many combinations (> maxCombinations).
func (ec *EquityCalculator) CalculateExact(holeCards [][]Card, board []Card, maxCombinations int) []EquityResult {
	// Build remaining deck
	usedCards := make(map[string]bool)
	for _, hole := range holeCards {
		for _, c := range hole {
			usedCards[cardKey(c)] = true
		}
	}
	for _, c := range board {
		usedCards[cardKey(c)] = true
	}

	remainingDeck := buildRemainingDeck(usedCards)
	cardsNeeded := 5 - len(board)

	// Check if enumeration is feasible
	combos := Combinations(remainingDeck, cardsNeeded)
	if len(combos) > maxCombinations {
		return nil
	}

	numPlayers := len(holeCards)
	wins := make([]int, numPlayers)
	ties := make([]int, numPlayers)
	total := len(combos)

	// Process combinations in parallel
	type result struct {
		winners []int
	}

	results := make(chan result, len(combos))
	sem := make(chan struct{}, ec.workers)
	var wg sync.WaitGroup

	for _, combo := range combos {
		wg.Add(1)
		go func(boardCards []Card) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			fullBoard := make([]Card, len(board), 5)
			copy(fullBoard, board)
			fullBoard = append(fullBoard, boardCards...)

			bestHands := make([]*Hand, numPlayers)
			for i, hole := range holeCards {
				allCards := make([]Card, 0, 7)
				allCards = append(allCards, hole...)
				allCards = append(allCards, fullBoard...)
				bestHands[i] = findBestHand(allCards)
			}

			results <- result{findWinnerIndices(bestHands)}
		}(combo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if len(r.winners) == 1 {
			wins[r.winners[0]]++
		} else {
			for _, idx := range r.winners {
				ties[idx]++
			}
		}
	}

	equityResults := make([]EquityResult, numPlayers)
	for i := 0; i < numPlayers; i++ {
		l := total - wins[i] - ties[i]
		equityResults[i] = EquityResult{
			Wins:   wins[i],
			Ties:   ties[i],
			Losses: l,
			Total:  total,
			Equity: (float64(wins[i]) + float64(ties[i])/float64(numPlayers)) / float64(total),
		}
	}

	return equityResults
}

// Helper functions

func cardKey(c Card) string {
	return string(rune(c.Rank)) + string(rune(c.Suit))
}

func buildRemainingDeck(usedCards map[string]bool) []Card {
	deck := make([]Card, 0, 52-len(usedCards))
	for _, suit := range AllSuits() {
		for _, rank := range AllRanks() {
			c := NewCard(rank, suit)
			if !usedCards[cardKey(c)] {
				deck = append(deck, c)
			}
		}
	}
	return deck
}

func shuffleDeck(deck []Card) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func findBestHand(cards []Card) *Hand {
	combos := CardCombinations(cards, 5)
	var best *Hand

	for _, combo := range combos {
		hand, err := NewHand(combo)
		if err != nil {
			continue
		}
		if best == nil || hand.Beats(best) {
			best = hand
		}
	}

	return best
}

func findWinnerIndices(hands []*Hand) []int {
	var winners []int
	var bestHand *Hand

	for i, hand := range hands {
		if hand == nil {
			continue
		}

		if bestHand == nil {
			winners = []int{i}
			bestHand = hand
			continue
		}

		cmp := hand.Compare(bestHand)
		if cmp > 0 {
			winners = []int{i}
			bestHand = hand
		} else if cmp == 0 {
			winners = append(winners, i)
		}
	}

	return winners
}
