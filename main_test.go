package main

import (
	"testing"
)

func Test(t *testing.T) {

	db := database{}
	db.store = make(map[string]string)
	db.setBuf = make(map[string]string)

	var val *string
	var n uint

	val = db.get("foo")
	if val != nil {
		t.Errorf("value should be nil, got %v", val)
	}

	db.set("a", "foo")
	db.set("b", "foo")
	n = db.count("foo")
	if n != 2 {
		t.Errorf("count should be 2, got %v", n)
	}

	n = db.count("bar")
	if n != 0 {
		t.Errorf("count should be 0, got %v", n)
	}

	db.delete("a")

	n = db.count("foo")
	if n != 1 {
		t.Errorf("count should be 0, got %v", n)
	}

	db.set("b", "baz")

	val = db.get("B")
	if val != nil {
		t.Errorf("value should be nil, got %v", val)
	}
}
