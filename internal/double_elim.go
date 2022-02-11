package gour

import (
	"sort"

	"github.com/yaricom/goNEAT/v2/neat/genetics"
)

type double_elimination struct {
	Contenders               []*genetics.Organism
	Losers_of_winner_bracket []*genetics.Organism
	Loser_bracket            []*genetics.Organism
	Winner_bracket           []*genetics.Organism
	Champion                 *genetics.Organism
}

// NB de joutes dans la winner bracket avant la loser's bracket = len(contenders) / 2 - 1
// Commencer par évaluer la winner bracket au complet. Noter les perdants de chaque joute
// L'ordre de correspondance entre les L et les losers est inversée. Les premiers de la losers bracket affrontent les derniers de la winners bracket
// Les affrontements entre les losers et le L se fait à toutes les deux séries de joutes jusqu'à détermination d'un gagnant, qui va se battre contre le gagnant de la winner's bracket

func EvaluateDoubleEliminationTournament(contenders []*genetics.Organism) double_elimination {
	contender_power_of_two := getNearestPowerOfTwo(len(contenders))
	has_nil_contenders := contender_power_of_two != len(contenders)
	for contender_power_of_two > len(contenders) {
		contenders = append(contenders, nil)
	}
	tournament := double_elimination{
		Contenders:               contenders,
		Losers_of_winner_bracket: make([]*genetics.Organism, 0),
		Loser_bracket:            make([]*genetics.Organism, 0),
		Winner_bracket:           make([]*genetics.Organism, 0),
	}

	// determine loser and winner brackets
	var left_player *genetics.Organism = nil
	left_player_set := false
	for _, right_player := range tournament.Contenders {
		if right_player != nil {
			right_player.IsWinner = false
		}
		if !left_player_set {
			left_player = right_player
			left_player_set = true
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
		left_player_set = false
		left_player = nil
	}

	// evaluate winner bracket
	for len(tournament.Winner_bracket) > 1 {
		next_winner_bracket := []*genetics.Organism{}
		for _, right_player := range tournament.Winner_bracket {
			if !left_player_set {
				left_player = right_player
				left_player_set = true
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
			left_player_set = false
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
			if !left_player_set {
				left_player = right_player
				left_player_set = true
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
			left_player_set = false
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

	if has_nil_contenders {
		new_contenders := []*genetics.Organism{}
		for _, contender := range tournament.Contenders {
			if contender != nil {
				new_contenders = append(new_contenders, contender)
			}
		}
		tournament.Contenders = new_contenders
	}

	sort.Slice(tournament.Contenders, func(i, j int) bool {
		return tournament.Contenders[i].Fitness < tournament.Contenders[j].Fitness
	})

	for _, o := range tournament.Contenders {
		o.Error = 1 / o.Fitness
	}

	return tournament
}

func remove(s []*genetics.Organism, i int) []*genetics.Organism {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getNearestPowerOfTwo(i int) int {
	var v uint32 = uint32(i)
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return int(v)
}

func Fight(left_player *genetics.Organism, right_player *genetics.Organism) int {
	if left_player == nil {
		return Right
	} else if right_player == nil {
		return Left
	}
	game := NewBoard(7)
	for game.Current_winner == 0 {
		if len(game.Current_player_path_moves) == 1 {
			for pawn := range game.Current_player_path_moves {
				game.Play(pawn)
			}
			continue
		} else {
			potential_futures := []*Potential_future{}
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
					potential_futures = append(potential_futures, &Potential_future{
						Pawn:  pawn,
						Score: score,
					})
				} else {
					score, err := GetPotentialFutureScore(right_player, current_board_right, potential_board)
					if err != nil {
						panic(err)
					}
					potential_futures = append(potential_futures, &Potential_future{
						Pawn:  pawn,
						Score: score,
					})
				}
			}
			sort.Slice(potential_futures, func(i, j int) bool {
				return potential_futures[i].Score < potential_futures[j].Score
			})

			// Play highest-scoring future
			game.Play(potential_futures[len(potential_futures)-1].Pawn)
		}
	}
	if game.Current_winner == Left {
		left_player.Fitness++
	} else {
		right_player.Fitness++
	}
	return game.Current_winner
}

func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}
