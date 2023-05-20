package main

import (
	"math/rand"
	"sort"
	"time"
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

// refactor idea: this could be a map like ScoreHints where I only store categories that have a score
// I could have a map with the keys predefined and each key would map to a score function
// then updating score hints would just be a matter of finding keys that aren't in playerscore map
// also the score event handler would do something like PlayerScore[category] = categories[category]() - calling the score func
// that would simplify how this data is consumed on the front end too, would reduce the repetive checks agains the Is<category>Score fields
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
	// seed the RNG
	rand.Seed(time.Now().UnixNano())

	// score card data defualts to 0
	return &Game{
		Round:      1,
		Turn:       "p1",
		RollsLeft:  3,
		DiceKept:   []int{},
		DiceInPlay: []int{1, 1, 1, 1, 1},
		Score:      &ScoreCard{Player1Score: &PlayerScore{}, Player2Score: &PlayerScore{}},
		ScoreHints: make(map[string]int),
	}
}

func (g *Game) rollDice() {
	if g.RollsLeft <= 0 {
		return
	}

	// roll whatever dice are in play
	roll := make([]int, len(g.DiceInPlay))

	for i := 0; i < len(g.DiceInPlay); i++ {
		roll[i] = rand.Intn(6) + 1
	}

	g.DiceInPlay = roll
	g.RollsLeft--

	// update score hints after roll
	g.updateScoreHints()
}

func (g *Game) updateScoreHints() {
	var scoreCard *PlayerScore

	if g.Turn == "p1" {
		scoreCard = g.Score.Player1Score
	} else {
		scoreCard = g.Score.Player2Score
	}

	if !scoreCard.IsAcesScore {
		g.ScoreHints["aces"] = g.scoreNumberedDice(1)
	}

	if !scoreCard.IsDeucesScore {
		g.ScoreHints["deuces"] = g.scoreNumberedDice(2)
	}

	if !scoreCard.IsThreesScore {
		g.ScoreHints["threes"] = g.scoreNumberedDice(3)
	}

	if !scoreCard.IsFoursScore {
		g.ScoreHints["fours"] = g.scoreNumberedDice(4)
	}

	if !scoreCard.IsFivesScore {
		g.ScoreHints["fives"] = g.scoreNumberedDice(5)
	}

	if !scoreCard.IsSixesScore {
		g.ScoreHints["sixes"] = g.scoreNumberedDice(6)
	}

	if !scoreCard.IsFourOfAKindScore {
		g.ScoreHints["fourOfAKind"] = g.scoreFourOfAKind()
	}

	if !scoreCard.IsFullHouseScore {
		g.ScoreHints["fullHouse"] = g.scoreFullHouse()
	}

	if !scoreCard.IsSmallStraightScore {
		g.ScoreHints["smallStraight"] = g.scoreSmallStraight()
	}

	if !scoreCard.IsLargeStraightScore {
		g.ScoreHints["largeStraight"] = g.scoreLargeStraight()
	}

	if !scoreCard.IsChanceScore {
		g.ScoreHints["chance"] = g.scoreChance()
	}

	if !scoreCard.IsYachtScore {
		g.ScoreHints["yacht"] = g.scoreYacht()
	}

}

func (g *Game) keepDie(index int) {
	g.DiceKept = append(g.DiceKept, g.DiceInPlay[index])

	// remove the kept die from the dice in play

	newDiceInPlay := []int{}
	for i := 0; i < len(g.DiceInPlay); i++ {
		if i == index {
			continue
		}

		newDiceInPlay = append(newDiceInPlay, g.DiceInPlay[i])
	}

	g.DiceInPlay = newDiceInPlay
}

func (g *Game) unkeepDie(index int) {
	g.DiceInPlay = append(g.DiceInPlay, g.DiceKept[index])

	// remove the kept die from the dice in play

	newDiceKept := []int{}
	for i := 0; i < len(g.DiceKept); i++ {
		if i == index {
			continue
		}

		newDiceKept = append(newDiceKept, g.DiceKept[i])
	}

	g.DiceKept = newDiceKept
}

func (g *Game) scoreNumberedDice(num int) int {
	allDice := append(g.DiceInPlay, g.DiceKept...)
	score := 0

	for _, die := range allDice {
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
		for _, die := range allDice {
			score += die
		}
	}

	return score
}

func (g *Game) scoreFullHouse() int {
	// TODO there is a bug here, something like 2 2 2 3 5 caused a full house

	// check if there are three one die and two of another
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)
	sort.Ints(allDice)

	// with a sorted list, there are two ways to have a full house
	// 1. die0 == die1 && die2 == die3 == die4
	// 2. die0 == die1 == die2 && die3 == die4
	if (allDice[0] == allDice[1] && allDice[2] == allDice[3] && allDice[3] == allDice[4]) ||
		(allDice[0] == allDice[1] && allDice[1] == allDice[2] && allDice[3] == allDice[4]) {
		score = 25
	}

	return score
}

func arrContains(arr []int, val int) bool {
	for i := range arr {
		if arr[i] == val {
			return true
		}
	}

	return false
}

func (g *Game) scoreSmallStraight() int {
	// four consecutive numbers
	score := 0
	allDice := append(g.DiceInPlay, g.DiceKept...)
	sort.Ints(allDice)

	// only look at unique dice, other wise 1 2 2 3 4 wouldn't count using the logic below
	uniqueDice := []int{}
	for _, die := range allDice {
		if !arrContains(uniqueDice, die) {
			uniqueDice = append(uniqueDice, die)
		}
	}

	numInARow := 1

	for i := 1; i < len(uniqueDice); i++ {
		if uniqueDice[i] == uniqueDice[i-1]+1 {
			numInARow++
		} else {
			numInARow = 1
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

	for _, die := range allDice {
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

	// put all dice back in play
	g.DiceInPlay = append(g.DiceInPlay, g.DiceKept...)
	g.DiceKept = []int{}

	// reset rolls left
	g.RollsLeft = 3

	// reset score hints
	g.ScoreHints = make(map[string]int)
}
