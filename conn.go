package main

import (
	"encoding/gob"
	"net"
	"time"
)

type simpleMessage struct {
	Message    string
	Action     string
	Arguments  []string
	connHandle net.Conn
}

var serverhandlekiller = make(chan bool)

func (s *simpleMessage) FullAnswer(msg, action string, args ...string) {
	go func() {
		s.connHandle.SetDeadline(time.Now().Add(1000 * time.Millisecond))
		gobber := gob.NewEncoder(s.connHandle)
		send := &simpleMessage{Message: msg, Action: action, Arguments: args}
		gobber.Encode(send)
	}()
}
func (s *simpleMessage) Answer(msg string) {
	s.FullAnswer(msg, "")
}

func (s *simpleMessage) Close() {
	s.connHandle.Close()
}

func handleServeConnection(c net.Conn, out chan simpleMessage) {
	gobber := gob.NewDecoder(c)
	res := &simpleMessage{}
	gobber.Decode(res)
	res.connHandle = c
	out <- *res
}

func Serve(address string) chan simpleMessage {
	c := make(chan simpleMessage)
	go func(out chan simpleMessage) {
		tcp, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return
		}
		ln, err := net.ListenTCP("tcp", tcp)
		if err != nil {
			return
		}
		for {
			select {
			case _ = <-serverhandlekiller:
				ln.Close()
				return
			default:
			}
			ln.SetDeadline(time.Now().Add(10 * time.Millisecond))
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go handleServeConnection(conn, out)
		}
	}(c)

	return c
}

func KillServer() {
	serverhandlekiller <- true
}

func SendText(target, msg string) chan simpleMessage {
	return SendAction(target, "", msg)
}

func SendAction(target, action, msg string, args ...string) chan simpleMessage {
	handler := make(chan simpleMessage)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return nil
	}
	m := &simpleMessage{Message: msg, Action: action, Arguments: args}
	gobber := gob.NewEncoder(conn)
	gobber.Encode(m)
	gobbler := gob.NewDecoder(conn)
	go func() {
		conn.SetDeadline(time.Now().Add(10000 * time.Millisecond))
		resp := &simpleMessage{}
		gobbler.Decode(resp)
		handler <- *resp
	}()
	return handler
}
