package main

import "strings"

type String string

func (v String) IndexOf(s string) int {
	return strings.Index(v.ToString(), s)
}

func (v String) Contains(s string) bool {
	return strings.Contains(v.ToString(), s)
}

func (v String) ToString() string {
	return string(v)
}
