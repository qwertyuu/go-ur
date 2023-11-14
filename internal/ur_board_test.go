package gour

import (
	"fmt"
	"testing"
)

func TestBoard(t *testing.T) {
	pcg32.Seed(1, 1)
	b := NewBoard(7)

	count := 0
	for b.Current_winner == 0 {
		for k := range *b.Current_player_path_moves {
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
		count++
		if count == 10 {
			break
		}
	}

	j := b.Copy()
	t.Log(j.Current_player_path_moves)
	t.Log(b.AsArray(Right))
	t.Log(j.AsArray(Right))
	j.Play(0)
	t.Log(b.AsArray(Right))
	t.Log(j.AsArray(Right))

}

func TestThousandGames(t *testing.T) {
	pcg32.Seed(1, 1)
	
	for i := 0; i < 1000; i++ {
		b := NewBoard(7)
		for b.Current_winner == 0 {
			// pick a random move using pcg32
			moves := make([]int, 0, len(*b.Current_player_path_moves))
			for k := range *b.Current_player_path_moves {
				moves = append(moves, k)
			}
			m := pcg32.Bounded(uint32(len(moves)))
			// play the move
			b.Play(int(moves[m]))
		}
	}
}

func TestDice(m *testing.T) {
	zero := 0
	one := 0
	two := 0
	three := 0
	four := 0
	n := float64(1000000)
	for i := 0.0; i < n; i++ {
		dice := throwDice()
		if dice == 0 {
			zero++
		}
		if dice == 1 {
			one++
		}
		if dice == 2 {
			two++
		}
		if dice == 3 {
			three++
		}
		if dice == 4 {
			four++
		}
	}

	fmt.Printf("0:%f, 1:%f, 2:%f, 3:%f, 4:%f\n", float64(zero)/n, float64(one)/n, float64(two)/n, float64(three)/n, float64(four)/n)
}
