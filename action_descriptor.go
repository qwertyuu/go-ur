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
