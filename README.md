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


## Help wanted
This project is fully open-sourced, any contribution is greatly welcomed! Help the midjouney community get more accurate modertation result through code!
- Help to enrich the banned words library
- Help to setup prompt test cases in test automation
