package checker

import (
	"strings"

	"mvdan.cc/xurls/v2"
)

func PreprocessPrompt(originPrompt string) (string, []string, []string) {
	// for Apple devices
	prompt := strings.ReplaceAll(originPrompt, "—", "--")
	prompt = strings.Trim(prompt, " ")
	words := strings.Fields(prompt)
	loweredWords := make([]string, len(words))
	for i, word := range words {
		loweredWords[i] = strings.ToLower(word)
	}
	urls := xurls.Strict().FindAllString(originPrompt, -1)
	return strings.Join(words, " "), loweredWords, urls
}
