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

There is an analysis part to this code, in order to (try to) measure the performance of the trained AIs against other AIs or against static policies (such as a random-playing player) for example.

The challenge of this task is to find metrics that predict or indicate real-world performance against humans. The ultimate metric would yield something like a floating-point number from 0 to 1, 0 being a random-playing AI and 1 being the ultimate unbeatable player (given good dice rolls).

### Championship analysis

The goal of this analysis is to measure AIs against each other in an efficient fashion, that is, not just brute-forcing thousands of 1v1s with the same pairs of players. To achieve this, we use double-elimination tournaments (so convenient!) that we run many times in a row. This has the advantage of mixing contenders together and to limit greatly the amount of games that are needed to evaluate the analyzed population. This analysis can take a large array of different AIs trained under different conditions under the only limit that they all fulfill the `Ur_player` interface (this is not yet implemented but it is the goal). If your AIs were trained using this codebase, this is already the case. You can just list the directories where your runs live in the `cmd\analysis\championship\main.go` file and have fun. The code will scan the folders, pick the organisms and launch a series of championships.

One use case would be: You tried 20 different combinations of ur.neat configurations and trained them using the `ur_neat_bootstrap` method. Now, you want to know what configuration led to the best AI out of those 20 combinations. You will then list all 20 or so folders that contain all the AIs that you trained in the `cmd\analysis\championship\main.go` file, run the analysis and take a look at the output. This will rank-order all your AIs in such a way that you will be able to choose which configurations to spend your precious time on and those not to try again.

The output is written to `wins.json` and contains a list of all organisms that were found/listed along with their number of championship wins and proportion of wins (0-1, 1 means it won all tournaments and 0 means it lost all tournaments). Usually, a contender with a high amount of wins for a large amount of tournaments is a (relatively) *good* contender.

NOTE: This method defines a *relative* way to rank-order your AIs. An AI being first against many other AIs might mean that it is good against humans, but it could also mean that all other AIs are VERY BAD and therefore is not a good indicator of real-world performance. So, I advise you use this technique in order to prune bad AIs instead of picking good AIs. The other analysis, the "versus" analysis, seems better suited for picking a *good* AI, but relies on you having picked good candidates before hand.

### Versus analysis

As you read in the last analysis type, a championship analysis is a good way to spot a single or few potentially good AIs out of many. Though, as stated in the NOTE, this is not a good way to find an AI that is objectively good at the game. For this, you need to go to a more simple approach, the 1v1s.

This way, you can go down to telling the code to run many games in a row against a fixed player, whatever its policy is.
As the developer of this code, I think a combination of using your last "best" AI and a Random-playing bot is a good way to have an objective view of the performance of the AI you trained. Usually, if an AI performs better both against your known "best" AI and a Random-playing bot, it *should* do better against a person.

One caveat is that this method is the slowest. Since you need to essentially brute-force the AIs against each other, you need to know which AIs to compare against eachother. Using this analysis on two totally new AIs might be the most computationally expensive way to rank-order your AIs. I would recomment using a Championship analysis in this case.

The strength of this analysis is that it is easy to understand. If you use static policies or a well-known AI as your opponent, this method is good for compiling numbersm keeping track of and comparing your best-playing AIs against eachother.

Another strength could be to compare humans against the metrics you used to compare your AIs and to rank-order humans against your AIs using those metrics. This way, you could statistically infer which is better, the computer or the human.

### Versus random

This analysis is different of the last two as it is an analysis of the game itself and not the AIs.
This analysis is very simple: Suppose 2 Ur players of equal skill over many games, is there a skew in winning or losing depending on which of the players is starting to play. 

In other words, is there an advantage or handicap to be the starting player in Ur.

To know, I coded the versus_random analysis. The "equal skill" players are two random-playing bots with the exact same policy. You can let that analysis run for as many games as you wish. I found that all the probabilities end up tending towards 50% everytime, so I deduced that there is no significant advantage or handicap to being the first or second player.

### Is the game only luck?

