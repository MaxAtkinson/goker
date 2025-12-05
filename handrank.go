package goker

// HandRank represents the ranking of a poker hand.
type HandRank int

const (
	HighCard HandRank = iota + 1
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (h HandRank) String() string {
	switch h {
	case HighCard:
		return "High Card"
	case Pair:
		return "Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

// Constants for binary arithmetic hand evaluation
const (
	numBits    = 15
	unusedBits = 2

	highCardValue      = 5
	pairValue          = 6
	twoPairValue       = 7
	threeOfAKindValue  = 9
	straightValue      = 31
	wheelStraightValue = 16444 // A2345 special case
	fullHouseValue     = 10
	fourOfAKindValue   = 17

	royalFlushValue = 31744 // AKQJT bitmap
)

// binaryArithmeticToHandRank maps modulo division results to hand ranks.
var binaryArithmeticToHandRank = map[int]HandRank{
	highCardValue:     HighCard,
	pairValue:         Pair,
	twoPairValue:      TwoPair,
	threeOfAKindValue: ThreeOfAKind,
	fullHouseValue:    FullHouse,
	fourOfAKindValue:  FourOfAKind,
}

// Tiebreaker bit shifts for scoring
var tiebreakerShifts = []int{16, 12, 8, 4, 0}
