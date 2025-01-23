package util

import (
	"fmt"
	"log"
	"strconv"

	"github.com/pterm/pterm"
)

func ShowInputPrompt(label string, defaultValue string) string {
	result, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(label).WithDefaultValue(defaultValue).Show()
	return result
}

func ShowNumberInputPrompt(label string, defaultValue int) int {
	result := ShowInputPrompt(label, fmt.Sprintf("%d", defaultValue))
	num, err := strconv.Atoi(result)
	if err != nil {
		log.Fatalf("输入的不是数字 %v \n", err)
	}
	return num
}
