package main

import (
	"fmt"
	gour "gour/internal"
	"time"
)

func main() {
	start := time.Now()
	ai, err := gour.LoadUrAI("trained/UR_evolving/2/ur_winner_genome_98-349")
	if err != nil {
		panic(err)
	}
	left_player := gour.Ai_ur_player{
		Ai:   ai,
		Name: "AI",
	}
	ai2, err := gour.LoadUrAI("trained/541/ur_winner_genome_58-39")
	right_player := gour.Ai_ur_player{
		Ai:   ai2,
		Name: "AI",
	}
	//right_player := gour.Random_ur_player{
	//	Name: "Random",
	//}

	ai_wins3, _ := gour.OneVSOne(&left_player, &right_player, 3, 3333)
	ai_wins5, _ := gour.OneVSOne(&left_player, &right_player, 5, 3333)
	ai_wins7, _ := gour.OneVSOne(&left_player, &right_player, 7, 3334)
	fmt.Printf("AI won %f times\n", float64(ai_wins3+ai_wins5+ai_wins7)/10000.0)
	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
