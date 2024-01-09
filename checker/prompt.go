package checker

import (
	"strings"

	"github.com/goapi-ai/midjourney-api-prompt-checker/model"
)

func CheckPrompt(prompt string, checkBannedWords bool, proxyUrl string) (result model.PromptCheckResult) {
	prompt, loweredWords, urls := PreprocessPrompt(prompt)
	if checkBannedWords {
		if err := CheckPromptBannedWords(strings.Join(loweredWords, " ")); err != nil {
			result.ErrorMessage = err.Error()
			return
		}
	}
	prompt, aspectRatio, err := CheckPromptParam(prompt, loweredWords)
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	if err := CheckImageUrl(prompt, urls, proxyUrl); err != nil {
		result.ErrorMessage = err.Error()
	}
	result.Prompt = prompt
	result.AspectRatio = aspectRatio
	return
}
