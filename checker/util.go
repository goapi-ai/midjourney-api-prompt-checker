package checker

import "strings"

func PreprocessPrompt(prompt string) (string, []string) {
	// for Apple devices
	prompt = strings.ReplaceAll(prompt, "—", "--")
	words := strings.Fields(strings.Trim(prompt, " "))
	loweredWords := make([]string, len(words))
	for i, word := range words {
		loweredWords[i] = strings.ToLower(word)
	}
	return strings.Join(words, " "), loweredWords
}
