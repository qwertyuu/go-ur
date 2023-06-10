package gour

type potential_board_descriptor struct { // 28 features
	board_state        [20]int
	my_pawn_in_play    float64 // 0-1 percent from total
	my_pawn_queue      float64 // 0-1 percent from total
	my_pawn_out        float64 // 0-1 percent from total
	enemy_pawn_in_play float64 // 0-1 percent from total
	enemy_pawn_queue   float64 // 0-1 percent from total
	enemy_pawn_out     float64 // 0-1 percent from total

	winner int // 1: me, 0, no one, -1 ennemy
	turn   int // 1: me, -1 ennemy
}

func GetPotentialBoardDescriptor(b *board, for_player int) *potential_board_descriptor {
	var my_pawn_in_play int
	var my_pawn_queue int
	var my_pawn_out int
	var enemy_pawn_in_play int
	var enemy_pawn_queue int
	var enemy_pawn_out int
	var winner int = 0
	var turn int = 0

	if for_player == Left {
		my_pawn_in_play = b.pawn_per_player - b.left_player_queue - b.left_player_out
		my_pawn_queue = b.left_player_queue
		my_pawn_out = b.left_player_out
		enemy_pawn_in_play = b.pawn_per_player - b.right_player_queue - b.right_player_out
		enemy_pawn_queue = b.right_player_queue
		enemy_pawn_out = b.right_player_out
		if b.Current_winner == Left {
			winner = 1
		} else if b.Current_winner == Right {
			winner = -1
		}
		if b.Current_player == Left {
			turn = 1
		} else if b.Current_player == Right {
			turn = -1
		}
	} else {
		my_pawn_in_play = b.pawn_per_player - b.right_player_queue - b.right_player_out
		my_pawn_queue = b.right_player_queue
		my_pawn_out = b.right_player_out
		enemy_pawn_in_play = b.pawn_per_player - b.left_player_queue - b.left_player_out
		enemy_pawn_queue = b.left_player_queue
		enemy_pawn_out = b.left_player_out
		if b.Current_winner == Right {
			winner = 1
		} else if b.Current_winner == Left {
			winner = -1
		}
		if b.Current_player == Right {
			turn = 1
		} else if b.Current_player == Left {
			turn = -1
		}
	}

	f_pawn_per_player := float64(b.pawn_per_player)
	return &potential_board_descriptor{
		board_state:        b.AsArray(for_player),
		my_pawn_in_play:    float64(my_pawn_in_play) / f_pawn_per_player,
		my_pawn_queue:      float64(my_pawn_queue) / f_pawn_per_player,
		my_pawn_out:        float64(my_pawn_out) / f_pawn_per_player,
		enemy_pawn_in_play: float64(enemy_pawn_in_play) / f_pawn_per_player,
		enemy_pawn_queue:   float64(enemy_pawn_queue) / f_pawn_per_player,
		enemy_pawn_out:     float64(enemy_pawn_out) / f_pawn_per_player,
		turn:               turn,
		winner:             winner,
	}
}