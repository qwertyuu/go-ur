package gour

import (
	"math/rand"
	"sort"
)

type double_elimination struct {
	Contenders               []Ur_player
	Losers_of_winner_bracket []Ur_player
	Loser_bracket            []Ur_player
	Winner_bracket           []Ur_player
	Champion                 Ur_player
	Contender_Amount         int
}

// NB de joutes dans la winner bracket avant la loser's bracket = len(contenders) / 2 - 1
// Commencer par évaluer la winner bracket au complet. Noter les perdants de chaque joute
// L'ordre de correspondance entre les L et les losers est inversée. Les premiers de la losers bracket affrontent les derniers de la winners bracket
// Les affrontements entre les losers et le L se fait à toutes les deux séries de joutes jusqu'à détermination d'un gagnant, qui va se battre contre le gagnant de la winner's bracket
// TODO: Make the parameter a list of Ur_players instead of a genetics.Organism so we can have tournaments including any policies
// Caveat: The number of contenders must be a power of two, else there will be empty players to fill the roster
// Caveat: At the end, the game is not repeated if the winner is from the loser's bracket, as to make the tournament a fixed number of games
func EvaluateDoubleEliminationTournament(contenders []Ur_player, pawn_amt int) double_elimination {
	contender_power_of_two := getNearestPowerOfTwo(len(contenders))
	has_nil_contenders := contender_power_of_two != len(contenders)
	for contender_power_of_two > len(contenders) {
		contenders = append(contenders, nil)
	}
	tournament := double_elimination{
		Contenders:               contenders,
		Losers_of_winner_bracket: make([]Ur_player, 0),
		Loser_bracket:            make([]Ur_player, 0),
		Winner_bracket:           make([]Ur_player, 0),
		Contender_Amount:         contender_power_of_two,
	}
	rand.Shuffle(len(tournament.Contenders), func(i, j int) {
		tournament.Contenders[i], tournament.Contenders[j] = tournament.Contenders[j], tournament.Contenders[i]
	})

	// determine loser and winner brackets
	var left_player Ur_player = nil
	left_player_set := false
	for _, right_player := range tournament.Contenders {
		if right_player != nil {
			right_player.SetWinner(false)
		}
		if !left_player_set {
			left_player = right_player
			left_player_set = true
			continue
		}
		left_wins, right_wins := OneVSOne(left_player, right_player, pawn_amt, 1)
		if left_wins > right_wins {
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
		next_winner_bracket := []Ur_player{}
		for _, right_player := range tournament.Winner_bracket {
			if !left_player_set {
				left_player = right_player
				left_player_set = true
				continue
			}
			left_wins, right_wins := OneVSOne(left_player, right_player, pawn_amt, 1)
			if left_wins > right_wins {
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

	// inverse loser bracket contestants to mix things up
	for i, j := 0, len(tournament.Loser_bracket)-1; i < j; i, j = i+1, j-1 {
		tournament.Loser_bracket[i], tournament.Loser_bracket[j] = tournament.Loser_bracket[j], tournament.Loser_bracket[i]
	}
	winner_bracket_loser_pointer := 0
	// evaluate loser bracket
	for len(tournament.Loser_bracket) > 1 {
		next_loser_bracket := []Ur_player{}
		for _, right_player := range tournament.Loser_bracket {
			if !left_player_set {
				left_player = right_player
				left_player_set = true
				continue
			}
			left_wins, right_wins := OneVSOne(left_player, right_player, pawn_amt, 1)
			loser_of_winner_bracket := tournament.Losers_of_winner_bracket[winner_bracket_loser_pointer]
			winner_bracket_loser_pointer++
			if left_wins > right_wins {
				left_wins, right_wins := OneVSOne(left_player, loser_of_winner_bracket, pawn_amt, 1)
				if left_wins > right_wins {
					next_loser_bracket = append(next_loser_bracket, left_player)
				} else {
					next_loser_bracket = append(next_loser_bracket, loser_of_winner_bracket)
				}
			} else {
				left_wins, right_wins := OneVSOne(right_player, loser_of_winner_bracket, pawn_amt, 1)
				if left_wins > right_wins {
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
	left_wins, right_wins := OneVSOne(tournament.Winner_bracket[0], tournament.Loser_bracket[0], 7, pawn_amt)
	if left_wins > right_wins {
		tournament.Champion = tournament.Winner_bracket[0]
	} else {
		tournament.Champion = tournament.Loser_bracket[0]
	}

	// Note that, contrary to the popular definition of double-elimination tournaments, here I chose not to make the loser-bracket winner win twice in order to win the championship.
	// This is due to the property of this kind of tournament to have a fixed number of faceoffs, so easier to mathematically anticipate.

	if has_nil_contenders {
		new_contenders := []Ur_player{}
		for _, contender := range tournament.Contenders {
			if contender != nil {
				new_contenders = append(new_contenders, contender)
			}
		}
		tournament.Contenders = new_contenders
	}

	sort.Slice(tournament.Contenders, func(i, j int) bool {
		return tournament.Contenders[i].GetWins() < tournament.Contenders[j].GetWins()
	})

	return tournament
}

func getNearestPowerOfTwo(i int) int {
	// Note, this will give the nearest upper power of two, not just the nearest power of two
	// So the return value is always >= i
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

func OneVSOne(left_player Ur_player, right_player Ur_player, number_of_pawns int, number_of_games int) (int, int) {
	// Handle nil players. This code exists to support double-elimination tournaments that do not have a power-of-two amount of contenders
	if left_player == nil {
		if right_player != nil {
			right_player.IncrementWins(number_of_games)
		}
		return 0, number_of_games
	}
	if right_player == nil {
		if left_player != nil {
			left_player.IncrementWins(number_of_games)
		}
		return number_of_games, 0
	}
	left_wins := 0
	right_wins := 0
	left_lost := 0
	right_lost := 0
	winner := make(chan int)
	for i := 0; i < number_of_games; i++ {
		l_cp := left_player.Copy()
		r_cp := right_player.Copy()
		go RoutineFight(l_cp, r_cp, number_of_pawns, winner)
	}
	for i := 0; i < number_of_games; i++ {
		// Await winner info
		Current_winner := <-winner
		if Current_winner == Left {
			//fmt.Printf("%s wins after %d moves\n", left_player.GetName(), moves)
			left_wins++
			right_lost++
		} else {
			//fmt.Printf("%s wins after %d moves\n", right_player.GetName(), moves)
			right_wins++
			left_lost++
		}
	}
	left_player.IncrementWins(left_wins)
	left_player.IncrementLosses(left_lost)
	right_player.IncrementWins(right_wins)
	right_player.IncrementLosses(right_lost)
	return left_wins, right_wins
}

func FightUntilWon(board *board, left_player Ur_player, right_player Ur_player) {
	for board.Current_winner == 0 {
		var current_player Ur_player
		if board.Current_player == Left {
			current_player = left_player
		} else {
			current_player = right_player
		}
		board.Play(current_player.GetMove(board))
	}
}

func RoutineFight(left_player Ur_player, right_player Ur_player, number_of_pawns int, winner chan int) {
	board := NewBoard(number_of_pawns)
	FightUntilWon(board, left_player, right_player)
	winner <- board.Current_winner
}

func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}
