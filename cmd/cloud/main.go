package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cmd "github.com/JIIL07/cloudFiles-manager/internal/cmd"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		dir, _ := os.Getwd()
		fmt.Printf("%v>", dir)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		args := strings.Split(input, " ")
		cmd.RootCmd.SetArgs(args[1:])

		cmd.RootCmd.Execute()
	}
}
