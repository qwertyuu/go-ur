package gour

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaricom/goNEAT/v2/examples/utils"
	experiment2 "github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

func TestBoard(t *testing.T) {
	pcg32.Seed(1, 1)
	b := NewBoard(7)

	count := 0
	for b.Current_winner == 0 {
		for k := range b.Current_player_path_moves {
			b.Play(k)
			//t.Log(b.String())
			t.Log(b.AsArray(Right))
			//print(b.String())
			//print(b.Current_dice)
			//print("\n")
			//t.Log(b.Current_player)
			//t.Log(b.Current_player_path_moves)
			//t.Log(b.right_player_out)
			//t.Log(b.left_player_out)
			break
		}
		count++
		if count == 10 {
			break
		}
	}

	j := b.Copy()
	t.Log(j.Current_player_path_moves)
	t.Log(b.AsArray(Right))
	t.Log(j.AsArray(Right))
	j.Play(0)
	t.Log(b.AsArray(Right))
	t.Log(j.AsArray(Right))

}

// The integration test running over multiple iterations in order to detect if any random errors occur.
func TestNeat(t *testing.T) {
	rand.Seed(1)
	pcg32.Seed(1, 1)

	outDirPath, contextPath, genomePath := "out/UR", "data/ur.neat", "data/urstartgenes"

	// Load Genome
	fmt.Println("Loading start genome for Ur experiment")
	opts, startGenome, err := LoadOptionsAndGenome(contextPath, genomePath)
	neat.LogLevel = neat.LogLevelInfo
	require.NoError(t, err)

	// Check if output dir exists
	err = utils.CreateOutputDir(outDirPath)
	require.NoError(t, err, "Failed to create output directory")

	// The 100 runs XOR experiment
	opts.NumRuns = 100
	experiment := experiment2.Experiment{
		Id:     0,
		Trials: make(experiment2.Trials, opts.NumRuns),
	}
	err = experiment.Execute(opts.NeatContext(), startGenome, NewUrGenerationEvaluator(outDirPath), nil)
	require.NoError(t, err, "Failed to perform XOR experiment")

	// Find winner statistics
	avgNodes, avgGenes, avgEvals, _ := experiment.AvgWinner()

	maxEvals := float64(opts.PopSize * opts.NumGenerations)
	assert.True(t, avgEvals < maxEvals)

	t.Logf("avg_nodes: %.1f, avg_genes: %.1f, avg_evals: %.1f\n", avgNodes, avgGenes, avgEvals)
	meanComplexity, meanDiversity, meanAge := 0.0, 0.0, 0.0
	for _, t := range experiment.Trials {
		meanComplexity += t.BestComplexity().Mean()
		meanDiversity += t.Diversity().Mean()
		meanAge += t.BestAge().Mean()
	}
	count := float64(len(experiment.Trials))
	meanComplexity /= count
	meanDiversity /= count
	meanAge /= count
	t.Logf("Mean best organisms: complexity=%.1f, diversity=%.1f, age=%.1f", meanComplexity, meanDiversity, meanAge)
}


func LoadOptionsAndGenome(contextPath, genomePath string) (*neat.Options, *genetics.Genome, error) {
	// Load context configuration
	configFile, err := os.Open(contextPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open context file")
	}
	context, err := neat.LoadNeatOptions(configFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load NEAT options")
	}

	// Load start Genome
	genomeFile, err := os.Open(genomePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open genome file")
	}
	r, err := genetics.NewGenomeReader(genomeFile, genetics.YAMLGenomeEncoding)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load NEAT genome")
	}
	genome, err := r.Read()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load NEAT genome")
	}
	return context, genome, nil
}