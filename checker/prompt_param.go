package checker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/goapi-ai/midjourney-api-prompt-checker/model"
)

const (
	ErrPermutationUnsupported = "Permutation Not Supported"
	ErrUnReconizedParam       = "Unrecognized Param"
	ErrInvalidParamValue      = "Invalid Param Value"
	ErrInvalidParamFormat     = "Invalid Param Format"
)

// basic param check, according to https://docs.midjourney.com/docs/parameter-list

// --aspect, or --ar Change the aspect ratio of a generation.
// v5: Aspect ratios greater than 2:1 are experimental and may produce unpredictable results.
// v6: range 1:3–3:1
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

// --chaos <number 0–100> Change how varied the results will be.
// Higher values produce more unusual and unexpected generations.
func CheckChaosParam(param string) bool {
	chaos, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	if chaos < 0 || chaos > 100 {
		return false
	}
	return true
}

// --iw <0–2> Sets image prompt weight relative to text weight. The default value is 1.
func CheckIWParam(param string) bool {
	iw, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return false
	}
	if iw < 0 || iw > 2 {
		return false
	}
	return true
}

// --quality <.25, .5, or 1>, or --q <.25, .5, or 1> How much rendering quality time you want to spend.
// The default value is 1. Higher values use more GPU minutes; lower values use less.
func CheckQualityParam(param string) bool {
	if param == ".25" || param == ".5" || param == "1" {
		return true
	}
	return false
}

// --repeat <1–40>, or --r <1–40> Create multiple Jobs from a single prompt.
// --repeat is useful for quickly rerunning a job multiple times.
func CheckRepeatParam(param string) bool {
	repeat, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	if repeat < 1 || repeat > 40 {
		return false
	}
	return true
}

// --seed <integer between 0–4294967295> The Midjourney bot uses a seed number to create a field of visual noise, like television static, as a starting point to generate the initial image grids.
// Seed numbers are generated randomly for each image but can be specified with the --seed or --sameseed parameter. Using the same seed number and prompt will produce similar ending images.
func CheckSeedParam(param string) bool {
	_, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return false
	}
	return true
}

// --stop <integer between 10–100> Use the --stop parameter to finish a Job partway through the process.
// Stopping a Job at an earlier percentage can create blurrier, less detailed results.
func CheckStopParam(param string) bool {
	stop, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	if stop < 10 || stop > 100 {
		return false
	}
	return true
}

// --stylize <number>, or --s <number> parameter influences how strongly Midjourney's default aesthetic style is applied to Jobs.
// 0-1000. The default value is 100. Higher values produce more stylized results.
func CheckStylizeParam(param string) bool {
	stylize, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	if stylize < 0 || stylize > 1000 {
		return false
	}
	return true
}

// --weird <number 0–3000>, or --w <number 0–3000> Explore unusual aesthetics with the experimental --weird parameter.
func CheckWeirdParam(param string) bool {
	weird, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	if weird < 0 || weird > 3000 {
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

// special check

func CheckPermutation(prompt string) error {
	if strings.Contains(prompt, "{") || strings.Contains(prompt, "}") {
		return errors.New(ErrPermutationUnsupported)
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
	// skip following checks if no param in prompt
	if !strings.Contains(newPrompt, "--") {
		return
	}
	// check special spaces cases between -- and param
	if err = CheckSpaces(newPrompt); err != nil {
		return
	}

	var (
		unsupportParams []string
		value           string
	)
	for index, subString := range words {
		if !strings.HasPrefix(subString, "--") {
			continue
		}

		// parse and check param in param list
		param := strings.TrimPrefix(subString, "--")
		if !CheckParamLegal(param) {
			err = fmt.Errorf("%s: --%s", ErrUnReconizedParam, param)
			return
		}

		// parse value if exists
		if index < len(words)-1 {
			value = words[index+1]
		}

		// check each param and corresponding value
		if param == "aspect" || param == "ar" {
			if !CheckAspectParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 1:1", ErrInvalidParamValue, param, value)
				return
			}
			// aspect ratio will be stored at extra param, used as default value lator in Custom Zoom/Remix Actions
			aspectRatio = value
			unsupportParams = append(unsupportParams, fmt.Sprintf(" --%s %s", param, value))
		}
		if param == "chaos" || param == "c" {
			if !CheckChaosParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 0, Range: 0-100", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "iw" {
			if !CheckIWParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 1, Range: 0-2", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "quality" || param == "q" {
			if !CheckQualityParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 1, Range: .25/.5/1", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "repeat" || param == "r" {
			if !CheckRepeatParam(value) {
				err = fmt.Errorf("%s: --%s %s. Range: 1-40", ErrInvalidParamValue, param, value)
				return
			}
			// repeat param is not supported now, auto removed
			unsupportParams = append(unsupportParams, fmt.Sprintf(" --%s %s", param, value))
		}
		if param == "seed" {
			if !CheckSeedParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: Random, Range: 0-4294967295", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "stop" {
			if !CheckStopParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 100, Range: 10-100", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "stylize" || param == "s" {
			if !CheckStylizeParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 100, Range: 0-1000", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "weird" || param == "w" {
			if !CheckWeirdParam(value) {
				err = fmt.Errorf("%s: --%s %s. Default: 0, Range: 0-3000", ErrInvalidParamValue, param, value)
				return
			}
		}
		if param == "version" || param == "v" {
			_, err = strconv.ParseFloat(value, 64)
			if err != nil {
				err = fmt.Errorf("%s: --%s %s", ErrInvalidParamValue, param, value)
				return
			}
		}
	}
	newPrompt = RemoveUnsupportParams(newPrompt, unsupportParams)
	return
}
