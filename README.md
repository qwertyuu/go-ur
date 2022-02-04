# go-ur
 go ur bot!

## Training
Run `go run cmd/train/main.go`. The results are going to be output in the `/out` directory.

## Inference
Find a genome you want to use for inference in the `/out` folder, replace it in `cmd/inference/ur_inference_server.go` and then

Run `go run cmd/inference/ur_inference_server.go` and call `http://localhost:8090/infer` using a payload like this:

```json
{
	"pawn_per_player": 3,
	"my_pawn_out": 1,
	"enemy_pawn_out": 1,
	"dice": 1,
	"my_pawn_positions": [7, 13],
	"enemy_pawn_positions": [6, 12]
}
```

my_pawn_positions and enemy_pawn_positions are position that reprensent the index of the path for the player (0 to 13 inclusive for both). They are NOT absolute board positions

TODO:
- run tournament in goroutines
- mess with the neat configurations