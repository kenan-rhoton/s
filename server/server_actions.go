package server

import (
	"fmt"
	"github.com/kenan-rhoton/s/conn"
	//"github.com/kenan-rhoton/s/crypt"
	database "github.com/kenan-rhoton/s/db"
	//"github.com/kenan-rhoton/s/files"
)

// User Registration Request
func Register(request conn.Request, db *database.DB) {
	if len(request.Arguments) > 0 && request.Arguments[0] != "" {
		_, err := db.GetKey(request.Arguments[0])
		if err != nil && err.Error() == "invalid target" {
			db.SaveID(request.Arguments[0], []byte(request.Message))
			request.Answer("Success")
		} else {
			request.FullAnswer("User exists", "error")
		}
	} else {
		request.FullAnswer("Missing argument", "error")
	}
}

// Request a User's Public Key
func GetUser(request conn.Request, db *database.DB) {
	if request.Message != "" {
		key, err := db.GetKey(request.Message)
		if err != nil {
			request.FullAnswer("User does not exist", "error")
		} else {
			request.Answer(string(key))
		}
	} else {
		request.FullAnswer("Missing argument", "error")
	}
}

// Request sending a message to a user
func SendMessage(request conn.Request, db *database.DB) {
	if len(request.Arguments) > 0 && request.Arguments[0] != "" {
		_, err := db.GetKey(request.Arguments[0])
		if err != nil && err.Error() == "invalid target" {
			request.FullAnswer("User does not exist", "error")
		} else {
			db.SaveMessage(request.Arguments[0], []byte(request.Message))
			request.Answer("Success")
		}
	} else {
		request.FullAnswer("Missing argument", "error")
	}
}

// Request retrieving all messages that match the client's public key
// TODO: Force the client to prove the key is his to improve security
func RetrieveMessages(request conn.Request, db *database.DB) {
	if len(request.Arguments) > 0 && request.Arguments[0] != "" {
		_, err := db.GetKey(request.Arguments[0])
		if err != nil && err.Error() == "invalid target" {
			request.FullAnswer("User does not exist", "error")
		} else {
			msgs, err := db.GetMessages(request.Arguments[0])
			if err != nil {
				request.FullAnswer("Could not retrieve messages", "error")
			} else {
				response := make([]string, 0)
				for _, m := range msgs {
					response = append(response, string(m))
				}
				request.FullAnswer("Success", "", response...)
			}
		}
	} else {
		request.FullAnswer("Missing argument", "error")
	}
}

// Launch a Server on port 8090
func StartServer() {
	db, err := database.UseDatabase("s")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	stream := conn.Serve(":8090")
	for {
		request := <-stream
		if request.Action == "" {
			continue
		}
		switch {
		case request.Action == "reg":
			Register(request, db)
		case request.Action == "get":
			GetUser(request, db)
		case request.Action == "send":
			SendMessage(request, db)
		case request.Action == "retrieve":
			RetrieveMessages(request, db)
		}
		request.Close()
	}

}
