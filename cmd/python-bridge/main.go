package main

import "C"
import (
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"log"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type board_contract struct {
	Pawn_per_player      int   `json:"pawn_per_player"`
	AI_pawn_out          int   `json:"ai_pawn_out"`
	Enemy_pawn_out       int   `json:"enemy_pawn_out"`
	Dice                 int   `json:"dice"`
	AI_pawn_positions    []int `json:"ai_pawn_positions"`
	Enemy_pawn_positions []int `json:"enemy_pawn_positions"`
}

var ai *genetics.Organism

//export infer
func infer(board_json_c *C.char) *C.char {
	var err error
	if ai == nil {
		log.Println("Loading AI")
		ai, err = gour.LoadUrAI("trained/UR_evolving/2/ur_winner_genome_98-349")
		if err != nil {
			panic(err)
		}
	}
	log.Println("Hello from infer")
	var board_input board_contract
	board_json := C.GoString(board_json_c)
	err = json.Unmarshal([]byte(board_json), &board_input)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return C.CString("")
	}
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
	//fmt.Println(board.String())
	fmt.Println(board.Current_player_path_moves)
	potential_futures := gour.GetMoveScoresOrdered(board, ai)

	inference, err := json.Marshal(struct {
		Pawn         int                      `json:"pawn"`
		FutureScores []*gour.Potential_future `json:"future_scores"`
	}{
		Pawn:         potential_futures[len(potential_futures)-1].Pawn,
		FutureScores: potential_futures,
	})
	if err != nil {
		log.Printf("Error creating inference output: %v", err)
		return C.CString("")
	}
	return C.CString(string(inference))
}

func main() {}
