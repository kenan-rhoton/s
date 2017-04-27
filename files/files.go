package files

import (
	"encoding/gob"
	"os"
)

func SaveAs(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	gobber := gob.NewEncoder(file)
	err = gobber.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func LoadFrom(data interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	gobber := gob.NewDecoder(file)
	err = gobber.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
