package main

import (
	"fmt"
	gour "gour/internal"

	"github.com/yaricom/goNEAT/v2/examples/utils"
	experiment2 "github.com/yaricom/goNEAT/v2/experiment"
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

	// The Ur runs
	experiment := experiment2.Experiment{}
	err = experiment.Execute(opts.NeatContext(), startGenome, gour.NewUrBootstrapGenerationEvaluator(outDirPath), nil)
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
	fmt.Printf("Mean best organisms: complexity=%.1f, diversity=%.1f, age=%.1f\n", meanComplexity, meanDiversity, meanAge)
}
