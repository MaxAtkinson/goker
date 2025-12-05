package goker

import "errors"

var (
	// ErrEmptyDeck is returned when attempting to draw from an empty deck.
	ErrEmptyDeck = errors.New("cannot draw from empty deck")

	// ErrInvalidHandSize is returned when a hand doesn't have exactly 5 cards.
	ErrInvalidHandSize = errors.New("hand must contain exactly 5 cards")

	// ErrDuplicateCards is returned when a hand contains duplicate cards.
	ErrDuplicateCards = errors.New("hand contains duplicate cards")

	// ErrInvalidHoleCards is returned when hole cards don't have exactly 2 cards.
	ErrInvalidHoleCards = errors.New("hole cards must contain exactly 2 cards")

	// ErrInvalidBoardState is returned when board operation is invalid for current state.
	ErrInvalidBoardState = errors.New("invalid board state for this operation")
)
