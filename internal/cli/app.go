package cli

import (
	"flag"
	"fmt"
	"strings"
)

func RunApp() {
	var stunIP string
	flag.StringVar(&stunIP, "stunip", "STUN:freestun.net:3479", "Sets the IP for the STUN server")

	flag.Parse()

	cmd := strings.ToLower(flag.Arg(0))
	argument := strings.ToLower(flag.Arg(1))
	if cmd == "send" || cmd == "s" {
		send(stunIP, argument)
	} else if cmd == "receive" || cmd == "r" {
		receive(argument)
	} else {
		fmt.Printf("Command \"%s\" not found.\n", cmd)
	}
}
