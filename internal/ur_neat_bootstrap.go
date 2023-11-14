package gour

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

type urBootstrapGenerationEvaluator struct {
	// The output path to store execution results
	OutputPath string
}

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

type Pawn_score_compare struct {
	PawnA  int
	PawnB  int
	ScoreA float64
	ScoreB float64
}

func GetFeatureNames() []string {
	feature_names := make([]string, 54)
	ti := 0

	// The numbers are by NEAT id, which start at 1 and 1 is "reserved" as a BIAS node, so we start counting here at 2.

	for i := 0; i < 20; i++ { // 2 to 21
		feature_names[ti] = fmt.Sprintf("current_board_state_%v", i)
		ti++
	}

	feature_names[ti] = "current_board_my_pawn_in_play" // 22
	ti++
	feature_names[ti] = "current_board_my_pawn_queue" // 23
	ti++
	feature_names[ti] = "current_board_my_pawn_out" // 24
	ti++
	feature_names[ti] = "current_board_enemy_pawn_in_play" // 25
	ti++
	feature_names[ti] = "current_board_enemy_pawn_queue" // 26
	ti++
	feature_names[ti] = "current_board_enemy_pawn_out" // 27
	ti++

	for i := 0; i < 20; i++ { // 28 to 47
		feature_names[ti] = fmt.Sprintf("potential_board_state_%v", i)
		ti++
	}

	feature_names[ti] = "potential_board_my_pawn_in_play" // 48
	ti++
	feature_names[ti] = "potential_board_my_pawn_queue" // 49
	ti++
	feature_names[ti] = "potential_board_my_pawn_out" // 50
	ti++
	feature_names[ti] = "potential_board_enemy_pawn_in_play" // 51
	ti++
	feature_names[ti] = "potential_board_enemy_pawn_queue" // 52
	ti++
	feature_names[ti] = "potential_board_enemy_pawn_out" // 53
	ti++
	feature_names[ti] = "potential_board_winner" // 54
	ti++
	feature_names[ti] = "potential_board_turn" // 55

	return feature_names
}

func Vectorize(potential_move_a potential_board_descriptor, potential_move_b potential_board_descriptor) []float64 {
	features_transformed := make([]float64, 56)
	ti := 0

	// The numbers are by NEAT id, which start at 1 and 1 is "reserved" as a BIAS node, so we start counting here at 2.

	for _, v := range potential_move_a.board_state { // 2 to 21
		features_transformed[ti] = float64(v)
		ti++
	}

	features_transformed[ti] = potential_move_a.my_pawn_in_play // 22
	ti++
	features_transformed[ti] = potential_move_a.my_pawn_queue // 23
	ti++
	features_transformed[ti] = potential_move_a.my_pawn_out // 24
	ti++
	features_transformed[ti] = potential_move_a.enemy_pawn_in_play // 25
	ti++
	features_transformed[ti] = potential_move_a.enemy_pawn_queue // 26
	ti++
	features_transformed[ti] = potential_move_a.enemy_pawn_out // 27
	ti++
	features_transformed[ti] = float64(potential_move_a.winner) // 28
	ti++
	features_transformed[ti] = float64(potential_move_a.turn) // 29

	for _, v := range potential_move_b.board_state { // 30 to 49
		features_transformed[ti] = float64(v)
		ti++
	}

	features_transformed[ti] = potential_move_b.my_pawn_in_play // 50
	ti++
	features_transformed[ti] = potential_move_b.my_pawn_queue // 51
	ti++
	features_transformed[ti] = potential_move_b.my_pawn_out // 52
	ti++
	features_transformed[ti] = potential_move_b.enemy_pawn_in_play // 53
	ti++
	features_transformed[ti] = potential_move_b.enemy_pawn_queue // 54
	ti++
	features_transformed[ti] = potential_move_b.enemy_pawn_out // 55
	ti++
	features_transformed[ti] = float64(potential_move_b.winner) // 56
	ti++
	features_transformed[ti] = float64(potential_move_b.turn) // 57

	return features_transformed
}

