package gour

import (
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type Ur_player interface {
	GetMove(*board) int
	GetName() string
	GetType() string
	SetWinner(bool)
	IncrementWins(int)
	GetWins() int
	Copy() Ur_player
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
	rand_key := pcg32.Bounded(uint32(len(*board.Current_player_path_moves)))
	count := uint32(0)
	last_k := -1
	for k := range *board.Current_player_path_moves {
		if count == rand_key {
			return k
		}
		last_k = k
		count++
	}
	return last_k
}

func (s *Random_ur_player) GetName() string {
	return s.Name
}

func (s *Random_ur_player) GetType() string {
	return "RANDOM"
}

func (s *Random_ur_player) Copy() Ur_player {
	return &Random_ur_player{
		Name: s.Name,
		wins: s.wins,
	}
}

type First_move_ur_player struct {
	Ur_player
	Name string
}

func (s *First_move_ur_player) GetMove(board *board) int {
	for k := range *board.Current_player_path_moves {
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

func (s *First_move_ur_player) Copy() Ur_player {
	return &First_move_ur_player{
		Name: s.Name,
	}
}

type Verbose_ai_ur_player struct {
	Ai_ur_player
	Played_vectors [][]float64
}

func (s *Verbose_ai_ur_player) GetMove(board *board) int {
	vectors := GetMovesVectors(board)
	s.Played_vectors = append(s.Played_vectors, vectors...)
	return s.Ai_ur_player.GetMove(board)
}

func (s *Verbose_ai_ur_player) Copy() Ur_player {
	ai := s.Ai_ur_player.Copy().(*Ai_ur_player)
	cpy := make([][]float64, len(s.Played_vectors))
	copy(cpy, s.Played_vectors)
	c := &Verbose_ai_ur_player{
		Ai_ur_player:   *ai,
		Played_vectors: cpy,
	}
	return c
}

type Ai_ur_player struct {
	Ur_player
	Ai   *genetics.Organism
	Name string
}

func (s *Ai_ur_player) GetMove(board *board) int {
	potential_futures, _ := GetMoveScoresOrdered(board, s.Ai)
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

func (s *Ai_ur_player) Copy() Ur_player {
	g, _ := s.Ai.Genotype.Duplicate(s.Ai.Genotype.Id)
	ai, _ := genetics.NewOrganism(s.Ai.Fitness, g, s.Ai.Generation)
	c := &Ai_ur_player{
		Name: s.Name,
		Ai:   ai,
	}
	return c
}
