package gour

// Package xor defines the XOR experiment which serves to actually check that network topology actually evolves and
// everything works as expected.
// Because XOR is not linearly separable, a neural network requires hidden units to solve it. The two inputs must be
// combined at some hidden unit, as opposed to only at the output node, because there is no function over a linear
// combination of the inputs that can separate the inputs into the proper classes. These structural requirements make
// XOR suitable for testing NEAT’s ability to evolve structure.

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/experiment/utils"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

// The fitness threshold value for successful solver
const bootstrapFitnessThreshold = 31

type urBootstrapGenerationEvaluator struct {
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
func NewUrBootstrapGenerationEvaluator(outputPath string) experiment.GenerationEvaluator {
	return &urBootstrapGenerationEvaluator{OutputPath: outputPath}
}

// GenerationEvaluate This method evaluates one epoch for given population and prints results into output directory if any.
func (e *urBootstrapGenerationEvaluator) GenerationEvaluate(pop *genetics.Population, epoch *experiment.Generation, context *neat.Options) (err error) {
	// Evaluate each organism on a test
	tournament := EvaluateDoubleEliminationTournament(pop.Organisms, 7)
	tournament = EvaluateDoubleEliminationTournament(pop.Organisms, 5)
	tournament = EvaluateDoubleEliminationTournament(pop.Organisms, 3)
	max_tournament_wins := int(math.Sqrt(float64(tournament.Contender_Amount))) + 1
	fmt.Printf("Max fitness: %d\n", max_tournament_wins*3)
	fmt.Printf("Expected fitness: %d\n", max_tournament_wins*2)
	best := tournament.Contenders[len(tournament.Contenders)-1]
	best.SetWinner(true)
	if best.GetWins() >= max_tournament_wins*2 {
		epoch.Solved = true
	}
	if best.GetType() == "NEAT" {
		best := best.(*Ai_ur_player)
		epoch.WinnerNodes = len(best.Ai.Genotype.Nodes)
		epoch.WinnerGenes = best.Ai.Genotype.Extrons()
		epoch.WinnerEvals = context.PopSize*epoch.Id + best.Ai.Genotype.Id
		epoch.Best = best.Ai
		neat.InfoLog(fmt.Sprintf("Number of species: %v", len(pop.Species)))
		neat.InfoLog(fmt.Sprintf("Best fitness: %v", best.Ai.Fitness))
	}

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

type Potential_future struct {
	Pawn  int     `json:"pawn"`
	Score float64 `json:"score"`
}

func GetPotentialFutureScore(organism *genetics.Organism, current_board current_board_descriptor, potential_board potential_board_descriptor) (float64, error) {
	netDepth, err := organism.Phenotype.MaxActivationDepthFast(0) // The max depth of the network to be activated
	if err != nil {
		neat.WarnLog(fmt.Sprintf(
			"Failed to estimate maximal depth of the network with loop:\n%s\nUsing default depth: %d",
			organism.Genotype, netDepth))
	}
	//neat.DebugLog(fmt.Sprintf("Network depth: %d for organism: %d\n", netDepth, organism.Genotype.Id))
	if netDepth == 0 {
		neat.DebugLog(fmt.Sprintf("ALERT: Network depth is ZERO for Genome: %s", organism.Genotype))
		return 0, nil
	}

	success := false // Check for successful activation

	features_transformed := make([]float64, 0)

	for _, v := range current_board.board_state { // 2 to 21
		features_transformed = append(features_transformed, float64(v))
	}
	features_transformed = append(features_transformed, current_board.my_pawn_in_play)    // 22
	features_transformed = append(features_transformed, current_board.my_pawn_queue)      // 23
	features_transformed = append(features_transformed, current_board.my_pawn_out)        // 24
	features_transformed = append(features_transformed, current_board.enemy_pawn_in_play) // 25
	features_transformed = append(features_transformed, current_board.enemy_pawn_queue)   // 26
	features_transformed = append(features_transformed, current_board.enemy_pawn_out)     // 27

	for _, v := range potential_board.board_state { // 28 to 47
		features_transformed = append(features_transformed, float64(v))
	}
	features_transformed = append(features_transformed, potential_board.my_pawn_in_play)    // 48
	features_transformed = append(features_transformed, potential_board.my_pawn_queue)      // 49
	features_transformed = append(features_transformed, potential_board.my_pawn_out)        // 50
	features_transformed = append(features_transformed, potential_board.enemy_pawn_in_play) // 51
	features_transformed = append(features_transformed, potential_board.enemy_pawn_queue)   // 52
	features_transformed = append(features_transformed, potential_board.enemy_pawn_out)     // 53
	features_transformed = append(features_transformed, float64(potential_board.winner))    // 54
	features_transformed = append(features_transformed, float64(potential_board.turn))      // 55

	if err = organism.Phenotype.LoadSensors(features_transformed); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to load sensors: %s", err))
		return 0, err
	}

	// Use depth to ensure full relaxation
	if success, err = organism.Phenotype.ForwardSteps(netDepth); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to activate network: %s", err))
		return 0, err
	}

	out := organism.Phenotype.Outputs[0].Activation

	// Flush network for subsequent use
	if _, err = organism.Phenotype.Flush(); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to flush network: %s", err))
		return 0, err
	}

	if !success {
		return 0, errors.New("Not success")
	}

	return out, nil
}

func GetMoveScoresOrdered(board *board, organism *genetics.Organism) []*Potential_future {
	current_board_descriptor := GetCurrentBoardDescriptor(board, Left)
	potential_futures := []*Potential_future{}
	for pawn := range board.Current_player_path_moves {
		potential_game := board.Copy()
		potential_game.Play(pawn)
		//fmt.Println(potential_game.String())
		potential_board := GetPotentialBoardDescriptor(potential_game, board.Current_player)
		score, err := GetPotentialFutureScore(organism, current_board_descriptor, potential_board)
		//fmt.Println(score)
		if err != nil {
			panic(err)
		}
		potential_futures = append(potential_futures, &Potential_future{
			Pawn:  pawn,
			Score: score,
		})
	}
	sort.Slice(potential_futures, func(i, j int) bool {
		return potential_futures[i].Score < potential_futures[j].Score
	})
	return potential_futures
}

func LoadUrAI(path string) (*genetics.Organism, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	genome, err := genetics.ReadGenome(file, 1)
	if err != nil {
		return nil, err
	}

	file.Close()

	ai, err := genetics.NewOrganism(0, genome, 0)
	if err != nil {
		return nil, err
	}

	return ai, nil
}
