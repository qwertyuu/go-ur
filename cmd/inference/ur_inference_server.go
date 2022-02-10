package main

import (
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"log"
	"net/http"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

func main() {
	var err error
	ai, err = gour.LoadUrAI("trained/275/ur_winner_genome_55-11")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/infer", infer)

	http.ListenAndServe(":8090", nil)
}

type board_contract struct {
	Pawn_per_player      int   `json:"pawn_per_player"`
	AI_pawn_out          int   `json:"ai_pawn_out"`
	Enemy_pawn_out       int   `json:"enemy_pawn_out"`
	Dice                 int   `json:"dice"`
	AI_pawn_positions    []int `json:"ai_pawn_positions"`
	Enemy_pawn_positions []int `json:"enemy_pawn_positions"`
}

var ai *genetics.Organism

func infer(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	var board_input board_contract
	err := json.NewDecoder(r.Body).Decode(&board_input)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	board := gour.RestoreBoard(
		board_input.Pawn_per_player,
		board_input.AI_pawn_out,
		board_input.Enemy_pawn_out,
		gour.Left,
		board_input.Dice,
		board_input.AI_pawn_positions,
		board_input.Enemy_pawn_positions,
	)
	board.Mirror_print_mode = true
	fmt.Println(board.String())
	fmt.Println(board.Current_player_path_moves)
	potential_futures := gour.GetMoveScoresOrdered(board, ai)

	json.NewEncoder(w).Encode(struct {
		Pawn         int                      `json:"pawn"`
		FutureScores []*gour.Potential_future `json:"future_scores"`
	}{
		Pawn:         potential_futures[len(potential_futures)-1].Pawn,
		FutureScores: potential_futures,
	})
}
