package main

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strings"
	"time"
)

var config struct {
	Target string
	ID     string
	Key    *rsa.PrivateKey
}

func main() {

	LoadFrom(&config, ".sconfig")

	// Switch for different actions
	switch {

	case len(os.Args) == 1:
		// No Args: Identify to server and ask for all your messages [comm, crypt, db] - 6
		if config.Target == "" {
			fmt.Println("first <sel> a server")
			break
		}
		if config.ID == "" {
			fmt.Println("first <reg> with the server")
			break
		}
		respHandle := SendAction(config.Target, "retrieve", "", config.ID)
		select {
		case <-time.After(time.Second * 10):
			fmt.Println("Server did not respond")
		case answer := <-respHandle:
			for _, msg := range answer.Arguments {
				m, err := Decode(msg, config.Key)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(m)
				}
			}
		}

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

		respHandle := SendAction(config.Target, "reg", Serialize(config.Key.PublicKey), os.Args[2])
		select {
		case <-time.After(time.Second * 10):
			fmt.Println("Server did not respond")
		case answer := <-respHandle:
			config.ID = os.Args[2]
			err := SaveAs(&config, ".sconfig")
			if err != nil {
				fmt.Println("Error occurred while saving config: ", err.Error())
				break
			}
			fmt.Println(answer.Action, answer.Message)
		}
	case os.Args[1] == "add":
		// Add: Get a registered user's publickey to speed up stuff and or whatever [comm, file]

	case os.Args[1] == "as":
		// Alias: Define a personal alias for a user [file]

	default:
		// Default: Send your message to a certain username on the server [comm, (file)] - 5
		if config.Target == "" {
			fmt.Println("first <sel> a server")
			break
		}
		if config.ID == "" {
			fmt.Println("first <reg> with the server")
			break
		}
		if len(os.Args) < 3 {
			fmt.Println("s <target> <message>")
			break
		}

		respHandle := SendAction(config.Target, "get", os.Args[1])

		select {
		case <-time.After(time.Second * 10):
			fmt.Println("Server did not respond")
		case answer := <-respHandle:
			if answer.Action != "" {
				fmt.Println(answer.Action, answer.Message)
				break
			} else {
				key := &rsa.PublicKey{}
				UnSerialize(key, answer.Message)
				msg := strings.Join(os.Args[2:], " ")
				enc, err := Encode(config.ID+" -> "+msg, key)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				SendAction(config.Target, "send", enc, os.Args[1])
			}
		}

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
		err := SaveAs(&config, ".sconfig")
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
		err = SaveAs(config, ".sconfig")
		if err != nil {
			fmt.Println("Error occurred while saving key: ", err.Error())
			break
		}
		fmt.Println(config)
		fmt.Println("Private Key successfully generated")
	}
}