func GetPotentialFutureScore(organism *genetics.Organism, features_transformed []float64) (float64, float64, error) {
	success := false // Check for successful activation

	netDepth, err := organism.Phenotype.MaxActivationDepthFast(0) // The max depth of the network to be activated
	if err != nil {
		neat.WarnLog(fmt.Sprintf(
			"Failed to estimate maximal depth of the network with loop:\n%s\nUsing default depth: %d",
			organism.Genotype, netDepth))
	}
	//neat.DebugLog(fmt.Sprintf("Network depth: %d for organism: %d\n", netDepth, organism.Genotype.Id))
	if netDepth == 0 {
		neat.DebugLog(fmt.Sprintf("ALERT: Network depth is ZERO for Genome: %s", organism.Genotype))
		return 0, 0, nil
	}

	if err = organism.Phenotype.LoadSensors(features_transformed); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to load sensors: %s", err))
		return 0, 0, err
	}

	// Use depth to ensure full relaxation
	if success, err = organism.Phenotype.ForwardSteps(netDepth); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to activate network: %s", err))
		return 0, 0, err
	}

	outA := organism.Phenotype.Outputs[0].Activation
	outB := organism.Phenotype.Outputs[1].Activation

	// Flush network for subsequent use
	if _, err = organism.Phenotype.Flush(); err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to flush network: %s", err))
		return 0, 0, err
	}

	if !success {
		return 0, 0, errors.New("not success")
	}

	return outA, outB, nil
}

//func GetScoresFromVectorized(organism *genetics.Organism, vectorized *mat.Dense) []float64 {
//	rows, _ := vectorized.Dims()
//	scores := make([]float64, rows)
//	for i := 0; i < rows; i++ {
//		row := vectorized.RawRowView(i)
//		score, err := GetPotentialFutureScore(organism, row)
//		if err != nil {
//			panic(err)
//		}
//		scores[i] = score
//	}
//	return scores
//}

func generateMovePairs(moves []int) [][]int {
	pairs := [][]int{}

	for i, moveA := range moves {
		for j := i + 1; j < len(moves); j++ {
			moveB := moves[j]
			pair := []int{moveA, moveB}
			pairs = append(pairs, pair)
		}
	}

	return pairs
}
func GetMoveScoresOrdered(board *board, organism *genetics.Organism) ([]*Potential_future, error) {
	// TODO: Add current board descriptor to the vectorize back, I think it is lacking
	//current_board := GetCurrentBoardDescriptor(board, Left)

	// list pawns from board.Current_player_path_moves
	moves := []int{}
	move_descriptor := make(map[int]potential_board_descriptor)
	for pawn := range *board.Current_player_path_moves {
		potential_game := board.Copy()
		potential_game.Play(pawn)
		potential_board := GetPotentialBoardDescriptor(potential_game, board.Current_player)
		move_descriptor[pawn] = potential_board
		moves = append(moves, pawn)
	}
	pairs := generateMovePairs(moves)

	pawn_score_compare := []Pawn_score_compare{}
	for _, pawnPair := range pairs {
		pawn_a := pawnPair[0]
		pawn_b := pawnPair[1]
		// TODO: Add current board descriptor to the vectorize back, I think it is lacking
		transformed_features := Vectorize(move_descriptor[pawn_a], move_descriptor[pawn_b])
		score_a, score_b, err := GetPotentialFutureScore(organism, transformed_features)
		//fmt.Println(score)
		if err != nil {
			return nil, err
			//panic(err)
		}
		pawn_score_compare = append(pawn_score_compare, Pawn_score_compare{
			PawnA:  pawn_a,
			PawnB:  pawn_b,
			ScoreA: score_a,
			ScoreB: score_b,
		})
	}

	// pick the best pawn out of the pairs by averaging the scores for each pawn and picking the highest
	pawn_score := make(map[int]float64)
	for _, v := range pawn_score_compare {
		pawn_score[v.PawnA] += v.ScoreA
		pawn_score[v.PawnB] += v.ScoreB
	}

	potential_futures := []*Potential_future{}
	for k, v := range pawn_score {
		potential_futures = append(potential_futures, &Potential_future{
			Pawn:  k,
			Score: v,
		})
	}

	sort.Slice(potential_futures, func(i, j int) bool {
		return potential_futures[i].Score < potential_futures[j].Score
	})
	return potential_futures, nil
}

//func GetMovesVectors(board *board) [][]float64 {
//	all_potential_board_trf := [][]float64{}
//	current_board_descriptor := GetCurrentBoardDescriptor(board, Left)
//
//	for pawn := range *board.Current_player_path_moves {
//		potential_game := board.Copy()
//		potential_game.Play(pawn)
//		potential_board := GetPotentialBoardDescriptor(potential_game, board.Current_player)
//		transformed_features := Vectorize(current_board_descriptor, potential_board)
//		all_potential_board_trf = append(all_potential_board_trf, transformed_features)
//	}
//	return all_potential_board_trf
//}

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
