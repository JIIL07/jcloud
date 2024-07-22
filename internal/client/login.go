package cloudfiles

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Login() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	return username, password
}
