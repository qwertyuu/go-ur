package main

import (
	"fmt"
	gour "gour/internal"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

var ai *genetics.Organism

func main() {
	ai, err := gour.LoadUrAI("trained/ur_winner_genome_61-206")
	if err != nil {
		panic(err)
	}
	left_player := gour.Ai_ur_player{
		Ai:   ai,
		Name: "AI",
	}
	right_player := gour.First_move_ur_player{
		Name: "First move picker",
	}

	ai_wins, _ := OneVSOne(&left_player, &right_player, 1000)
	fmt.Printf("AI won %f times", float64(ai_wins)/1000.0)

}

func OneVSOne(left_player gour.Ur_player, right_player gour.Ur_player, number_of_games int) (int, int) {
	left_wins := 0
	right_wins := 0
	for i := 0; i < number_of_games; i++ {
		board := gour.NewBoard(7)
		moves := 0
		for board.Current_winner == 0 {
			var current_player gour.Ur_player
			if board.Current_player == gour.Left {
				current_player = left_player
			} else {
				current_player = right_player
			}
			board.Play(current_player.GetMove(board))
			moves++
		}
		if board.Current_winner == gour.Left {
			fmt.Printf("%s wins after %d moves\n", left_player.GetName(), moves)
			left_wins++
		} else {
			fmt.Printf("%s wins after %d moves\n", right_player.GetName(), moves)
			right_wins++
		}
	}
	return left_wins, right_wins
}
