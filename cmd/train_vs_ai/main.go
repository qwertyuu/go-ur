package main

import (
	"fmt"
	gour "gour/internal"

	"github.com/yaricom/goNEAT/v2/examples/utils"
	"github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/neat"
)

func main() {
	outDirPath, contextPath, genomePath := "out/UR", "data/ur.neat", "data/urstartgenes.yml"

	// Load Genome
	fmt.Println("Loading start genome for Ur experiment")
	opts, startGenome, err := gour.LoadOptionsAndGenome(contextPath, genomePath)
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
	experiment := experiment.Experiment{}
	err = experiment.Execute(opts.NeatContext(), startGenome, gour.NewUrVsAiGenerationEvaluator(outDirPath, 30, false), nil)

	if err != nil {
		panic(err)
	}
}
