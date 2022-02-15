package main

import (
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type ai_score struct {
	Path  string
	Score int
}

func main() {
	contenders := []*genetics.Organism{}
	organism_map := make(map[*genetics.Organism]string)
	win_counts := make(map[string]int)
	i := 0
	for {
		path := fmt.Sprintf("out/UR_reference_ai/%d", i)
		i++
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}
		genome_path, err := get_genome_from_dir(path)
		if err != nil {
			panic(err)
		}
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
	total_tournaments := 1000
	for i := 0; i < total_tournaments; i++ {
		tournament := gour.EvaluateDoubleEliminationTournament(contenders, 7)
		best := tournament.Contenders[len(tournament.Contenders)-1]
		if best.GetType() == "NEAT" {
			best := best.(*gour.Ai_ur_player)
			organism_path := organism_map[best.Ai]
			fmt.Println(best.Ai.Fitness)
			fmt.Println(organism_map[best.Ai])
			_, ok := win_counts[organism_path]
			if ok {
				win_counts[organism_path]++
			} else {
				win_counts[organism_path] = 1
			}
		}

		fmt.Printf("%d/%d\n", i, total_tournaments)
		for _, contender := range contenders {
			contender.Fitness = 0
			contender.Error = 0
		}
	}

	ai_scores := make([]ai_score, 0, len(win_counts))

	for path, score := range win_counts {
		ai_scores = append(ai_scores, ai_score{
			Path:  path,
			Score: score,
		})
	}
	sort.Slice(ai_scores, func(i, j int) bool {
		return ai_scores[i].Score > ai_scores[j].Score
	})
	json_wins, _ := json.Marshal(ai_scores)
	fmt.Println(string(json_wins))
	_ = ioutil.WriteFile("wins.json", json_wins, 0644)
}

func get_genome_from_dir(dir string) (string, error) {
	var file string
	err := filepath.WalkDir(dir, func(path string, f fs.DirEntry, err error) error {
		if strings.HasPrefix(filepath.Base(path), "ur_winner_genome") && filepath.Ext(path) == "" {
			file = path
		}
		return nil
	})

	return file, err
}
