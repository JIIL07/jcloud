package main

import (
	"bufio"
	"fmt"
	"github.com/JIIL07/cloudFiles-manager/internal/client/cmd"
	"os"
	"strings"
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
