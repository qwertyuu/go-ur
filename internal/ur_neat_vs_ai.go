package gour

// Package xor defines the XOR experiment which serves to actually check that network topology actually evolves and
// everything works as expected.
// Because XOR is not linearly separable, a neural network requires hidden units to solve it. The two inputs must be
// combined at some hidden unit, as opposed to only at the output node, because there is no function over a linear
// combination of the inputs that can separate the inputs into the proper classes. These structural requirements make
// XOR suitable for testing NEAT’s ability to evolve structure.

import (
	"fmt"
	"sort"

	"github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/experiment/utils"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

// The fitness threshold value for successful solver
const urVsAifitnessThreshold = 100

type urVsAiGenerationEvaluator struct {
	// The output path to store execution results
	OutputPath string
}

// NewUrGenerationEvaluator is to create new generations evaluator to be used for the XOR experiment execution.
// XOR is very simple and does not make a very interesting scientific experiment; however, it is a good way to
// check whether your system works.
// Make sure recurrency is disabled for the XOR test. If NEAT is able to add recurrent connections, it may solve XOR by
// memorizing the order of the training set. (Which is why you may even want to randomize order to be most safe) All
// documented experiments with XOR are without recurrent connections. Interestingly, XOR can be solved by a recurrent
// network with no hidden nodes.
//
// This method performs evolution on XOR for specified number of generations and output results into outDirPath
// It also returns number of nodes, genes, and evaluations performed per each run (context.NumRuns)
func NewUrVsAiGenerationEvaluator(outputPath string) experiment.GenerationEvaluator {
	return &urVsAiGenerationEvaluator{OutputPath: outputPath}
}

// GenerationEvaluate This method evaluates one epoch for given population and prints results into output directory if any.
func (e *urVsAiGenerationEvaluator) GenerationEvaluate(pop *genetics.Population, epoch *experiment.Generation, context *neat.Options) (err error) {
	// Evaluate each organism on a test
	reference_ais := []string{
		"trained\\541\\ur_winner_genome_58-39",
	}
	// TODO: add number of moves as fitness (less moves, better fitness)
	for _, _ = range reference_ais {
		//reference_ai, _ := LoadUrAI(reference_ai_path)
		//reference := Ai_ur_player{
		//	Ai:   reference_ai,
		//	Name: "reference",
		//}
		reference := Random_ur_player{
			Name: "Random",
		}

		for i := 0; i < len(pop.Organisms); i++ {
			organism := Ai_ur_player{
				Ai:   pop.Organisms[i],
				Name: "organism",
			}
			OneVSOne(&organism, &reference, 3, 33)
			OneVSOne(&organism, &reference, 5, 33)
			OneVSOne(&organism, &reference, 7, 34)
		}
	}

	sort.Slice(pop.Organisms, func(i, j int) bool {
		return pop.Organisms[i].Fitness < pop.Organisms[j].Fitness
	})
	best := pop.Organisms[len(pop.Organisms)-1]
	best.IsWinner = true
	if best.Fitness >= float64(urVsAifitnessThreshold*len(reference_ais)) {
		epoch.Solved = true
	}
	epoch.WinnerNodes = len(best.Genotype.Nodes)
	epoch.WinnerGenes = best.Genotype.Extrons()
	epoch.WinnerEvals = context.PopSize*epoch.Id + best.Genotype.Id
	epoch.Best = best

	neat.InfoLog(fmt.Sprintf("Best fitness: %v", best.Fitness))

	// Fill statistics about current epoch
	epoch.FillPopulationStatistics(pop)

	// Only print to file every print_every generation
	if epoch.Solved || epoch.Id%context.PrintEvery == 0 {
		if _, err = utils.WritePopulationPlain(e.OutputPath, pop, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
			return err
		}
	}

	if epoch.Solved {
		// print winner organism
		org := epoch.Best
		if depth, err := org.Phenotype.MaxActivationDepthFast(0); err == nil {
			neat.InfoLog(fmt.Sprintf("Activation depth of the winner: %d\n", depth))
		}

		genomeFile := "ur_winner_genome"
		// Prints the winner organism's Genome to the file!
		if orgPath, err := utils.WriteGenomePlain(genomeFile, e.OutputPath, org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's genome, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's genome dumped to: %s\n", epoch.Id, orgPath))
		}

		// Prints the winner organism's Phenotype to the DOT file!
		if orgPath, err := utils.WriteGenomeDOT(genomeFile, e.OutputPath, org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's phenome DOT graph, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's phenome DOT graph dumped to: %s\n",
				epoch.Id, orgPath))
		}

		// Prints the winner organism's Phenotype to the Cytoscape JSON file!
		if orgPath, err := utils.WriteGenomeCytoscapeJSON(genomeFile, e.OutputPath, org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's phenome Cytoscape JSON graph, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's phenome Cytoscape JSON graph dumped to: %s\n",
				epoch.Id, orgPath))
		}
	}

	return err
}
