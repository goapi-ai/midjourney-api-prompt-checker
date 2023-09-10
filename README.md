# midjourney-api-prompt-checker
A golang package to filter midjourney prompt by checking banned words and validating params. It's battle-tested in [GoAPI Midjourney API](https://www.goapi.ai/midjourney-api).

## Sources of truth
The code strictly followed a bunch of Midjourney materials to perform a proper prompt check, which are:
- [Midjourney Banned Prompt](https://github.com/PlexPt/midjourney-banned-prompt)
- [Midjourney Parameter List](https://docs.midjourney.com/docs/parameter-list)

## Disclaimer
- This project does not use any AI function, its just simple rule and schema. Therefore, some validation results maybe different from Midjourney's AI moderator.
- There's a 'soft-ban' mechanism in Midjourney(Midjourney bot will not cancel the task but send you result in ephemeral message), this project does not handle the 'soft ban' scenario.

## How to use
Run the test program by executing `go run main.go`.

The program will need you to enter an input. You can input test prompt such as `a cute cat --ar 16:9 --v 5.2`.

After entering an input, the program will process the action and display the results:
- `prompt`: processed prompt, lowered, trimmed, etc.
- `aspect ratio`: extracted aspect ratio from input
- `Banned Prompt`: words that not allowed in Midjourney, you can check the list at model/banned_word.go
- `Invalid Param Format`/`Invalid Param Value`: input prompt not following Midjourney param rules


## Help wanted
This project is fully open-sourced, any contribution is greatly welcomed! Help the midjouney community get more accurate modertation result through code!
- Help to enrich the banned words library
- Help to setup prompt test cases in test automation
