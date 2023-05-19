package main

import (
	"math/rand"
	"sort"
)

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

func (g *Game) rollDice() {
	// roll whatever dice are in play
	roll := make([]int, len(g.DiceInPlay))

	for i := 0; i < len(g.DiceInPlay); i++ {
		roll[i] = rand.Intn(6) + 1
	}

	g.DiceInPlay = roll

	// update score hints after roll
}

func (g *Game) keepDie(dice int) {
	g.DiceKept = append(g.DiceKept, dice)
}

func (g *Game) scoreNumberedDice(num int) int {
	allDice := append(g.DiceInPlay, g.DiceKept...)
	score := 0

	for die := range allDice {
		if die == num {
			score += num
		}
	}

	return score
}

func (g *Game) scoreFourOfAKind() int {
	// check if there are four dice of the same number
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)
	hasFourOfAKind := false

	// check four of a kind by getting score for each number
	if g.scoreNumberedDice(1) == 4 || g.scoreNumberedDice(2) == 8 ||
		g.scoreNumberedDice(3) == 12 || g.scoreNumberedDice(4) == 16 ||
		g.scoreNumberedDice(5) == 20 || g.scoreNumberedDice(6) == 24 {
		hasFourOfAKind = true
	}

	// if they have four of a kind, score is the sum of all dice
	if hasFourOfAKind {
		for die := range allDice {
			score += die
		}
	}

	return score
}

func (g *Game) scoreFullHouse() int {
	// check if there are three one die and two of another
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)
	sort.Ints(allDice)

	// with a sorted list, there are two ways to have a full house
	// 1. die0 == die1 && die2 == die3 == die4
	// 2. die0 == die1 == die2 && die3 == die4
	if allDice[0] == allDice[1] &&
		(allDice[0] == allDice[2] ||
			(allDice[2] == allDice[3] && allDice[3] == allDice[4])) {
		score = 25
	}

	return score
}

func (g *Game) scoreSmallStraight() int {
	// four consecutive numbers
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)
	sort.Ints(allDice)

	numInARow := 0

	for i := 1; i < len(allDice); i++ {
		if allDice[i] == allDice[i-1]+1 {
			numInARow++
		} else {
			numInARow = 0
		}

		if numInARow == 4 {
			score = 30
			break
		}
	}

	return score
}

func (g *Game) scoreLargeStraight() int {
	// five consecutive numbers
	score := 40
	allDice := append(g.DiceInPlay, g.DiceKept...)
	sort.Ints(allDice)

	for i := 1; i < len(allDice); i++ {
		if allDice[i] != allDice[i-1]+1 {
			score = 0
			break
		}
	}

	return score
}

func (g *Game) scoreChance() int {
	// sum of all dice
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)

	for die := range allDice {
		score += die
	}

	return score
}

func (g *Game) scoreYacht() int {
	// five of a kind
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)

	if allDice[0] == allDice[1] && allDice[1] == allDice[2] &&
		allDice[2] == allDice[3] && allDice[3] == allDice[4] {
		score = 50
	}

	return score
}

func (g *Game) scoreRoll(category string) {
	var scoreCard *PlayerScore

	if g.Turn == "p1" {
		scoreCard = g.Score.Player1Score
	} else {
		scoreCard = g.Score.Player2Score
	}

	switch category {
	case "aces":
		scoreCard.Aces = g.scoreNumberedDice(1)
		scoreCard.IsAcesScore = true
	case "deuces":
		scoreCard.Deuces = g.scoreNumberedDice(2)
		scoreCard.IsDeucesScore = true
	case "threes":
		scoreCard.Threes = g.scoreNumberedDice(3)
		scoreCard.IsThreesScore = true
	case "fours":
		scoreCard.Fours = g.scoreNumberedDice(4)
		scoreCard.IsFoursScore = true
	case "fives":
		scoreCard.Fives = g.scoreNumberedDice(5)
		scoreCard.IsFivesScore = true
	case "sixes":
		scoreCard.Sixes = g.scoreNumberedDice(6)
		scoreCard.IsSixesScore = true
	case "fourOfAKind":
		scoreCard.FourOfAKind = g.scoreFourOfAKind()
		scoreCard.IsFourOfAKindScore = true
	case "fullHouse":
		scoreCard.FullHouse = g.scoreFullHouse()
		scoreCard.IsFullHouseScore = true
	case "smallStraight":
		scoreCard.SmallStraight = g.scoreSmallStraight()
		scoreCard.IsSmallStraightScore = true
	case "largeStraight":
		scoreCard.LargeStraight = g.scoreLargeStraight()
		scoreCard.IsLargeStraightScore = true
	case "chance":
		scoreCard.Chance = g.scoreChance()
		scoreCard.IsChanceScore = true
	case "yacht":
		scoreCard.Yacht = g.scoreYacht()
		scoreCard.IsYachtScore = true
	}

	// switch turns
	if g.Turn == "p1" {
		g.Turn = "p2"
	} else {
		g.Turn = "p1"

		// round increments when we come back to player 1
		g.Round++
	}
}
