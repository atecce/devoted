package db

import (
	"testing"
)

func TestExample1(t *testing.T) {

	db := NewDatabase()

	var val *string
	var n uint

	val = db.Get("foo")
	if val != nil {
		t.Errorf("value should be nil, got %v", val)
	}

	db.Set("a", "foo")
	db.Set("b", "foo")
	n = db.Count("foo")
	if n != 2 {
		t.Errorf("count should be 2, got %v", n)
	}

	n = db.Count("bar")
	if n != 0 {
		t.Errorf("count should be 0, got %v", n)
	}

	db.Delete("a")

	n = db.Count("foo")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.Set("b", "baz")

	val = db.Get("B")
	if val != nil {
		t.Errorf("value should be nil, got %v", val)
	}
}

func TestExample2(t *testing.T) {

	db := NewDatabase()

	var val *string
	var n uint

	db.Set("a", "foo")
	db.Set("a", "foo")
	n = db.Count("foo")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	val = db.Get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.Delete("a")
	val = db.Get("a")
	if val != nil {
		t.Errorf("val should be nil, got %v", val)
	}

	n = db.Count("foo")
	if n != 0 {
		t.Errorf("count should be 0, got %v", n)
	}
}

func TestExample3(t *testing.T) {

	db := NewDatabase()

	var val *string

	db.Begin()

	db.Set("a", "foo")
	val = db.Get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.Begin()

	db.Set("a", "bar")
	val = db.Get("a")
	if val == nil || *val != "bar" {
		t.Errorf("val should be bar, got %v", val)
	}

	db.Rollback()
	val = db.Get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.Rollback()
	val = db.Get("a")
	if val != nil {
		t.Errorf("val should be nil, got %v", val)
	}
}

func TestExample4(t *testing.T) {

	db := NewDatabase()

	var val *string
	var n uint

	db.Set("a", "foo")
	db.Set("b", "baz")

	db.Begin()

	val = db.Get("a")
	if val == nil || *val != "foo" {
		t.Errorf("val should be foo, got %v", val)
	}

	db.Set("a", "bar")
	n = db.Count("bar")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.Begin()

	n = db.Count("bar")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.Delete("a")

	val = db.Get("a")
	if val != nil {
		t.Errorf("val should be nil, got %v", *val)
	}

	n = db.Count("bar")
	if n != 0 {
		t.Errorf(`db.count("bar") should be 0, got %v`, n)
	}

	db.Rollback()

	val = db.Get("a")
	if val == nil || *val != "bar" {
		t.Errorf("val should be bar, got %v", *val)
	}

	n = db.Count("bar")
	if n != 1 {
		t.Errorf("count should be 1, got %v", n)
	}

	db.Commit()

	val = db.Get("a")
	if val == nil || *val != "bar" {
		t.Errorf("val should be bar, got %v", *val)
	}

	val = db.Get("b")
	if val == nil || *val != "baz" {
		t.Errorf("val should be baz, got %v", *val)
	}
}
