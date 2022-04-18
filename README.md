# go-ur

A "framework" to train bots to play The Royal Game of Ur, using NEAT, and to measure/analyze their performance.

## Training

There are 3 ways, for now, to train bots to play Ur using this code.
You can find them under `cmd/`, they are:

- `train_bootstrap`
- `train_vs_ai`
- `train_vs_ai_evolved`

### train_bootstrap

This is the only way to train against other generations themselves, starting from scratch (hence "bootstrap").

The concept is this one, defined in `internal/ur_neat_bootstrap.go`:

- Make all organisms of the population challenge each other in a double-elimination tournament (see `internal/double_elim.go` for more details about the tournaments)
- When an organism wins one fight, it gets a point
- Pick the organism with most won fight at the end of the tournament and use as the "best" organism
- End the run if this "best" organism fulfills fitness expectations

This usually yields okay-ish contenders but nothing unbeatable, unless you are very lucky. This is the best way to get a reference AI to train other generations using other methods.

You can mess with the evaluation of this training process in `internal/ur_neat_bootstrap.go` and play with the max fitness for a winner or the amount of games to play against each other.

### train_vs_ai

This is a fixed way to train an AI using 1v1s against one or more reference AIs (anything that implements the `Ur_player` interface from `internal/ur_player.go` can be used as a reference player)

The concept is this one, defined in `internal/ur_neat_vs_ai.go`:

- Define a list of reference AIs (either using a fixed policy like a Random-player or fixed "first-possible-move-picker" player, or another trained AI)
- Execute a bunch of 1v1s of the organism vs. the reference player and note the wins (noting the wins is done automatically by the `OneVSOne` function)
- Rank-order the organisms by their fitness (or wins) and pick the best as a measure
- End the run if this "best" organism fulfills fitness expectations

This is a way to get good AIs but you need to be very thoughtful of the metrics you use to pick a winner and which reference AIs you use, too. I've had abysmal results and promising results just by using different reference AIs for exemple.

### train_vs_ai_evolved

This is the most advanced way to train an AI with this code. It works almost the same way as `train_vs_ai` with the added benefit of scaling the number of winning games as the training goes on.

The concept is this one, defined in `internal/ur_neat_vs_ai.go` (note that it is the same as train_vs_ai, but with a twist):

- Define a list of reference AIs (either using a fixed policy like a Random-player or fixed "first-possible-move-picker" player, or another trained AI)
- Execute a bunch of 1v1s of the organism vs. the reference player and note the wins (noting the wins is done automatically by the `OneVSOne` function)
- Rank-order the organisms by their fitness (or wins) and pick the best as a measure
- When a winner is found, double the number of games played in the next iteration

The last step is crucial for this "evolved" part. First, the AI will play one set of games against the reference AIs (note that one game means 1 * number of reference AIs * number of games per reference AIs, like 3, 5 or 7 pawns. The way it is coded right now, there are 2 reference AIs with 3 rounds each (3, 5, then 7 pawns) so the first round of 1 set of games is actually 1*3*2 which is 6 standoffs). When a good organism is found, the amount of set of games to play is doubled. This means, in our example, that our organisms will play then 12 games (2*3*2) against the reference AIs. This will go on until the `num_generations` parameter from `data/ur.neat` generations have been attained because this method has no fitness stopper like the other two. It will continue to set the bar higher as time goes on. So if you want a super-player AI, just put the `num_generations` way up and mess with other settings to get a sustainable evolution for many many generations, kick back and enjoy.


This is the best way to get a very good AI player that I found. The evolved part of the learning process makes it easier for the computer to find fit organisms for large amount of wins against the reference AIs instead of focusing on winning a static amount of games. It still requires a long time to train but actually gets good results. This is the way the AI in production right now was trained for `https://ur.whnet.ca` under the bot name `Neato`.


## Analysis

TODO :)

(note for myself: talk about analysis_championship and analysis_versus and why use them)


## Inference
Find a genome you want to use for inference in the `out/` or `trained/` folder (note that the `out/` folder is only available locally on your computer and `trained/` contains a set of AIs that were pre-trained and kept for reference purposes), replace it in `cmd/inference/ur_inference_server.go` and then

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
`future_scores` are ordered ascending, you should have the highest-scoring scored pawn in the "pawn" attribute in the root.

Notes: 

- ai_pawn_positions and enemy_pawn_positions are position that reprensent the index of the path for the player (0 to 13 inclusive for both). They are NOT absolute board positions
- "ai" refers to the bot's point of view in this case (ai_pawn_positions means the pawn positions of the bot, enemy is the other player)

TODO:
- argument to pick trained genome for inference
- better readme.md for other people than myself