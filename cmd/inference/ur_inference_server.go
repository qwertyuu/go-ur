package main

import (
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"log"
	"net/http"
	"sync"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type ai_pool struct {
	total          int
	org_pool_count int
	org_pool       [100]*genetics.Organism
}

func NewPool() ai_pool {
	return ai_pool{
		org_pool_count: 0,
	}
}

func (s *ai_pool) Get() *genetics.Organism {
	m.Lock()
	defer m.Unlock()
	if s.org_pool_count == 0 {
		ai, err := gour.LoadUrAI("trained/UR_evolving/2/ur_winner_genome_98-349")
		if err != nil {
			panic(err)
		}
		s.total++
		return ai
	}
	s.org_pool_count--
	return s.org_pool[s.org_pool_count]
}

func (s *ai_pool) Return(org *genetics.Organism) {
	m.Lock()
	defer m.Unlock()
	s.org_pool[s.org_pool_count] = org
	s.org_pool_count++
}

var pool ai_pool
var m sync.Mutex

func main() {
	pool = NewPool()
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
	ai := pool.Get()
	potential_futures, err := gour.GetMoveScoresOrdered(board, ai)
	pool.Return(ai)
	if err != nil || len(potential_futures) == 0 {
		fmt.Println(r.Body)
		panic(err)
	}

	json.NewEncoder(w).Encode(struct {
		Pawn         int                      `json:"pawn"`
		FutureScores []*gour.Potential_future `json:"future_scores"`
	}{
		Pawn:         potential_futures[len(potential_futures)-1].Pawn,
		FutureScores: potential_futures,
	})
}
