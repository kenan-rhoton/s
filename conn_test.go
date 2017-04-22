package main

import (
	"testing"
	"time"
)

func TestClientServerCommunication(t *testing.T) {
	testdata := []string{
		"Potatoes",
		"Burguers",
		"I like 'em all!",
		"Ché mirá!",
		"€å∫€µ√∑Ω",
		"This \n probably fails\n\n",
	}
	res := Serve(":8090")
	for _, v := range testdata {
		SendText("localhost:8090", v)
		ret := <-res
		if ret.Message != v {
			t.Errorf("Expected \"%s\" but got \"%s\"", v, ret.Message)
		}
		ret.Close()
	}
	KillServer()
}

func TestClientServerActions(t *testing.T) {
	testdata := [][]string{
		{"Ché mirá!", "add"},
		{"€å∫€µ√∑Ω", "del"},
		{"This \n probably fails\n\n", "funky"},
	}
	res := Serve(":8090")
	for _, v := range testdata {
		SendAction("localhost:8090", v[1], v[0])
		ret := <-res
		if ret.Message != v[0] {
			t.Errorf("Expected \"%s\" but got \"%s\"", v[0], ret.Message)
		}
		if ret.Action != v[1] {
			t.Errorf("Expected \"%s\" but got \"%s\"", v[1], ret.Action)
		}
		ret.Close()
	}
	KillServer()
}

func TestClientServerResponse(t *testing.T) {
	res := Serve(":8090")
	go func() {
		ret := <-res
		if ret.Message != "ping" {
			t.Errorf("Expected \"%s\" but got \"%s\"", "ping", ret.Message)
		}
		ret.Answer("pong")
		ret.Close()
	}()
	handler := SendText("localhost:8090", "ping")
	select {
	case <-time.After(time.Second * 1):
		t.Errorf("no answer")
		t.FailNow()
	case answer := <-handler:
		if answer.Message != "pong" {
			t.Errorf("Expected \"%s\" but got \"%s\"", "pong", answer.Message)
		}
	}

	KillServer()
}

func TestClientServerResponseVerbose(t *testing.T) {
	res := Serve(":8090")
	go func() {
		ret := <-res
		if ret.Message != "ping" {
			t.Errorf("Expected \"%s\" but got \"%s\"", "ping", ret.Message)
		}
		ret.FullAnswer("pong", "pang", "peng")
		time.Sleep(time.Second * 1)
		ret.Close()
	}()
	handler := SendText("localhost:8090", "ping")
	select {
	case <-time.After(time.Second * 1):
		t.Errorf("no answer")
		t.FailNow()
	case answer := <-handler:
		if answer.Message != "pong" {
			t.Errorf("Expected \"%s\" but got \"%s\"", "pong", answer.Message)
		}
		if answer.Action != "pang" {
			t.Errorf("Expected \"%s\" but got \"%s\"", "pang", answer.Message)
		}
		if len(answer.Arguments) < 1 || answer.Arguments[0] != "peng" {
			t.Errorf("Expected \"%s\" but got \"%v\"", "pang", answer.Arguments)
		}
	}

	KillServer()
}
