package gour

import (
	"github.com/MichaelTJones/pcg"
)

var pcg32 = pcg.NewPCG32()

func RandBool() bool {
	return pcg32.Random()&0x01 == 0
}

const (
	Left  = 1
	Right = -1
)

type board struct {
	board [20]int

	left_player_queue  int
	right_player_queue int

	left_pawns  [7]int
	right_pawns [7]int

	left_pawn_path_positions  [7]int
	right_pawn_path_positions [7]int

	left_player_out  int
	right_player_out int

	pawn_per_player int

	Current_player       int
	Current_player_path_moves map[int]int
	Current_dice              int
}

func NewBoard(pawn_per_player int) *board {

	p := board{
		pawn_per_player:    pawn_per_player,
		left_player_queue:  pawn_per_player,
		right_player_queue: pawn_per_player,
	}

	if RandBool() {
		p.Current_player = Left
	} else {
		p.Current_player = Right
	}

	p.Current_dice = ThrowDice()

	for p.Current_dice == 0 {
		p.Current_player *= -1
		p.Current_dice = ThrowDice()
	}

	p.Current_player_path_moves = p.PlayerValidMoves(p.Current_dice, p.Current_player)
	return &p
}

func ThrowDice() int {
	dice := 0
	for i := 0; i < 4; i++ {
		if RandBool() {
			dice++
		}
	}
	return dice
}

func (r *board) AsArray(for_player int) [20]int {
	var board_array [20]int

	copy(r.board[:], board_array[:])

	if for_player == Right {
		for i := 0; i < len(board_array); i++ {
			board_array[i] = -board_array[i]
		}
	}
	return board_array
}

func leftPlayerPath() [14]int {
	return [14]int{9, 6, 3, 0, 1, 4, 7, 10, 12, 13, 15, 18, 17, 14}
}

func rightPlayerPath() [14]int {
	return [14]int{11, 8, 5, 2, 1, 4, 7, 10, 12, 13, 15, 18, 19, 16}
}

func (r *board) Play(pawn int) {
	pawn_in_play := 0
	var path_positions *[7]int
	var enemy_path_positions *[7]int
	var player_queue *int
	var enemy_player_queue *int
	var path [14]int
	var out *int
	if r.Current_player == Left {
		pawn_in_play = r.pawn_per_player - r.left_player_queue - r.left_player_out
		player_queue = &r.left_player_queue
		enemy_player_queue = &r.right_player_queue
		path_positions = &r.left_pawn_path_positions
		enemy_path_positions = &r.right_pawn_path_positions
		out = &r.left_player_out
		path = leftPlayerPath()
	} else {
		pawn_in_play = r.pawn_per_player - r.right_player_queue - r.right_player_out
		player_queue = &r.right_player_queue
		enemy_player_queue = &r.left_player_queue
		path_positions = &r.right_pawn_path_positions
		enemy_path_positions = &r.left_pawn_path_positions
		out = &r.right_player_out
		path = rightPlayerPath()
	}

	if pawn > pawn_in_play || pawn > r.pawn_per_player {
		panic("Pawn out of range")
	}

	new_pawn_path_position := r.Current_player_path_moves[pawn]

	// Apply move to course and to board
	if pawn == -1 { // new pawn
		new_pawn_board_position := path[new_pawn_path_position]
		*player_queue--
		path_positions[pawn_in_play] = new_pawn_path_position
		pawn_in_play++
		r.board[new_pawn_board_position] = r.Current_player
	} else if new_pawn_path_position == 14 { // out! (14 is len(path))
		current_pawn_board_position := path[path_positions[pawn]]
		*out++
		r.board[current_pawn_board_position] = 0
		// move all pawns to the left to remove current pawn data
		for i := pawn + 1; i < 7; i++ {
			path_positions[i-1] = path_positions[i]
		}
		// fill whatever value was at the end with 0 to create space
		path_positions[len(path_positions)-1] = 0
		pawn_in_play--
	} else { // moving a pawn
		new_pawn_board_position := path[new_pawn_path_position]
		current_pawn_board_position := path[path_positions[pawn]]
		// remove pawn from board where it was
		r.board[current_pawn_board_position] = 0

		// update path position
		path_positions[pawn] = new_pawn_path_position
		if r.board[new_pawn_board_position] == -r.Current_player { // enemy gets eaten!
			*enemy_player_queue++
			enemy_pawn := -1
			for i := 0; i < 7; i++ { // find enemy pawn that is where we want to land
				if enemy_path_positions[i] == new_pawn_path_position {
					enemy_pawn = i
					break
				}
			}
			if enemy_pawn == -1 {
				panic("Enemy pawn not found?")
			}
			// move all pawns to the left to remove current pawn data
			for i := enemy_pawn + 1; i < 7; i++ {
				enemy_path_positions[i-1] = enemy_path_positions[i]
			}
			// fill whatever value was at the end with 0 to create space
			enemy_path_positions[len(enemy_path_positions)-1] = 0
		}

		// Add to board where it now is
		r.board[new_pawn_board_position] = r.Current_player
	}

	r.Current_dice = ThrowDice()

	for r.Current_dice == 0 {
		r.Current_player *= -1
		r.Current_dice = ThrowDice()
	}

	// Player plays again
	if !(new_pawn_path_position == 3 || new_pawn_path_position == 7 || new_pawn_path_position == 13) {
		r.Current_player *= -1
	}

	r.Current_player_path_moves = r.PlayerValidMoves(r.Current_dice, r.Current_player)

}

func (r *board) PlayerValidMoves(dice int, player int) map[int]int {
	if dice == 0 {
		return make(map[int]int, 0)
	}

	var path [14]int
	pawn_in_play := 0
	pawn_in_queue := 0
	var pawns_path_positions [7]int
	if player == Left {
		path = leftPlayerPath()
		pawn_in_play = r.pawn_per_player - r.left_player_queue - r.left_player_out
		pawn_in_queue = r.left_player_queue
		pawns_path_positions = r.left_pawn_path_positions
	} else {
		path = rightPlayerPath()
		pawn_in_play = r.pawn_per_player - r.right_player_queue - r.right_player_out
		pawn_in_queue = r.right_player_queue
		pawns_path_positions = r.right_pawn_path_positions
	}

	// pawn => course position
	possible_course_moves := make(map[int]int, 0)
	if pawn_in_queue > 0 && r.board[path[dice-1]] == 0 {
		possible_course_moves[-1] = dice - 1 // In
	}
	for i := 0; i < pawn_in_play; i++ {
		pawn_position_in_course := pawns_path_positions[i]
		pawn_course_position_after_dice := pawn_position_in_course + dice
		if pawn_course_position_after_dice == 14 { // fixed size, equivalent to len(path)
			possible_course_moves[i] = 14 // Out
		} else if pawn_course_position_after_dice < 14 { // fixed size, equivalent to len(path)
			board_position_after_dice := path[pawn_course_position_after_dice]
			if board_position_after_dice == 10 { // Center rosette
				if r.board[board_position_after_dice] == 0 {
					possible_course_moves[i] = pawn_course_position_after_dice // Can play on center
				}
			} else if r.board[board_position_after_dice] != player {
				possible_course_moves[i] = pawn_course_position_after_dice // Can play anywhere
			}
		}
	}
	return possible_course_moves
}
