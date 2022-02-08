package main

import (
	"fmt"
	gour "gour/internal"
	"os"
	"path/filepath"
	"strings"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

var ai *genetics.Organism

func main() {
	contenders := []*genetics.Organism{}
	organism_map := make(map[*genetics.Organism]string)
	win_counts := make(map[string]int)
	i := 0
	for len(contenders) < 64 {
		path := fmt.Sprintf("out/UR/%d", i)
		i++
		genome_path, err := get_genome_from_dir(path)
		if genome_path == "" {
			continue
		}
		fmt.Println(genome_path)
		if err != nil {
			panic(err)
		}
		ai, err := gour.LoadUrAI(genome_path)
		if err != nil {
			panic(err)
		}
		organism_map[ai] = genome_path
		contenders = append(contenders, ai)
	}
	for i := 0; i < 100; i++ {
		tournament := gour.EvaluateDoubleEliminationTournament(contenders)
		best := tournament.Contenders[len(tournament.Contenders)-1]
		organism_path := organism_map[best]
		_, ok := win_counts[organism_path]
		if ok {
			win_counts[organism_path]++
		} else {
			win_counts[organism_path] = 1
		}

		fmt.Println(best.Fitness)
		fmt.Println(organism_map[best])
		for _, contender := range contenders {
			contender.Fitness = 0
			contender.Error = 0
		}
	}

	fmt.Println(win_counts)
}

func get_genome_from_dir(dir string) (string, error) {
	var file string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), "ur_winner_genome") && filepath.Ext(path) == "" {
			file = path
		}
		return nil
	})

	return file, err
}
