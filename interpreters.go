package main

import (
	"regexp"
	"fmt"
	"os"
	"strings"
)

type Interpreter struct {
	name string
	fn func(string) (bool, int, []byte);
}

var interpreters []Interpreter = []Interpreter{
	Interpreter{ "Character", IntCharacter },
	Interpreter{ "Java literal", IntJava },
	Interpreter{ "HTML decimal character reference", IntHTMLDec },
	Interpreter{ "HTML hexadecimal character reference", IntHTMLHex },
	Interpreter{ "HTML entity", IntHTMLEntity },
	Interpreter{ "Bytes", IntBytes },
}

func IntCharacter(arg string) (bool, int, []byte) {
	varg := []int(arg)
	if len(varg) == 1 {
		return true, varg[0], nil
	}
	return false, -1, nil
}

func interpretCodepoint(arg string, hex bool) (bool, int, []byte) {
	var num int
	var err os.Error
	if hex {
		_, err = fmt.Sscanf(arg, "%x", &num)
	} else {
		_, err = fmt.Sscanf(arg, "%d", &num)
	}
	if err != nil { return false, -1, nil }
	if num < 0 { return false, -1 , nil }
	if num > 0x10ffff { return false, -1, nil }
	return true, num, nil
}

var javaRE *regexp.Regexp = regexp.MustCompile("^\\\\[uU][0-9a-fA-F]+$")

func IntJava(arg string) (bool, int, []byte) {
	if !javaRE.MatchString(arg) { return false, -1, nil }
	return interpretCodepoint(arg[2:], true)
}

var HTMLDecRE *regexp.Regexp = regexp.MustCompile("^&[0-9]+;$")

func IntHTMLDec(arg string) (bool, int, []byte) {
	if !HTMLDecRE.MatchString(arg) { return false, -1, nil }
	return interpretCodepoint(arg[1:len(arg)-1], false)
}

var HTMLHexRE *regexp.Regexp = regexp.MustCompile("^&#[0-9a-fA-F]+;$")

func IntHTMLHex(arg string) (bool, int, []byte) {
	if !HTMLHexRE.MatchString(arg) { return false, -1, nil }
	return interpretCodepoint(arg[2:len(arg)-1], true)
}

var HTMLEntityRE *regexp.Regexp = regexp.MustCompile("^&[a-fA-F]+;$")

func IntHTMLEntity(arg string) (bool, int, []byte) {
	if !HTMLEntityRE.MatchString(arg) { return false, -1, nil }
	name := arg[1:len(arg)-1]
	codepoint, ok := entityLookup[name]
	if !ok { return false, -1, nil }
	return true, codepoint, nil
}

var BytesRE *regexp.Regexp = regexp.MustCompile("^[0-9a-fA-F ]+$")

func IntBytes(arg string) (bool, int, []byte) {
	if !BytesRE.MatchString(arg) { return false, -1, nil }

	r := make([]byte, 0)
	
	for _, segment := range strings.Split(arg, " ", -1) {
		var n byte
		
		switch len(segment) {
		case 2:
			_, err := fmt.Sscanf(segment, "%x", &n)
			if err != nil { return false, -1, nil }
			r = append(r, n)
		case 4:
			_, err := fmt.Sscanf(segment[0:2], "%x", &n)
			if err != nil { return false, -1, nil }
			r = append(r, n)
			_, err = fmt.Sscanf(segment[2:], "%x", &n)
			if err != nil { return false, -1, nil }
			r = append(r, n)
		default:
			return false, -1, nil
		} 
	}
	
	return true, -1, r
}

