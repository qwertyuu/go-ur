package gour

import (
	"fmt"
	"sort"

	"github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/experiment/utils"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type UrVsAiGenerationEvaluator struct {
	// The output path to store execution results
	OutputPath string

	// the number of games to play, per game size (so will be x3)
	NumberOfGames int

	// if we are in an evolving context or not
	Evolve bool
}

func NewUrVsAiGenerationEvaluator(outputPath string, numberOfGames int, evolve bool) *UrVsAiGenerationEvaluator {
	return &UrVsAiGenerationEvaluator{
		OutputPath:    outputPath,
		NumberOfGames: numberOfGames,
		Evolve:        evolve,
	}
}

func (e *UrVsAiGenerationEvaluator) GetOutputPath() string {
	return e.OutputPath
}

// GenerationEvaluate This method evaluates one epoch for given population and prints results into output directory if any.
func (e *UrVsAiGenerationEvaluator) GenerationEvaluate(pop *genetics.Population, epoch *experiment.Generation, context *neat.Options) (err error) {
	// Evaluate each organism on a test
	//ai_1, err := LoadUrAI("out\\UR\\0\\ur_winner_genome_63-236")
	fmt.Printf("Training against AI for %v games\n", e.NumberOfGames)
	reference_ais := []Ur_player{
		//&Ai_ur_player{
		//	Ai:   ai_1,
		//	Name: "reference",
		//},
		&Random_ur_player{
			Name: "Random",
		},
	}
	bestFitness := float64(e.NumberOfGames * 3 * len(reference_ais))
	fmt.Printf("Target fitness: %v\n", bestFitness)
	// TODO: add number of moves as fitness (less moves, better fitness)
	for _, reference_ai := range reference_ais {
		for i := 0; i < len(pop.Organisms); i++ {
			organism := &Ai_ur_player{
				Ai:   pop.Organisms[i],
				Name: "organism",
			}
			OneVSOne(organism, reference_ai, 3, e.NumberOfGames)
			OneVSOne(organism, reference_ai, 5, e.NumberOfGames)
			OneVSOne(organism, reference_ai, 7, e.NumberOfGames)
		}
	}

	sort.Slice(pop.Organisms, func(i, j int) bool {
		return pop.Organisms[i].Fitness < pop.Organisms[j].Fitness
	})
	best := pop.Organisms[len(pop.Organisms)-1]
	best.IsWinner = true
	if best.Fitness >= float64(e.NumberOfGames*3*len(reference_ais)) {
		epoch.Solved = true
	}
	epoch.WinnerNodes = len(best.Genotype.Nodes)
	epoch.WinnerGenes = best.Genotype.Extrons()
	epoch.WinnerEvals = context.PopSize*epoch.Id + best.Genotype.Id
	epoch.Best = best

	neat.InfoLog(fmt.Sprintf("Number of species: %v", len(pop.Species)))
	neat.InfoLog(fmt.Sprintf("Best fitness: %v", best.Fitness))

	// Fill statistics about current epoch
	epoch.FillPopulationStatistics(pop)

	// Only print to file every print_every generation
	if epoch.Solved || epoch.Id%context.PrintEvery == 0 {
		if _, err = utils.WritePopulationPlain(e.GetOutputPath(), pop, epoch); err != nil {
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
		if orgPath, err := utils.WriteGenomePlain(genomeFile, e.GetOutputPath(), org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's genome, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's genome dumped to: %s\n", epoch.Id, orgPath))
		}

		// Prints the winner organism's Phenotype to the DOT file!
		if orgPath, err := utils.WriteGenomeDOT(genomeFile, e.GetOutputPath(), org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's phenome DOT graph, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's phenome DOT graph dumped to: %s\n",
				epoch.Id, orgPath))
		}

		// Prints the winner organism's Phenotype to the Cytoscape JSON file!
		if orgPath, err := utils.WriteGenomeCytoscapeJSON(genomeFile, e.GetOutputPath(), org, epoch); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism's phenome Cytoscape JSON graph, reason: %s\n", err))
		} else {
			neat.InfoLog(fmt.Sprintf("Generation #%d winner's phenome Cytoscape JSON graph dumped to: %s\n",
				epoch.Id, orgPath))
		}
		if e.Evolve {
			e.NumberOfGames *= 2
			epoch.Solved = false
		}
	}

	return err
}
