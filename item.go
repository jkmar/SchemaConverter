package main

type Item interface {
	Type() string
	Parse(string, map[interface{}]interface{})
}
