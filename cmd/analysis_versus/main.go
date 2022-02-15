package main

import (
	"fmt"
	gour "gour/internal"
)

func main() {
	//ai, err := gour.LoadUrAI("out/UR_best_iteration_1/2/ur_winner_genome_58-31")
	ai, err := gour.LoadUrAI("trained/197/ur_winner_genome_59-31")
	if err != nil {
		panic(err)
	}
	left_player := gour.Ai_ur_player{
		Ai:   ai,
		Name: "AI",
	}
	right_player := gour.Random_ur_player{
		Name: "First move picker",
	}

	ai_wins, _ := gour.OneVSOne(&left_player, &right_player, 7, 10000)
	fmt.Printf("AI won %f times", float64(ai_wins)/10000.0)

}