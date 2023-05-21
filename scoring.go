package main

import (
	"sort"
)

/*

functions for scoring dice

*/

func scoreAces(dice []int) int {
	return scoreNumberedDice(1, dice)
}

func scoreDeuces(dice []int) int {
	return scoreNumberedDice(2, dice)
}

func scoreThrees(dice []int) int {
	return scoreNumberedDice(3, dice)
}

func scoreFours(dice []int) int {
	return scoreNumberedDice(4, dice)
}

func scoreFives(dice []int) int {
	return scoreNumberedDice(5, dice)
}

func scoreSixes(dice []int) int {
	return scoreNumberedDice(6, dice)
}

func scoreNumberedDice(num int, dice []int) int {
	score := 0

	for _, die := range dice {
		if die == num {
			score += num
		}
	}

	return score
}

func scoreFourOfAKind(dice []int) int {
	// check if there are four dice of the same number
	score := 0
	hasFourOfAKind := false

	// check four of a kind by getting score for each number
	if scoreNumberedDice(1, dice) == 4 || scoreNumberedDice(2, dice) == 8 ||
		scoreNumberedDice(3, dice) == 12 || scoreNumberedDice(4, dice) == 16 ||
		scoreNumberedDice(5, dice) == 20 || scoreNumberedDice(6, dice) == 24 {
		hasFourOfAKind = true
	}

	// if they have four of a kind, score is the sum of all dice
	if hasFourOfAKind {
		for _, die := range dice {
			score += die
		}
	}

	return score
}

func scoreFullHouse(dice []int) int {
	// TODO there is a bug here, something like 2 2 2 3 5 caused a full house

	// check if there are three one die and two of another
	score := 0
	sort.Ints(dice)

	// with a sorted list, there are two ways to have a full house
	// 1. die0 == die1 && die2 == die3 == die4
	// 2. die0 == die1 == die2 && die3 == die4
	if (dice[0] == dice[1] && dice[2] == dice[3] && dice[3] == dice[4]) ||
		(dice[0] == dice[1] && dice[1] == dice[2] && dice[3] == dice[4]) {
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

func scoreSmallStraight(dice []int) int {
	// four consecutive numbers
	score := 0
	sort.Ints(dice)

	// only look at unique dice, other wise 1 2 2 3 4 wouldn't count using the logic below
	uniqueDice := []int{}
	for _, die := range dice {
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

func scoreLargeStraight(dice []int) int {
	// five consecutive numbers
	score := 40
	sort.Ints(dice)

	for i := 1; i < len(dice); i++ {
		if dice[i] != dice[i-1]+1 {
			score = 0
			break
		}
	}

	return score
}

func scoreChance(dice []int) int {
	// sum of all dice
	score := 0

	for _, die := range dice {
		score += die
	}

	return score
}

func scoreYacht(dice []int) int {
	// five of a kind
	score := 0

	if dice[0] == dice[1] && dice[1] == dice[2] &&
		dice[2] == dice[3] && dice[3] == dice[4] {
		score = 50
	}

	return score
}
