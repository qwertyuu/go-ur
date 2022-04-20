package main

import (
	"fmt"
	gour "gour/internal"
	"time"
)

func main() {
	// Analyze the probability of winning a game if you are the first player vs. the second, in case the game is skewed
	start := time.Now()
	left_player := gour.Random_ur_player{
		Name: "Random left",
	}

	right_player := gour.Random_ur_player{
		Name: "Random right",
	}

	number_of_games := 100000

	first_player_wins_3, _ := gour.OneVSOne(&left_player, &right_player, 3, number_of_games)
	fmt.Printf("first player won %f times in a 3 pawn game\n", float64(first_player_wins_3)/float64(number_of_games))
	first_player_wins_5, _ := gour.OneVSOne(&left_player, &right_player, 5, number_of_games)
	fmt.Printf("first player won %f times in a 5 pawn game\n", float64(first_player_wins_5)/float64(number_of_games))
	first_player_wins_7, _ := gour.OneVSOne(&left_player, &right_player, 7, number_of_games)
	fmt.Printf("first player won %f times in a 7 pawn game\n", float64(first_player_wins_7)/float64(number_of_games))
	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
