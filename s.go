package main

import (
	"os"
)

func main() {
	// Switch for different actions
	switch {
	// No Args: Identify to server and ask for all your messages [comm, crypt, db] - 6
	case len(os.Args) == 1:
		// Register: Send your publickey to the server and request a username [comm, crypt, db] - 4
	case os.Args[1] == "reg":
		// Add: Get a registered user's publickey to speed up stuff and or whatever [comm, file]
	case os.Args[1] == "add":
		// Alias: Define a personal alias for a user [file]
	case os.Args[1] == "as":
		// Default: Send your message to a certain username on the server [comm, (file)] - 5
	default:
		// Serve: Act as a server [comm, db, crypt] - 1
	case os.Args[1] == "srv":
		// Select: Choose a server [file] - 3
	case os.Args[1] == "sel":
		// Generate: generate a key [crypt, file] - 2
	case os.Args[1] == "gen":
	}
}
