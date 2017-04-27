package client

import (
	"crypto/rsa"
	"fmt"
	"github.com/kenan-rhoton/s/conn"
	"github.com/kenan-rhoton/s/crypt"
	"github.com/kenan-rhoton/s/files"
	"os"
	"strings"
	"time"
)

var config struct {
	Target string
	ID     string
	Key    *rsa.PrivateKey
}

func init() {
	files.LoadFrom(&config, ".sconfig")
}

// FetchMessages: Identify to server and ask for all your messages
func FetchMessages() {
	if config.Target == "" {
		fmt.Println("first <sel> a server")
		return
	}
	if config.ID == "" {
		fmt.Println("first <reg> with the server")
		return
	}
	respHandle := conn.SendAction(config.Target, "retrieve", "", config.ID)
	select {
	case <-time.After(time.Second * 10):
		fmt.Println("Server did not respond")
	case answer := <-respHandle:
		for _, msg := range answer.Arguments {
			m, err := crypt.Decode(msg, config.Key)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(m)
			}
		}
	}

}

// Register: Send your publickey to the server and request a username
func Register(name string) {
	if config.Target == "" {
		fmt.Println("first <sel> a server")
		return
	}
	if len(os.Args) < 3 {
		fmt.Println("reg <alias>")
		return
	}

	respHandle := conn.SendAction(
		config.Target,
		"reg",
		Serialize(config.Key.PublicKey),
		name)
	select {
	case <-time.After(time.Second * 10):
		fmt.Println("Server did not respond")
	case answer := <-respHandle:
		config.ID = name
		err := files.SaveAs(&config, ".sconfig")
		if err != nil {
			fmt.Println(
				"Error occurred while saving config: ",
				err.Error())
			return
		}
		fmt.Println(answer.Action, answer.Message)
	}
}

// Send: Send your message to a certain username on the server
func Send(target string, msg []string) {
	if config.Target == "" {
		fmt.Println("first <sel> a server")
		return
	}
	if config.ID == "" {
		fmt.Println("first <reg> with the server")
		return
	}

	respHandle := conn.SendAction(config.Target, "get", target)

	select {
	case <-time.After(time.Second * 10):
		fmt.Println("Server did not respond")
	case answer := <-respHandle:
		if answer.Action != "" {
			fmt.Println(answer.Action, answer.Message)
			return
		} else {
			key := &rsa.PublicKey{}
			UnSerialize(key, answer.Message)
			m := strings.Join(msg, " ")
			enc, err := crypt.Encode(config.ID+" -> "+m, key)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			conn.SendAction(config.Target, "send", enc, target)
		}
	}

}

// Select: Choose a server [file] - 3
func Select(target string) {
	config.Target = target
	err := files.SaveAs(&config, ".sconfig")
	if err != nil {
		fmt.Println("Error occurred while saving config: ", err.Error())
		return
	}
	fmt.Println("Selected new server: ", target)
}

// Generate: generate a key [crypt, file] - 2
func Generate() {
	var err error
	config.Key, err = crypt.GenerateKey(4096)
	if err != nil {
		fmt.Println("Error occurred while generating key: ", err.Error())
		return
	}
	err = files.SaveAs(config, ".sconfig")
	if err != nil {
		fmt.Println("Error occurred while saving key: ", err.Error())
		return
	}
	fmt.Println(config)
	fmt.Println("Private Key successfully generated")
}
