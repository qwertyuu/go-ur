package gour

type current_board_descriptor struct {
	board_state        [20]int
	my_pawn_in_play    float32 // 0-1 percent from total
	my_pawn_queue      float32 // 0-1 percent from total
	my_pawn_out        float32 // 0-1 percent from total
	enemy_pawn_in_play float32 // 0-1 percent from total
	enemy_pawn_queue   float32 // 0-1 percent from total
	enemy_pawn_out     float32 // 0-1 percent from total
}

type potential_board_descriptor struct {
	board_state        [20]int
	my_pawn_in_play    float32 // 0-1 percent from total
	my_pawn_queue      float32 // 0-1 percent from total
	my_pawn_out        float32 // 0-1 percent from total
	enemy_pawn_in_play float32 // 0-1 percent from total
	enemy_pawn_queue   float32 // 0-1 percent from total
	enemy_pawn_out     float32 // 0-1 percent from total

	winner int // 1: me, 0, no one, -1 ennemy
	turn   int // 1: me, -1 ennemy
}

func GetPotentialBoardDescriptor(b *board, for_player int) potential_board_descriptor {
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

	f_pawn_per_player := float32(b.pawn_per_player)
	return potential_board_descriptor{
		board_state:        b.AsArray(for_player),
		my_pawn_in_play:    float32(my_pawn_in_play) / f_pawn_per_player,
		my_pawn_queue:      float32(my_pawn_queue) / f_pawn_per_player,
		my_pawn_out:        float32(my_pawn_out) / f_pawn_per_player,
		enemy_pawn_in_play: float32(enemy_pawn_in_play) / f_pawn_per_player,
		enemy_pawn_queue:   float32(enemy_pawn_queue) / f_pawn_per_player,
		enemy_pawn_out:     float32(enemy_pawn_out) / f_pawn_per_player,
		turn:               turn,
		winner:             winner,
	}
}

func GetCurrentBoardDescriptor(b *board, for_player int) current_board_descriptor {
	var my_pawn_in_play int
	var my_pawn_queue int
	var my_pawn_out int
	var enemy_pawn_in_play int
	var enemy_pawn_queue int
	var enemy_pawn_out int

	if for_player == Left {
		my_pawn_in_play = b.pawn_per_player - b.left_player_queue - b.left_player_out
		my_pawn_queue = b.left_player_queue
		my_pawn_out = b.left_player_out
		enemy_pawn_in_play = b.pawn_per_player - b.right_player_queue - b.right_player_out
		enemy_pawn_queue = b.right_player_queue
		enemy_pawn_out = b.right_player_out
	} else {
		my_pawn_in_play = b.pawn_per_player - b.right_player_queue - b.right_player_out
		my_pawn_queue = b.right_player_queue
		my_pawn_out = b.right_player_out
		enemy_pawn_in_play = b.pawn_per_player - b.left_player_queue - b.left_player_out
		enemy_pawn_queue = b.left_player_queue
		enemy_pawn_out = b.left_player_out
	}

	f_pawn_per_player := float32(b.pawn_per_player)
	return current_board_descriptor{
		board_state:        b.AsArray(for_player),
		my_pawn_in_play:    float32(my_pawn_in_play) / f_pawn_per_player,
		my_pawn_queue:      float32(my_pawn_queue) / f_pawn_per_player,
		my_pawn_out:        float32(my_pawn_out) / f_pawn_per_player,
		enemy_pawn_in_play: float32(enemy_pawn_in_play) / f_pawn_per_player,
		enemy_pawn_queue:   float32(enemy_pawn_queue) / f_pawn_per_player,
		enemy_pawn_out:     float32(enemy_pawn_out) / f_pawn_per_player,
	}
}

/*
Rosette (rejouer)
Position absolue, 0 sortie et 1 a 14 position
Libère une rosette
Libère le centre
Arrive au centre
Manger un ennemi
Ya ti un pion ennemi a 1,2,3,4 cases devant
Ya ti un pion ennemi a 1,2,3,4 cases derrière
Ya ti un pion allié a 1,2,3,4 cases devant
Ya ti un pion allié à 1,2,3,4 cases derrière
Zone de combat
Ajouter un nouveau pion
*/
