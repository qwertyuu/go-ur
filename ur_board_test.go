package gour

import (
	"testing"
)

func TestBoard(t *testing.T) {
	b := NewBoard(7)

	for b.Current_winner == 0 {
		for k := range b.Current_player_path_moves {
			b.Play(k)
			//t.Log(b.String())
			t.Log(b.AsArray(Right))
			//print(b.String())
			//print(b.Current_dice)
			//print("\n")
			//t.Log(b.Current_player)
			//t.Log(b.Current_player_path_moves)
			//t.Log(b.right_player_out)
			//t.Log(b.left_player_out)
			break
		}
	}

}
