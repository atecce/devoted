package main

import (
	"testing"
)

func TestExample1(t *testing.T) {

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
		t.Errorf("count should be 1, got %v", n)
	}

	db.set("b", "baz")

	val = db.get("B")
	if val != nil {
		t.Errorf("value should be nil, got %v", val)
	}
}

func TestExample2(t *testing.T) {

	db := database{}
	db.store = make(map[string]string)
	db.setBuf = make(map[string]string)

	var val *string
	var n uint

	db.set("a", "foo")
	db.set("a", "foo")
	n = db.count("foo")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	val = db.get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.delete("a")
	val = db.get("a")
	if val != nil {
		t.Errorf("val should be nil, got %v", val)
	}

	n = db.count("foo")
	if n != 0 {
		t.Errorf("count should be 0, got %v", n)
	}
}

func TestExample3(t *testing.T) {

	db := database{}
	db.store = make(map[string]string)
	db.setBuf = make(map[string]string)

	var val *string
	// var n uint

	db.begin()

	db.set("a", "foo")
	val = db.get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.begin()

	db.set("a", "bar")
	val = db.get("a")
	if val == nil || *val != "bar" {
		t.Errorf("val should be bar, got %v", val)
	}

	db.rollback()
	val = db.get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.rollback()
	val = db.get("a")
	if val != nil {
		t.Errorf("val should be nil, got %v", val)
	}
}

func TestExample4(t *testing.T) {

	db := database{}
	db.store = make(map[string]string)
	db.setBuf = make(map[string]string)

	var val *string
	var n uint

	db.set("a", "foo")
	db.set("b", "baz")

	db.begin()

	val = db.get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.set("a", "bar")
	n = db.count("bar")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.begin()
	n = db.count("bar")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.delete("a")
	val = db.get("a")
	if val != nil {
		t.Errorf("val should be nil, got %v", *val)
	}

	n = db.count("bar")
	if n != 0 {
		t.Errorf("count should be 0, got %v", n)
	}

	db.rollback()
	val = db.get("a")
	if val == nil || *val != "bar" {
		t.Errorf("val should be bar, got %v", *val)
	}

	n = db.count("bar")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.commit()

	val = db.get("a")
	if val == nil || *val != "bar" {
		t.Errorf("val should be bar, got %v", *val)
	}

	val = db.get("b")
	if val == nil || *val != "baz" {
		t.Errorf("val should be baz, got %v", *val)
	}
}