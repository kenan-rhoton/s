package main

import (
	"fmt"
	"os"
)

func main() {

	// Switch for different actions
	switch {

	case len(os.Args) == 1:
		// No Args: Identify to server and ask for all your messages [comm, crypt, db] - 6

	case os.Args[1] == "reg":
		// Register: Send your publickey to the server and request a username [comm, crypt, db] - 4

	case os.Args[1] == "add":
		// Add: Get a registered user's publickey to speed up stuff and or whatever [comm, file]

	case os.Args[1] == "as":
		// Alias: Define a personal alias for a user [file]

	default:
		// Default: Send your message to a certain username on the server [comm, (file)] - 5

	case os.Args[1] == "srv":
		// Serve: Act as a server [comm, db, crypt] - 1
		StartServer()

	case os.Args[1] == "sel":
		// Select: Choose a server [file] - 3

	case os.Args[1] == "gen":
		// Generate: generate a key [crypt, file] - 2
		key, err := GenerateKey(4096)
		if err != nil {
			fmt.Println("Error occurred while generating key: ", err.Error())
			break
		}
		err = SaveAs(key, "_privkey")
		if err != nil {
			fmt.Println("Error occurred while saving key: ", err.Error())
		}
		fmt.Println("Private Key successfully generated")
	}
}
