package main

import (
	"testing"
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
	res := make(chan string)
	go Serve(res)
	for _, v := range testdata {
		Send("localhost", v)
		ret := <-res
		if ret != v {
			t.Errorf("Expected \"%s\" but got \"%s\"", v, ret)
		}
	}
}
