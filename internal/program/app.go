package program

import (
	"flag"
	"fmt"
	"strings"
)

func RunApp() {
	config := ReadConfig()

	flag.Parse()

	cmd := strings.ToLower(flag.Arg(0))
	filepath := flag.Arg(1)

	if cmd == "send" || cmd == "s" {
		send(config, filepath)
	} else if cmd == "receive" || cmd == "r" {
		receive(*config)
	} else if cmd == "" {
		fmt.Println("Please enter a command.")
	} else {
		fmt.Printf("Command \"%s\" not found.\n", cmd)
	}
}
