package main

type Game struct {
	Round      int            `json:"round"`
	Winner     string         `json:"winner"` // "p1" or "p2", game is over if this is not null
	Turn       string         `json:"turn"`   // "p1" or "p2"
	Player1    *Player        `json:"p1"`
	Player2    *Player        `json:"p2"`
	RollsLeft  int            `json:"rollsLeft"`
	DiceKept   []int          `json:"diceKept"`
	DiceInPlay []int          `json:"diceInPlay"`
	Score      *ScoreCard     `json:"scoreCard"`
	ScoreHints map[string]int `json:"scoreHints"` // only populated with possible selections
}

type ScoreCard struct {
	Player1Score *PlayerScore `json:"p1"`
	Player2Score *PlayerScore `json:"p2"`
}

// since the ints default to 0, the booleans tell whether the 0 is there actual score
// or if the player just hasn't entered a score for that category yet
type PlayerScore struct {
	Aces                 int  `json:"aces"`
	IsAcesScore          bool `json:"isAcesScore"`
	Deuces               int  `json:"deuces"`
	IsDeucesScore        bool `json:"isDeucesScore"`
	Threes               int  `json:"threes"`
	IsThreesScore        bool `json:"isThreesScore"`
	Fours                int  `json:"fours"`
	IsFoursScore         bool `json:"isFoursScore"`
	Fives                int  `json:"fives"`
	IsFivesScore         bool `json:"isFivesScore"`
	Sixes                int  `json:"sixes"`
	IsSixesScore         bool `json:"isSixesScore"`
	FourOfAKind          int  `json:"fourOfAKind"`
	IsFourOfAKindScore   bool `json:"isFourOfAKindScore"`
	FullHouse            int  `json:"fullHouse"`
	IsFullHouseScore     bool `json:"isFullHouseScore"`
	SmallStraight        int  `json:"smallStraight"`
	IsSmallStraightScore bool `json:"isSmallStraightScore"`
	LargeStraight        int  `json:"largeStraight"`
	IsLargeStraightScore bool `json:"isLargeStraightScore"`
	Chance               int  `json:"chance"`
	IsChanceScore        bool `json:"isChanceScore"`
	Yacht                int  `json:"yacht"`
	IsYachtScore         bool `json:"isYachtScore"`
}

func NewGame() *Game {
	// score card data defualts to 0
	return &Game{
		Round:      1,
		Turn:       "p1",
		RollsLeft:  3,
		DiceKept:   []int{},
		DiceInPlay: []int{1, 1, 1, 1, 1},
		Score:      &ScoreCard{Player1Score: &PlayerScore{}, Player2Score: &PlayerScore{}},
	}
}
