package main

import (
	"fmt"
	gour "gour/internal"
)

func main() {
	ai, err := gour.LoadUrAI("trained/10/ur_winner_genome_58-141")
	if err != nil {
		panic(err)
	}
	left_player := gour.Ai_ur_player{
		Ai:   ai,
		Name: "AI",
	}
	//ai2, err := gour.LoadUrAI("trained/10/ur_winner_genome_58-141")
	//right_player := gour.Ai_ur_player{
	//	Ai:   ai2,
	//	Name: "AI",
	//}
	right_player := gour.Random_ur_player{
		Name: "Random",
	}

	ai_wins, _ := gour.OneVSOne(&left_player, &right_player, 7, 10000)
	fmt.Printf("AI won %f times", float64(ai_wins)/10000.0)

}