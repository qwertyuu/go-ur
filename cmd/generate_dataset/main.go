package main

import (
	"encoding/csv"
	"fmt"
	gour "gour/internal"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ai, err := gour.LoadUrAI("trained/UR_evolving/2/ur_winner_genome_98-349")
	if err != nil {
		panic(err)
	}

	left_player := &gour.Verbose_ai_ur_player{
		Ai_ur_player: gour.Ai_ur_player{
			Ai:   ai,
			Name: "AI",
		},
		Played_vectors: make([][]float64, 0),
	}
	//ai2, err := gour.LoadUrAI("trained/541/ur_winner_genome_58-39")
	//right_player := gour.Ai_ur_player{
	//	Ai:   ai2,
	//	Name: "AI",
	//}
	right_player := &gour.Random_ur_player{
		Name: "Random",
	}
	amt_pawns := 7
	amt_games := 200

	for i := 0; i < amt_games; i++ {
		board := gour.NewBoard(amt_pawns)
		gour.FightUntilWon(board, left_player, right_player)
	}

	log.Println(len(left_player.Played_vectors))
	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)

	records := [][]string{
		gour.GetFeatureNames(),
	}

	for _, row := range left_player.Played_vectors {
		csv_strings := []string{}
		for _, value := range row {
			csv_strings = append(csv_strings, fmt.Sprint(value))
		}
		records = append(records, csv_strings)
	}

	f, err := os.Create(fmt.Sprintf("dataset_AI_%s_%v.csv", right_player.GetName(), amt_pawns))
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}
}
