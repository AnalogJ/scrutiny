package utils

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func StdinQueryPassword(question string) (string, error) {

	fmt.Println(color.BlueString(question))
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	text := strings.TrimSpace(string(bytePassword))
	return text, nil
}

func StdinQuery(question string) string {

	fmt.Println(color.BlueString(question))
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func StdinQueryBoolean(question string) bool {

	text := StdinQuery(question)
	text = strings.ToLower(text)

	if text == "true" || text == "yes" || text == "y" {
		return true
	} else if text == "false" || text == "no" || text == "n" {
		return false
	} else {
		color.Yellow("WARNING: invalid response only true/yes/y/false/no/n allowed not `%v`.\nAssuming `no`", text)
		return false
	}
}

func StdinQueryInt(question string) (int, error) {

	text := StdinQuery(question)
	return StringToInt(text)
}
