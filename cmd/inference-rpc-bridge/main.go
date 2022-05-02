package main

import (
	"bytes"
	"encoding/json"
	gour "gour/internal"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/sbinet/npyio"
	"github.com/spiral/goridge"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
	"gonum.org/v1/gonum/mat"
)

type GoUr struct{}

func (s *GoUr) InferNumpy(payload []byte, ret *string) error {
	start := time.Now()
	buf := bytes.NewBuffer(payload)
	r, err := npyio.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}

	shape := r.Header.Descr.Shape
	raw := make([]float64, shape[0]*shape[1])

	err = r.Read(&raw)
	if err != nil {
		log.Fatal(err)
	}

	m := mat.NewDense(shape[0], shape[1], raw)
	scores := gour.GetScoresFromVectorized(ai, m)
	scores_json, err := json.Marshal(scores)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	log.Printf("Returning inference for %v in %s", len(scores), elapsed)
	*ret = string(scores_json)
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

var ai *genetics.Organism
