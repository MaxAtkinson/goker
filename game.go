package goker

import "fmt"

// Game represents a Texas Hold'em poker game.
type Game struct {
	Deck    *Deck
	Board   *Board
	Players []*Player
}

// NewGame creates a new game with the specified number of players.
func NewGame(numPlayers int) *Game {
	g := &Game{
		Deck:    NewDeck(),
		Board:   NewBoard(),
		Players: make([]*Player, numPlayers),
	}

	for i := 0; i < numPlayers; i++ {
		g.Players[i] = NewPlayer(fmt.Sprintf("Player %d", i+1))
	}

	g.DealHoleCards()
	return g
}

// DealHoleCards deals 2 cards to each player.
func (g *Game) DealHoleCards() error {
	for _, player := range g.Players {
		cards, err := g.Deck.DrawMany(2)
		if err != nil {
			return err
		}
		if err := player.SetHoleCards(cards); err != nil {
			return err
		}
	}
	return nil
}

// DealFlop deals the flop (3 community cards) with a burn.
func (g *Game) DealFlop() error {
	if g.Board.State() != Preflop {
		return ErrInvalidBoardState
	}
	if err := g.Deck.Burn(); err != nil {
		return err
	}
	cards, err := g.Deck.DrawMany(3)
	if err != nil {
		return err
	}
	return g.Board.SetFlop(cards)
}

// DealTurn deals the turn card with a burn.
func (g *Game) DealTurn() error {
	if g.Board.State() != Flop {
		return ErrInvalidBoardState
	}
	if err := g.Deck.Burn(); err != nil {
		return err
	}
	card, err := g.Deck.Draw()
	if err != nil {
		return err
	}
	return g.Board.SetTurn(card)
}

// DealRiver deals the river card with a burn.
func (g *Game) DealRiver() error {
	if g.Board.State() != Turn {
		return ErrInvalidBoardState
	}
	if err := g.Deck.Burn(); err != nil {
		return err
	}
	card, err := g.Deck.Draw()
	if err != nil {
		return err
	}
	return g.Board.SetRiver(card)
}

// DealNextStreet advances to the next stage of the game.
func (g *Game) DealNextStreet() error {
	switch g.Board.State() {
	case Preflop:
		return g.DealFlop()
	case Flop:
		return g.DealTurn()
	case Turn:
		return g.DealRiver()
	default:
		return ErrInvalidBoardState
	}
}

// GetCandidateHands returns all possible 5-card hands for a player.
func (g *Game) GetCandidateHands(player *Player) ([]*Hand, error) {
	if len(g.Board.Cards) < 3 {
		return nil, ErrInvalidBoardState
	}

	// Combine hole cards with board cards
	allCards := make([]Card, 0, len(player.HoleCards)+len(g.Board.Cards))
	allCards = append(allCards, player.HoleCards...)
	allCards = append(allCards, g.Board.Cards...)

	// Generate all 5-card combinations
	combos := CardCombinations(allCards, 5)
	hands := make([]*Hand, 0, len(combos))

	for _, combo := range combos {
		hand, err := NewHandWithPlayer(combo, player)
		if err != nil {
			continue // Skip invalid hands
		}
		hands = append(hands, hand)
	}

	return hands, nil
}

// GetBestHand returns the best possible hand for a player.
func (g *Game) GetBestHand(player *Player) (*Hand, error) {
	hands, err := g.GetCandidateHands(player)
	if err != nil {
		return nil, err
	}
	if len(hands) == 0 {
		return nil, ErrInvalidBoardState
	}

	best := hands[0]
	for _, hand := range hands[1:] {
		if hand.Beats(best) {
			best = hand
		}
	}
	return best, nil
}

// GetWinners returns the winning player(s) with their best hands.
// Returns multiple players in case of a split pot.
func (g *Game) GetWinners() ([]*Player, []*Hand, error) {
	if g.Board.State() != River {
		return nil, nil, ErrInvalidBoardState
	}

	type playerHand struct {
		player *Player
		hand   *Hand
	}

	bestHands := make([]playerHand, 0, len(g.Players))

	for _, player := range g.Players {
		best, err := g.GetBestHand(player)
		if err != nil {
			return nil, nil, err
		}
		bestHands = append(bestHands, playerHand{player, best})
	}

	// Find the winning hand(s)
	var winners []playerHand
	for _, ph := range bestHands {
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
