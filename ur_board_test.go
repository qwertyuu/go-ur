package gour

import (
	"testing"
)

func TestBoard(t *testing.T) {
	b := NewBoard(7)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)

	b.Play(-1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)

	b.Play(-1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)

	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)

	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)

	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)

	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(-1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(-1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(-1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(0)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	
	b.Play(1)
	t.Log(b.Current_dice)
	t.Log(b.Current_player)
	t.Log(b.Current_player_path_moves)
	t.Log(b.right_player_out)
	t.Log(b.left_player_out)

}