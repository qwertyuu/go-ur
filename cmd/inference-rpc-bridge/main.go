package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"log"
	"net"
	"net/rpc"

	"github.com/sbinet/npyio"
	"github.com/spiral/goridge"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
	"gonum.org/v1/gonum/mat"
)

type GoUr struct{}

func (s *GoUr) Infer(payload string, r *string) error {
	*r = infer(payload)
	return nil
}

func (s *GoUr) InferNumpy(payload []byte, ret *string) error {
	buf := bytes.NewBuffer(payload)
	r, err := npyio.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("npy-header: %v\n", r.Header)
	shape := r.Header.Descr.Shape
	raw := make([]float64, shape[0]*shape[1])

	err = r.Read(&raw)
	if err != nil {
		log.Fatal(err)
	}

	m := mat.NewDense(shape[0], shape[1], raw)
	//fmt.Printf("data = %v\n", mat.Formatted(m))
	scores := gour.GetScoresFromVectorized(ai, m)
	//fmt.Printf("scores = %v\n", scores)
	a, _ := json.Marshal(scores)
	*ret = string(a)
	return nil
}

func main() {
	log.Println("Loading AI")
	var err error
	ai, err = gour.LoadUrAI("trained/UR_evolving/2/ur_winner_genome_98-349")
	if err != nil {
		panic(err)
	}
	ln, err := net.Listen("tcp", "localhost:6001")
	if err != nil {
		panic(err)
	}

	rpc.Register(new(GoUr))

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeCodec(goridge.NewCodec(conn))
	}
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

func infer(board_json string) string {
	log.Println("Hello from infer")
	var board_input board_contract
	err := json.Unmarshal([]byte(board_json), &board_input)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return ""
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
		return ""
	}
	return string(inference)
}
