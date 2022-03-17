package main

import (
	"fmt"
	gour "gour/internal"
	"os"

	"github.com/pkg/errors"
	"github.com/yaricom/goNEAT/v2/examples/utils"
	experiment2 "github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

func main() {
	outDirPath, contextPath, genomePath := "out/UR", "data/ur.neat", "data/urstartgenes.yml"

	// Load Genome
	fmt.Println("Loading start genome for Ur experiment")
	opts, startGenome, err := LoadOptionsAndGenome(contextPath, genomePath)
	if err != nil {
		panic(err)
	}
	neat.LogLevel = neat.LogLevelInfo

	// Check if output dir exists
	err = utils.CreateOutputDir(outDirPath)
	if err != nil {
		panic(err)
	}

	// The Ur runs
	experiment := experiment2.Experiment{}
	err = experiment.Execute(opts.NeatContext(), startGenome, gour.NewUrVsAiGenerationEvaluator(outDirPath), nil)
	if err != nil {
		panic(err)
	}

	// Find winner statistics
	avgNodes, avgGenes, avgEvals, _ := experiment.AvgWinner()

	fmt.Printf("avg_nodes: %.1f, avg_genes: %.1f, avg_evals: %.1f\n", avgNodes, avgGenes, avgEvals)
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
	fmt.Printf("Mean best organisms: complexity=%.1f, diversity=%.1f, age=%.1f", meanComplexity, meanDiversity, meanAge)
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
