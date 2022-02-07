package main

import (
	"fmt"
	gour "gour/internal"
	"os"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

var ai *genetics.Organism

func main() {
	file, err := os.Open("trained/ur_winner_genome_61-206")
	if err != nil {
		panic(err)
	}

	genome, err := genetics.ReadGenome(file, 1)
	if err != nil {
		panic(err)
	}

	file.Close()

	ai, err = genetics.NewOrganism(0, genome, 0)
	if err != nil {
		panic(err)
	}
	ai_player := gour.Ai_ur_player{
		Ai: ai,
	}
	random_player := gour.First_move_ur_player{}

	ai_wins := 0.0
	random_wins := 0.0
	for i := 0; i < 1000; i++ {
		board := gour.NewBoard(7)
		moves := 0
		for board.Current_winner == 0 {
			if board.Current_player == gour.Left {
				board.Play(ai_player.GetMove(board))
			} else {
				board.Play(random_player.GetMove(board))
			}
			moves++
		}
		if board.Current_winner == gour.Left {
			fmt.Printf("AI Wins after %d moves\n", moves)
			ai_wins++
		} else {
			fmt.Printf("Random Wins after %d moves\n", moves)
			random_wins++
		}
	}
	fmt.Printf("AI won %f times", ai_wins)

}
