package db

import (
	"testing"
)

type testID struct {
	Alias string
	PKey  []byte
}

type testMessage struct {
	target  string
	message []byte
}

func TestUpdateUser(t *testing.T) {
	db, err := UseDatabase("test")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close()
	err = db.Reset()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = db.SaveID("bo", []byte("?"))
	if err != nil {
		t.Errorf(err.Error())
	}
	err = db.SaveID("bo", []byte("!"))
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSaveIdentity(t *testing.T) {
	db, err := UseDatabase("test")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close()
	err = db.Reset()
	if err != nil {
		t.Errorf(err.Error())
	}
	testdata := []testID{
		{"boop", []byte("n12ui3h102h312undiqna0cpc+ÀW+")},
		{"bep", []byte("FA SFOA MFMWE QTM QRQ  ZX;c.-xc,af e q")},
	}
	for _, v := range testdata {
		err = db.SaveID(v.Alias, v.PKey)
		if err != nil {
			t.Errorf(err.Error())
		}
		res, err := db.GetKey(v.Alias)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !same(res, v.PKey) {
			t.Errorf("Wrong key: expected \"%s\" but got \"%s\"\n", v.PKey, res)
		}
	}
}

func same(one, two []byte) bool {
	if len(one) != len(two) {
		return false
	}
	for i := 0; i < len(one); i++ {
		if one[i] != two[i] {
			return false
		}
	}
	return true
}

func TestSaveMessage(t *testing.T) {
	db, err := UseDatabase("test")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close()
	err = db.Reset()
	if err != nil {
		t.Errorf(err.Error())
	}
	testdata := []testMessage{
		{"testuser", []byte("lol")},
		{"testuser", []byte("rofl")},
		{"testuser", []byte("duude")},
		{"testuser", []byte("monkeey")},
	}
	err = db.SaveID("testuser", []byte("whatever, really"))
	if err != nil {
		t.Errorf(err.Error())
	}
	messagelist := make([][]byte, 0)
	for _, v := range testdata {
		err = db.SaveMessage(v.target, v.message)
		if err != nil {
			t.Errorf(err.Error())
		}
		messagelist = append(messagelist, v.message)
	}
	msgs, err := db.GetMessages("testuser")
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, m := range messagelist {
		found := false
		for _, msg := range msgs {
			if same(m, msg) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Did not receive message %s", m)
		}
	}
}

func TestSaveMessageWithInvalidTarget(t *testing.T) {
	db, err := UseDatabase("test")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer db.Close()
	err = db.Reset()
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := "invalid target"
	err = db.SaveMessage("testuser", []byte("doesn't matter"))
	if err == nil {
		t.Errorf("Did not get an error with nonexistent user!")
	} else {
		if err.Error() != expected {
			t.Errorf("Got wrong error, expected %s, got %s", expected, err.Error())
		}
	}
}
