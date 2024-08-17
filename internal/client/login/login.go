package login

import (
	"fmt"
	"os"
)

func Login() (string, string) {
	var username, password string
	fmt.Print("Enter username and password (space separated): ")
	fmt.Fscanln(os.Stdin, &username, &password)

	return username, password
}
