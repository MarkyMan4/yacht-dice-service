package main

import (
	"math/rand"
	"time"
)

type scoreFunc func(dice []int) int

var categories map[string]scoreFunc = map[string]scoreFunc{
	"aces":          scoreAces,
	"deuces":        scoreDeuces,
	"threes":        scoreThrees,
	"fours":         scoreFours,
	"fives":         scoreFives,
	"sixes":         scoreSixes,
	"fourOfAKind":   scoreFourOfAKind,
	"fullHouse":     scoreFullHouse,
	"smallStraight": scoreSmallStraight,
	"largeStraight": scoreLargeStraight,
	"chance":        scoreChance,
	"yacht":         scoreYacht,
}

type Game struct {
	Round      int            `json:"round"`
	Turn       string         `json:"turn"` // "p1" or "p2"
	Player1    *Player        `json:"p1"`
	Player2    *Player        `json:"p2"`
	RollsLeft  int            `json:"rollsLeft"`
	DiceKept   []int          `json:"diceKept"`
	DiceInPlay []int          `json:"diceInPlay"`
	ScoreCard  *PlayerScores  `json:"scoreCard"`
	ScoreHints map[string]int `json:"scoreHints"` // only populated with possible selections
	Winner     string         `json:"winner"`     // "p1" or "p2", game is over if this is not null
	Totals     *PlayerTotals  `json:"totals"`
}

type PlayerScores struct {
	Player1Score map[string]int `json:"p1"`
	Player2Score map[string]int `json:"p2"`
}

type PlayerTotals struct {
	Player1Total int `json:"p1"`
	Player2Total int `json:"p2"`
}

// refactor idea: this could be a map like ScoreHints where I only store categories that have a score
// I could have a map with the keys predefined and each key would map to a score function
// then updating score hints would just be a matter of finding keys that aren't in playerscore map
// also the score event handler would do something like PlayerScore[category] = categories[category]() - calling the score func
// that would simplify how this data is consumed on the front end too, would reduce the repetive checks agains the Is<category>Score fields

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
		ScoreCard: &PlayerScores{
			Player1Score: make(map[string]int),
			Player2Score: make(map[string]int),
		},
		ScoreHints: make(map[string]int),
		Totals:     &PlayerTotals{},
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
	var scoreCard map[string]int

	if g.Turn == "p1" {
		scoreCard = g.ScoreCard.Player1Score
	} else {
		scoreCard = g.ScoreCard.Player2Score
	}

	for cat := range categories {
		// if the category is not in scores, add a score hing
		if _, ok := scoreCard[cat]; !ok {
			allDice := append(g.DiceInPlay, g.DiceKept...)
			g.ScoreHints[cat] = categories[cat](allDice)
		}
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

func (g *Game) updatePlayerTotals() {
	p1Total := 0
	p2Total := 0

	for _, v := range g.ScoreCard.Player1Score {
		p1Total += v
	}

	for _, v := range g.ScoreCard.Player2Score {
		p2Total += v
	}

	g.Totals.Player1Total = p1Total
	g.Totals.Player2Total = p2Total
}

func (g *Game) endGame() {
	// determine the winner and end the game
	if g.Totals.Player1Total > g.Totals.Player2Total {
		g.Winner = "p1"
	} else if g.Totals.Player2Total > g.Totals.Player1Total {
		g.Winner = "p2"
	} else {
		g.Winner = "tie"
	}

	g.DiceInPlay = []int{}
	g.DiceKept = []int{}
}

func (g *Game) scoreRoll(category string) {
	var scoreCard map[string]int

	if g.Turn == "p1" {
		scoreCard = g.ScoreCard.Player1Score
	} else {
		scoreCard = g.ScoreCard.Player2Score
	}

	scorer, categoryExists := categories[category]
	_, hasScore := scoreCard[category]

	// make sure the value in the event is a valid catetory
	// also make sure the player doesn't already have a score for that category
	if categoryExists && !hasScore {
		allDice := append(g.DiceInPlay, g.DiceKept...)
		scoreCard[category] = scorer(allDice)
	}

	g.updatePlayerTotals()

	// if that was the last round (12), determine the winner
	if g.Round == 12 && g.Turn == "p2" {
		g.endGame()
		return
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
