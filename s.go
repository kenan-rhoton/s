package main

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"
)

var config struct {
	Target string
	ID     string
	Key    *rsa.PrivateKey
}

func main() {

	LoadFrom(&config, "_config")

	// Switch for different actions
	switch {

	case len(os.Args) == 1:
		// No Args: Identify to server and ask for all your messages [comm, crypt, db] - 6

	case os.Args[1] == "reg":
		// Register: Send your publickey to the server and request a username [comm, crypt, db] - 4
		if config.Target == "" {
			fmt.Println("first <sel> a server")
			break
		}
		if len(os.Args) < 3 {
			fmt.Println("reg <alias>")
			break
		}

		respHandle := SendAction(config.Target, "reg", Serialize(config.Key.PublicKey))
		select {
		case <-time.After(time.Second * 10):
			fmt.Println("Server did not respond")
		case answer := <-respHandle:
			fmt.Println(answer.Action, answer.Message)
		}
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
		if len(os.Args) < 3 {
			fmt.Println("sel <ipaddress>")
			break
		}
		config.Target = os.Args[2]
		err := SaveAs(&config, "_config")
		if err != nil {
			fmt.Println("Error occurred while saving config: ", err.Error())
			break
		}
		fmt.Println("Selected new server: ", os.Args[2])
	case os.Args[1] == "gen":
		// Generate: generate a key [crypt, file] - 2
		var err error
		config.Key, err = GenerateKey(4096)
		if err != nil {
			fmt.Println("Error occurred while generating key: ", err.Error())
			break
		}
		err = SaveAs(&config, "_config")
		if err != nil {
			fmt.Println("Error occurred while saving key: ", err.Error())
			break
		}
		fmt.Println("Private Key successfully generated")
	}
}
