package main

import (
	"encoding/gob"
	"net"
)

type simpleMessage struct {
	Message string
}

func handleServeConnection(c net.Conn, out chan string) {
	gobber := gob.NewDecoder(c)
	res := &simpleMessage{}
	gobber.Decode(res)
	out <- res.Message
}

func Serve(out chan string) {
	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleServeConnection(conn, out)
	}
}

func Send(target, msg string) {
	conn, err := net.Dial("tcp", target+":8090")
	if err != nil {
		return
	}
	m := &simpleMessage{Message: msg}
	gobber := gob.NewEncoder(conn)
	gobber.Encode(m)
	//writer := bufio.NewWriter(conn)
	//writer.WriteString(msg)
	//writer.Flush()
	conn.Close()
}
