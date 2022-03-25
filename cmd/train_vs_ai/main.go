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

	opts.NumRuns = 900

	// The Ur runs
	experiment := experiment2.Experiment{}
	err = experiment.Execute(opts.NeatContext(), startGenome, gour.NewUrVsAiGenerationEvaluator(outDirPath, 30, false), nil)
	if err != nil {
		panic(err)
	}
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
