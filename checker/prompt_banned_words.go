package checker

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/goapi-ai/midjourney-api-prompt-checker/model"
)

func CheckPromptBannedWords(prompt string) error {
	words := strings.FieldsFunc(prompt, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, bannedWord := range model.BannedWords {
		// check if the prompt contains banned phrase, and check if any word is banned
		if strings.Contains(bannedWord, " ") && strings.Contains(prompt, bannedWord) {
			return fmt.Errorf("Banned Prompt: %s", bannedWord)
		} else {
			for _, word := range words {
				if word == bannedWord {
					return fmt.Errorf("Banned Prompt: %s", bannedWord)
				}
			}
		}
	}
	return nil
}
