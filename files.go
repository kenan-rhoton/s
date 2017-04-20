package main

import (
	"encoding/gob"
	"os"
)

func SaveAs(data interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	gobber := gob.NewEncoder(file)
	gobber.Encode(data)
	return nil
}

func LoadFrom(data interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	gobber := gob.NewDecoder(file)
	gobber.Decode(data)
	return nil
}
