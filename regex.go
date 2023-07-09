package main

import (
	"fmt"
	"strings"
)

type RegEx interface {
	Derive(char rune) RegEx
	IsNullable() bool
}

type EmptySet struct{}

func (e EmptySet) Derive(char rune) RegEx {
	return e
}

func (e EmptySet) IsNullable() bool {
	return false
}

type EmptyString struct{}

func (e EmptyString) Derive(char rune) RegEx {
	return EmptySet{}
}

func (e EmptyString) IsNullable() bool {
	return true
}

type Singleton struct {
	char rune
}

func (s Singleton) Derive(char rune) RegEx {
	if s.char == char {
		return EmptyString{}
	} else {
		return EmptySet{}
	}
}

func (s Singleton) IsNullable() bool {
	return false
}

type Union struct {
	r1, r2 RegEx
}

func (u Union) Derive(char rune) RegEx {
	return Union{u.r1.Derive(char), u.r2.Derive(char)}
}

func (u Union) IsNullable() bool {
	return u.r1.IsNullable() || u.r2.IsNullable()
}

type Concatenation struct {
	r1, r2 RegEx
}

func (c Concatenation) Derive(char rune) RegEx {
	if c.r1.IsNullable() {
		return Union{Concatenation{c.r1.Derive(char), c.r2}, c.r2.Derive(char)}
	} else {
		return Concatenation{c.r1.Derive(char), c.r2}
	}
}

func (c Concatenation) IsNullable() bool {
	return c.r1.IsNullable() && c.r2.IsNullable()
}

type Star struct {
	r RegEx
}

func (s Star) Derive(char rune) RegEx {
	return Concatenation{s.r.Derive(char), Star{s.r}}
}

func (s Star) IsNullable() bool {
	return true
}

func parseRegex(pattern string) (RegEx, string) {
	return Singleton{rune(pattern[0])}, pattern[1:]
}

func matches(input string, pattern string) bool {
	r, _ := parseRegex(pattern)
	for _, char := range input {
		r = r.Derive(char)
	}
	return r.IsNullable()
}

func main() {
	fmt.Println(matches("abc", "abc"))
	fmt.Println(matches("abc", "def"))
	fmt.Println(matches("abc", "a.c"))
	fmt.Println(matches("abc", "a.*c"))
}