Many people which who I played the game stated the following: This is a game of chance. If I get good throws, I will most likely win. But I feel it's also a game of strategy, because my choices seem to influence the outcome of the game. Therefore I am unsure whether the game is only luck, only strategy, or both and if both, how much of both.

There are many hypothesis here.

Here are many things we can deduce:

- If the game is only strategy, a random-playing AI would have a 0% winrate against the world best player
  - This hypothesis can be easily thought about and dismissed as we could imaging the best player in the world only rolling 0s for the whole game, and the random-playing AI rolling normally, which would give a win to the random-playing AI.
- If the game is only luck, a random-playing AI would have the same chance of winning as any other policies
  - This hypothesis is easy to dismiss too as it would be easy to imagine always "eating" pawns of the random-playing AI as they come out their starting zone given the appropriate dice rolls. We cannot, for sure, estimate how much we would win against such a primitive player, but we can at least deduce that the win rate would be more than 50% for the more mature policy.

Therefore, the game is neither only luck and only strategy.

In the numbers I have run through, here are my conclusions: Even my best AI, one that trained over many generations and beat the older best AIs by a large margin (71.08% winrate), the best it can do against a random-playing bot is `96.64%` winrate against a random-playing AI over 10k games. This means that the game tends towards 0% luck but is also more than 0 as we deduced. I would say then that the game is ~3% luck.

A fun thing to think about is this one:

My "best" AI was not the best against a random-playing AI (my best recorded score so far was `96.72%`), but since it did better against the older "best" AI in 1v1 (71.08% winrate vs. `67.98%` winrate for the other best AI against a random-playing AI), it also seemed to generalize better against humans (tested in real games against real peoples).

As you can see here, one of the biggest challenges in this whole code base is to find an appropriate metric that will predict how well an AI will do against a human. Only optimizing for winrate against a random-playing AI will yield an appreciable opponent, but not a very good opponent. Using different reference points and policies seem to be the best approach so far.

I am eager to find better approaches in the future. Please try some out by forking and don't hesitate to open a PR to propose improvements in this matter.

### Other theorized way to rank-order AIs, the ELO score

In the world of chess or online games, there is a kind of metric called the ELO score. This is an avenue that I looked at to eventually implement here as a reference metric to use in my analysis and rankings in order to get the absolute strength of my trained AIs.

I see many challenges with this method, such as keeping a record of the ELO of each AIs somewhere permanent, so that I can read and update them as time goes on. This seems crumbersome, so I let that idea slide.

If used only to compare AIs against eachother, it might suffer that some AIs will be analyzed more than others, so they will get excessively large ELO scores from winning so many games. Also, since there is a luck element to the game as analyzed in the previous section, the ELO score of a very good player could drop significantly against a random-playing bot policy just by bad luck, which is not desirable. Though, this is only a feeling I have, not something I researched.

If used to compare humans and AIs, there would need to be some centralized way of keeping track of player ELOs, which would lead to many problems such as hosting the game freely and having secure gaming environments with trusted ELO counting. This seems like a job for a company to do, or someone very dedicated to hosting and developing a trustworthy Ur platform.

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

Output pawns are referenced by index of the `ai_pawn_positions` input. A pawn of value `-1` means to add a new pawn (play a new pawn)
`future_scores` are ordered ascending, you should have the highest-scoring scored pawn in the "pawn" attribute in the root.

Notes: 

- ai_pawn_positions and enemy_pawn_positions are position that reprensent the index of the path for the player (0 to 13 inclusive for both). They are NOT absolute board positions (TODO: Add an image here to illustrate this concept, along with an explanation)
- "ai" refers to the bot's point of view in this case (ai_pawn_positions means the pawn positions of the bot, enemy is the other player, whatever its nature)

# Compiling to C for python bridge

`go build -buildmode=c-shared -o go-ur_infer.so .\cmd\python-bridge\main.go`

then run `python loadc.py`

# TODO

- argument to pick trained genome for inference
- better readme.md for other people than myself
- support `Ur_player` interface in double-elimination tournaments