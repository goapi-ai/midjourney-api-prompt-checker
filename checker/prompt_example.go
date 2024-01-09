package checker

import (
	"bufio"
	"fmt"
	"os"
)

func RunPromptCheckerExample() {
	fmt.Println("prompt checker example starts. input `quit` to exit")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> please input prompt: ")
		scanner.Scan()
		prompt := scanner.Text()
		if prompt == "quit" {
			break
		}
		checkResult := CheckPrompt(prompt, true, "")
		if checkResult.ErrorMessage != "" {
			fmt.Println(checkResult.ErrorMessage)
		} else {
			fmt.Printf("prompt: %s\naspect ratio: %s\n", checkResult.Prompt, checkResult.AspectRatio)
		}
	}
}
