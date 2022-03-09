package gour

import (
	"math/rand"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type Ur_player interface {
	GetMove(*board) int
	GetName() string
	GetType() string
	SetWinner(bool)
	IncrementWins(int)
	GetWins() int
}

type Random_ur_player struct {
	Ur_player
	Name string
	wins int
}

func (s *Random_ur_player) IncrementWins(wins int) {
	s.wins += wins
}

func (s *Random_ur_player) GetMove(board *board) int {
	all_pawns := []int{}
	for k := range board.Current_player_path_moves {
		all_pawns = append(all_pawns, k)
	}
	return all_pawns[rand.Intn(len(all_pawns))]
}

func (s *Random_ur_player) GetName() string {
	return s.Name
}

func (s *Random_ur_player) GetType() string {
	return "RANDOM"
}

type First_move_ur_player struct {
	Ur_player
	Name string
}

func (s *First_move_ur_player) GetMove(board *board) int {
	for k := range board.Current_player_path_moves {
		return k
	}
	return -2 // never happens
}

func (s *First_move_ur_player) GetName() string {
	return s.Name
}

func (s *First_move_ur_player) GetType() string {
	return "FIRST_MOVE"
}

type Ai_ur_player struct {
	Ur_player
	Ai   *genetics.Organism
	Name string
}

func (s *Ai_ur_player) GetMove(board *board) int {
	potential_futures := GetMoveScoresOrdered(board, s.Ai)
	return potential_futures[len(potential_futures)-1].Pawn
}

func (s *Ai_ur_player) GetName() string {
	return s.Name
}

func (s *Ai_ur_player) GetType() string {
	return "NEAT"
}

func (s *Ai_ur_player) GetWins() int {
	return int(s.Ai.Fitness)
}

func (s *Ai_ur_player) IncrementWins(wins int) {
	s.Ai.Fitness += float64(wins)
}

func (s *Ai_ur_player) SetWinner(winner bool) {
	s.Ai.IsWinner = winner
}
