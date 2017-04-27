package main

import (
	"fmt"
	"github.com/kenan-rhoton/s/client"
	"github.com/kenan-rhoton/s/server"
	"os"
)

func main() {

	// Switch for different actions
	switch {

	case len(os.Args) == 1:
		client.FetchMessages()
	case os.Args[1] == "reg":
		if len(os.Args) < 3 {
			fmt.Println("s reg <name>")
			break
		}
		client.Register(os.Args[2])
	case os.Args[1] == "add":
		// Add: Get a registered user's publickey to speed up stuff and or whatever [comm, file]
		// TODO

	case os.Args[1] == "as":
		// Alias: Define a personal alias for a user [file]
		// TODO

	default:
		if len(os.Args) < 3 {
			fmt.Println("s <target> <message>")
			break
		}
		client.Send(os.Args[1], os.Args[2:])
	case os.Args[1] == "srv":
		// StartServer: Act as a server
		server.StartServer()

	case os.Args[1] == "sel":
		if len(os.Args) < 3 {
			fmt.Println("sel <ipaddress>")
			break
		}
		client.Select(os.Args[2])
	case os.Args[1] == "gen":
		client.Generate()
	}
}
