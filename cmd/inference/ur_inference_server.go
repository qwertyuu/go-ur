package main

import (
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

func main() {
	file, err := os.Open("trained/ur_winner_genome_57-65")
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

	http.HandleFunc("/infer", infer)

	http.ListenAndServe(":8090", nil)
}

type board_contract struct {
	Pawn_per_player      int   `json:"pawn_per_player"`
	My_pawn_out          int   `json:"my_pawn_out"`
	Enemy_pawn_out       int   `json:"enemy_pawn_out"`
	Dice                 int   `json:"dice"`
	My_pawn_positions    []int `json:"my_pawn_positions"`
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
		board_input.My_pawn_out,
		board_input.Enemy_pawn_out,
		gour.Left,
		board_input.Dice,
		board_input.My_pawn_positions,
		board_input.Enemy_pawn_positions,
	)
	fmt.Println(board.String())
	fmt.Println(board.Current_player_path_moves)
	current_board_descriptor := gour.GetCurrentBoardDescriptor(board, gour.Left)
	potential_futures := []*gour.Potential_future{}
	for pawn := range board.Current_player_path_moves {
		potential_game := board.Copy()
		potential_game.Play(pawn)
		fmt.Println(potential_game.String())
		potential_board := gour.GetPotentialBoardDescriptor(potential_game, board.Current_player)
		score, err := gour.GetPotentialFutureScore(ai, current_board_descriptor, potential_board)
		fmt.Println(score)
		if err != nil {
			panic(err)
		}
		potential_futures = append(potential_futures, &gour.Potential_future{
			Pawn:  pawn,
			Score: score,
		})
	}
	sort.Slice(potential_futures, func(i, j int) bool {
		return potential_futures[i].Score < potential_futures[j].Score
	})

	json.NewEncoder(w).Encode(struct {
		Pawn         int                      `json:"pawn"`
		FutureScores []*gour.Potential_future `json:"future_scores"`
	}{
		Pawn:         potential_futures[len(potential_futures)-1].Pawn,
		FutureScores: potential_futures,
	})
}
