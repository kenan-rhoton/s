package main

import (
	"fmt"
)

func StartServer() {
	db, err := UseDatabase("s")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	stream := Serve(":8090")
	for {
		request := <-stream
		if request.Action == "" {
			continue
		}
		switch {
		case request.Action == "reg":
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
		case request.Action == "get":
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
		case request.Action == "send":
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
		case request.Action == "retrieve":
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
		request.Close()
	}

}
