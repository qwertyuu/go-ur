package gour

import (
	"errors"
	"fmt"
	"sort"

	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type double_elimination struct {
	Contenders               []*genetics.Organism
	Losers_of_winner_bracket []*genetics.Organism
	Loser_bracket            []*genetics.Organism
	Winner_bracket           []*genetics.Organism
	Champion                 *genetics.Organism
}

type potential_future struct {
	pawn  int
	score float64
}

// NB de joutes dans la winner bracket avant la loser's bracket = len(contenders) / 2 - 1
// Commencer par évaluer la winner bracket au complet. Noter les perdants de chaque joute
// L'ordre de correspondance entre les L et les losers est inversée. Les premiers de la losers bracket affrontent les derniers de la winners bracket
// Les affrontements entre les losers et le L se fait à toutes les deux séries de joutes jusqu'à détermination d'un gagnant, qui va se battre contre le gagnant de la winner's bracket

func EvaluateDoubleEliminationTournament(contenders []*genetics.Organism) double_elimination {
	if !IsPowerOfTwo(len(contenders)) {
		panic("oups")
	}
	tournament := double_elimination{
		Contenders:               contenders,
		Losers_of_winner_bracket: make([]*genetics.Organism, 0),
		Loser_bracket:            make([]*genetics.Organism, 0),
		Winner_bracket:           make([]*genetics.Organism, 0),
	}

	// determine loser and winner brackets
	var left_player *genetics.Organism = nil
	for _, right_player := range tournament.Contenders {
		right_player.IsWinner = false
		if left_player == nil {
			left_player = right_player
			continue
		}
		winner_position := Fight(left_player, right_player)
		if winner_position == Left {
			tournament.Winner_bracket = append(tournament.Winner_bracket, left_player)
			tournament.Loser_bracket = append(tournament.Loser_bracket, right_player)
		} else {
			tournament.Winner_bracket = append(tournament.Winner_bracket, right_player)
			tournament.Loser_bracket = append(tournament.Loser_bracket, left_player)
		}
		left_player = nil
	}

	// evaluate winner bracket
	for len(tournament.Winner_bracket) > 1 {
		next_winner_bracket := []*genetics.Organism{}
		for _, right_player := range tournament.Winner_bracket {
			if left_player == nil {
				left_player = right_player
				continue
			}
			winner_position := Fight(left_player, right_player)
			if winner_position == Left {
				next_winner_bracket = append(next_winner_bracket, left_player)
				tournament.Losers_of_winner_bracket = append(tournament.Losers_of_winner_bracket, right_player)
			} else {
				next_winner_bracket = append(next_winner_bracket, right_player)
				tournament.Losers_of_winner_bracket = append(tournament.Losers_of_winner_bracket, left_player)
			}
			left_player = nil
		}
		tournament.Winner_bracket = next_winner_bracket
	}

	// inverse loser bracket contestant to mix things up
	for i, j := 0, len(tournament.Loser_bracket)-1; i < j; i, j = i+1, j-1 {
		tournament.Loser_bracket[i], tournament.Loser_bracket[j] = tournament.Loser_bracket[j], tournament.Loser_bracket[i]
	}
	winner_bracket_loser_pointer := 0
	// evaluate loser bracket
	for len(tournament.Loser_bracket) > 1 {
		next_loser_bracket := []*genetics.Organism{}
		for _, right_player := range tournament.Loser_bracket {
			if left_player == nil {
				left_player = right_player
				continue
			}
			winner_position := Fight(left_player, right_player)
			loser_of_winner_bracket := tournament.Losers_of_winner_bracket[winner_bracket_loser_pointer]
			winner_bracket_loser_pointer++
			if winner_position == Left {
				winner_of_second_fight := Fight(left_player, loser_of_winner_bracket)
				if winner_of_second_fight == Left {
					next_loser_bracket = append(next_loser_bracket, left_player)
				} else {
					next_loser_bracket = append(next_loser_bracket, loser_of_winner_bracket)
				}
			} else {
				winner_of_second_fight := Fight(right_player, loser_of_winner_bracket)
				if winner_of_second_fight == Left {
					next_loser_bracket = append(next_loser_bracket, right_player)
				} else {
					next_loser_bracket = append(next_loser_bracket, loser_of_winner_bracket)
				}
			}
			left_player = nil
		}
		tournament.Loser_bracket = next_loser_bracket
	}

	// final epic fight
	final_winner := Fight(tournament.Winner_bracket[0], tournament.Loser_bracket[0])
	if final_winner == Left {
		tournament.Champion = tournament.Winner_bracket[0]
	} else {
		tournament.Champion = tournament.Loser_bracket[0]
	}

	sort.Slice(tournament.Contenders, func(i, j int) bool {
		return tournament.Contenders[i].Fitness < tournament.Contenders[j].Fitness
	})

	for _, o := range tournament.Contenders {
		o.Error = 1 / o.Fitness
	}

	return tournament
}

func Fight(left_player *genetics.Organism, right_player *genetics.Organism) int {
	game := NewBoard(7)
	for game.Current_winner == 0 {
		if len(game.Current_player_path_moves) == 1 {
			for pawn := range game.Current_player_path_moves {
				game.Play(pawn)
			}
			continue
		} else {
			potential_futures := []*potential_future{}
			current_board_left := GetCurrentBoardDescriptor(game, Left)
			current_board_right := GetCurrentBoardDescriptor(game, Right)
			for pawn := range game.Current_player_path_moves {
				potential_game := game.Copy()
				potential_game.Play(pawn)
				potential_board := GetPotentialBoardDescriptor(potential_game, game.Current_player)
				if game.Current_player == Left {
					score, err := GetPotentialFutureScore(left_player, current_board_left, potential_board)
					if err != nil {
						panic(err)
					}
					potential_futures = append(potential_futures, &potential_future{
						pawn:  pawn,
						score: score,
					})
				} else {
					score, err := GetPotentialFutureScore(right_player, current_board_right, potential_board)
					if err != nil {
						panic(err)
					}
					potential_futures = append(potential_futures, &potential_future{
						pawn:  pawn,
						score: score,
					})
				}
			}
			sort.Slice(potential_futures, func(i, j int) bool {
				return potential_futures[i].score < potential_futures[j].score
			})

			// Play highest-scoring future
			game.Play(potential_futures[len(potential_futures)-1].pawn)
		}
	}
	if game.Current_winner == Left {
		left_player.Fitness++
	} else {
		right_player.Fitness++
	}
	return game.Current_winner
}

func GetPotentialFutureScore(organism *genetics.Organism, current_board current_board_descriptor, potential_board potential_board_descriptor) (float64, error) {
	// TODO: Implement from ur_neat.go orgEvaluate()
	netDepth, err := organism.Phenotype.MaxActivationDepthFast(0) // The max depth of the network to be activated
	if err != nil {
		neat.WarnLog(fmt.Sprintf(
			"Failed to estimate maximal depth of the network with loop:\n%s\nUsing default depth: %d",
			organism.Genotype, netDepth))
	}
	neat.DebugLog(fmt.Sprintf("Network depth: %d for organism: %d\n", netDepth, organism.Genotype.Id))
	if netDepth == 0 {
		neat.DebugLog(fmt.Sprintf("ALERT: Network depth is ZERO for Genome: %s", organism.Genotype))
		return 0, nil
	}

	success := false // Check for successful activation

	features_transformed := make([]float64, 0)

	for _, v := range current_board.board_state { // 0 to 19
		features_transformed = append(features_transformed, float64(v))
	}
	features_transformed = append(features_transformed, current_board.my_pawn_in_play) // 20
	features_transformed = append(features_transformed, current_board.my_pawn_queue) // 21
	features_transformed = append(features_transformed, current_board.my_pawn_out) // 22
	features_transformed = append(features_transformed, current_board.enemy_pawn_in_play) // 23
	features_transformed = append(features_transformed, current_board.enemy_pawn_queue) // 24
	features_transformed = append(features_transformed, current_board.enemy_pawn_out) // 25

	for _, v := range potential_board.board_state { // 26 to 45
		features_transformed = append(features_transformed, float64(v))
	}
	features_transformed = append(features_transformed, potential_board.my_pawn_in_play) // 46
	features_transformed = append(features_transformed, potential_board.my_pawn_queue) // 47
	features_transformed = append(features_transformed, potential_board.my_pawn_out) // 48
	features_transformed = append(features_transformed, potential_board.enemy_pawn_in_play) // 49
	features_transformed = append(features_transformed, potential_board.enemy_pawn_queue) // 50
	features_transformed = append(features_transformed, potential_board.enemy_pawn_out) // 51
	features_transformed = append(features_transformed, float64(potential_board.winner)) // 52
	features_transformed = append(features_transformed, float64(potential_board.turn)) // 53

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

func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}
