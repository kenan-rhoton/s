package conn

import (
	"encoding/gob"
	"net"
	"time"
)

type Request struct {
	Message    string
	Action     string
	Arguments  []string
	connHandle net.Conn
}

var serverhandlekiller = make(chan bool)

func (s *Request) FullAnswer(msg, action string, args ...string) {
	go func() {
		s.connHandle.SetDeadline(time.Now().Add(1000 * time.Millisecond))
		gobber := gob.NewEncoder(s.connHandle)
		send := &Request{Message: msg, Action: action, Arguments: args}
		gobber.Encode(send)
	}()
}
func (s *Request) Answer(msg string) {
	s.FullAnswer(msg, "")
}

func (s *Request) Close() {
	go func(m *Request) {
		time.Sleep(time.Second * 1)
		m.connHandle.Close()
	}(s)
}

func handleServeConnection(c net.Conn, out chan Request) {
	gobber := gob.NewDecoder(c)
	res := &Request{}
	gobber.Decode(res)
	res.connHandle = c
	out <- *res
}

func Serve(address string) chan Request {
	c := make(chan Request)
	go func(out chan Request) {
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

func SendText(target, msg string) chan Request {
	return SendAction(target, "", msg)
}

func SendAction(target, action, msg string, args ...string) chan Request {
	handler := make(chan Request)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return nil
	}
	m := &Request{Message: msg, Action: action, Arguments: args}
	gobber := gob.NewEncoder(conn)
	gobber.Encode(m)
	gobbler := gob.NewDecoder(conn)
	go func() {
		conn.SetDeadline(time.Now().Add(10000 * time.Millisecond))
		resp := &Request{}
		gobbler.Decode(resp)
		handler <- *resp
	}()
	return handler
}
