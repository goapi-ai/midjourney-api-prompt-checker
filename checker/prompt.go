package checker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/goapi-ai/midjourney-api-prompt-checker/model"
)

const (
	StyleRawVersion = 5.1

	ErrUnReconizedParam   = "Unrecognized Param"
	ErrInvalidParamValue  = "Invalid Param Value"
	ErrInvalidParamFormat = "Invalid Param Format"
)

func PreprocessPrompt(prompt string) (string, []string) {
	// for Apple devices
	prompt = strings.ReplaceAll(prompt, "—", "--")
	words := strings.Fields(strings.Trim(prompt, " "))
	loweredWords := make([]string, len(words))
	for i, word := range words {
		loweredWords[i] = strings.ToLower(word)
	}
	return strings.Join(loweredWords, " "), loweredWords
}

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

func CheckAspectParam(param string) bool {
	aspects := strings.Split(param, ":")
	if len(aspects) != 2 {
		return false
	}
	if _, err := strconv.Atoi(aspects[0]); err != nil {
		return false
	}
	if _, err := strconv.Atoi(aspects[1]); err != nil {
		return false
	}
	return true
}

func CheckZoomParam(param string) bool {
	ratio, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return false
	}
	if ratio < 1.0 || ratio > 2.0 {
		return false
	}
	return true
}

func CheckRepeatParam(param string) bool {
	repeat, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	if repeat < 2 || repeat > 40 {
		return false
	}
	return true
}

func CheckPermutation(prompt string) error {
	if strings.Contains(prompt, "{") || strings.Contains(prompt, "}") {
		return errors.New("Permutation Not Supported")
	}
	return nil
}

func CheckSpaces(prompt string) error {
	curPrompt := prompt
	for i := strings.Index(curPrompt, "--"); i != -1; i = strings.Index(curPrompt, "--") {
		// -- must follows a space
		if i == 0 || curPrompt[i-1] != ' ' {
			return fmt.Errorf("%s: there should be space before --", ErrInvalidParamFormat)
		}
		// param must follow --, no space allowed in between
		if i+2 < len(curPrompt) && curPrompt[i+2] == ' ' {
			return fmt.Errorf("%s: there should be no space after --", ErrInvalidParamFormat)
		}
		curPrompt = curPrompt[i+2:]
	}
	return nil
}

func CheckParamLegal(param string) bool {
	for _, allowedParam := range model.Params {
		if param == allowedParam {
			return true
		}
	}
	return false
}

func RemoveUnsupportParams(prompt string, params []string) string {
	for _, param := range params {
		prompt = strings.ReplaceAll(prompt, param, "")
	}
	return prompt
}

func CheckPromptParam(prompt string, words []string) (newPrompt, aspectRatio string, err error) {
	newPrompt = prompt
	// permutation param is not supported now
	if err = CheckPermutation(newPrompt); err != nil {
		return
	}
	if !strings.Contains(newPrompt, "--") {
		return
	}
	if err = CheckSpaces(newPrompt); err != nil {
		return
	}

	var (
		styleValue      string
		versionValue    float64
		unsupportParams []string
	)
	for index, subString := range words {
		if strings.HasPrefix(subString, "--") {
			param := strings.TrimPrefix(subString, "--")
			if !CheckParamLegal(param) {
				err = fmt.Errorf("%s: --%s", ErrUnReconizedParam, param)
				return
			}
			if param == "turbo" || param == "fast" {
				// not supporeted to change mode with param
				unsupportParams = append(unsupportParams, fmt.Sprintf(" --%s", param))
			}
			if index < len(words)-1 {
				value := words[index+1]
				if param == "repeat" || param == "r" {
					if !CheckRepeatParam(value) {
						err = fmt.Errorf("%s: --%s %s", ErrInvalidParamValue, param, value)
						return
					}
					// repeat param is not supported now
					unsupportParams = append(unsupportParams, fmt.Sprintf(" --%s %s", param, value))
				}
				if param == "aspect" || param == "ar" {
					if !CheckAspectParam(value) {
						err = fmt.Errorf("%s: --%s %s", ErrInvalidParamValue, param, value)
						return
					}
					// aspect ratio will be stored at extra param, used as default value lator in Custom Zoom/Remix Actions
					aspectRatio = value
					unsupportParams = append(unsupportParams, fmt.Sprintf(" --%s %s", param, value))
				}
				if param == "version" || param == "v" {
					versionValue, err = strconv.ParseFloat(value, 64)
					if err != nil {
						err = fmt.Errorf("%s: --%s %s", ErrInvalidParamValue, param, value)
						return
					}
				}
				if param == "style" {
					styleValue = value
				}
			}
		}
	}
	if versionValue >= StyleRawVersion && styleValue != "raw" {
		newPrompt += " --style raw"
		if styleValue != "" {
			unsupportParams = append(unsupportParams, fmt.Sprintf(" --style %s", styleValue))
		}
	}
	newPrompt = RemoveUnsupportParams(newPrompt, unsupportParams)
	return
}

type PromptCheckResult struct {
	Prompt       string
	ErrorMessage string
	AspectRatio  string
}

func CheckPrompt(prompt string) PromptCheckResult {
	prompt, words := PreprocessPrompt(prompt)
	if err := CheckPromptBannedWords(prompt); err != nil {
		return PromptCheckResult{
			ErrorMessage: err.Error(),
		}
	}
	prompt, aspectRatio, err := CheckPromptParam(prompt, words)
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	return PromptCheckResult{
		Prompt:       prompt,
		ErrorMessage: errorMessage,
		AspectRatio:  aspectRatio,
	}
}
