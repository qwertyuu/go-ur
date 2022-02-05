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
	"ai_pawn_out": 1,
	"enemy_pawn_out": 1,
	"dice": 1,
	"ai_pawn_positions": [7, 13],
	"enemy_pawn_positions": [6, 12]
}
```

Output format:
```json
{
	"pawn": 1,
	"future_scores": [
		{
			"pawn": -1,
			"score": -1,
		},
		{
			"pawn": 0,
			"score": -1,
		},
		{
			"pawn": 1,
			"score": 1,
		}
	]
}
```

Output pawns are referenced by index of the `ai_pawn_positions` input.
`future_scores` are ordered ascending, you should have the latest in the root "pawn".

Notes: 

- ai_pawn_positions and enemy_pawn_positions are position that reprensent the index of the path for the player (0 to 13 inclusive for both). They are NOT absolute board positions
- "ai" refers to the bot's point of view in this case (ai_pawn_positions means the pawn positions of the bot, enemy is the other player)

TODO:
- run tournament in goroutines
- mess with the neat configurations
- argument to pick trained genome for inference
- run tournament between all winners of all generations to find best of "best"!
- run tournament against "random"-playing bots for a baseline