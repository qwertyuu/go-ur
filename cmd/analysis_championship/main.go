package main

import (
	"encoding/json"
	"fmt"
	gour "gour/internal"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type ai_score struct {
	Path       string
	Score      int
	Proportion float64
}

func main() {
	contenders := []*genetics.Organism{}
	organism_map := make(map[*genetics.Organism]string)
	win_counts := make(map[string]int)
	paths_to_scan := []string{}
	out_paths, _ := get_genome_dirs_from_dir("out/UR_oom_beat_28")
	paths_to_scan = append(paths_to_scan, out_paths...)
	out_paths, _ = get_genome_dirs_from_dir("out/UR_lost_current_beat_28")
	paths_to_scan = append(paths_to_scan, out_paths...)
	out_paths, _ = get_genome_dirs_from_dir("out/UR_beat_29")
	paths_to_scan = append(paths_to_scan, out_paths...)
	out_paths, _ = get_genome_dirs_from_dir("trained/541")
	paths_to_scan = append(paths_to_scan, out_paths...)

	for _, path := range paths_to_scan {
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
	rand.Seed(1)
	for i := 0; i < total_tournaments; i++ {
		tournament := gour.EvaluateDoubleEliminationTournament(contenders, 7)
		best := tournament.Contenders[len(tournament.Contenders)-1]
		if best.GetType() == "NEAT" {
			best := best.(*gour.Ai_ur_player)
			organism_path := organism_map[best.Ai]
			//fmt.Println(best.Ai.Fitness)
			//fmt.Println(organism_map[best.Ai])
			_, ok := win_counts[organism_path]
			if ok {
				win_counts[organism_path]++
			} else {
				win_counts[organism_path] = 1
			}
		}
		if i%100 == 0 && i > 0 {
			json_wins := get_json_wins(win_counts, i)
			fmt.Println(json_wins)
			fmt.Printf("%d/%d\n", i, total_tournaments)
		}
		for _, contender := range contenders {
			contender.Fitness = 0
			contender.Error = 0
		}
	}
	json_wins := get_json_wins(win_counts, total_tournaments)
	fmt.Println(json_wins)
	_ = ioutil.WriteFile("wins.json", []byte(json_wins), 0644)
}

func get_json_wins(win_counts map[string]int, total int) string {
	ai_scores := make([]ai_score, 0, len(win_counts))
	for path, score := range win_counts {
		ai_scores = append(ai_scores, ai_score{
			Path:       path,
			Score:      score,
			Proportion: float64(score) / float64(total),
		})
	}
	sort.Slice(ai_scores, func(i, j int) bool {
		return ai_scores[i].Score > ai_scores[j].Score
	})
	json_wins, _ := json.Marshal(ai_scores)
	return string(json_wins)
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

func get_genome_dirs_from_dir(dir string) ([]string, error) {
	folders := []string{}
	err := filepath.WalkDir(dir, func(subpath string, f fs.DirEntry, err error) error {
		if strings.HasPrefix(filepath.Base(subpath), "ur_winner_genome") && filepath.Ext(subpath) == "" {
			folders = append(folders, subpath)
		}
		return nil
	})

	return folders, err
}
