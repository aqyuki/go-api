package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func StringFromPrompt(msg string, allowEmpty bool) (string, error) {
	validator := func(input string) error {
		if len(input) == 0 && !allowEmpty {
			return fmt.Errorf("input cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    msg,
		Validate: validator,
	}
	res, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return res, nil
}
