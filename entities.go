package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

var entities map[int]string = make(map[int]string)
var entityLookup map[string]int = make(map[string]int)

func InitHTMLEntities() {
	file, err := os.Open("entities.txt", os.O_RDONLY, 0)
	must(err)
	defer file.Close()
	in := bufio.NewReader(file)

	for line, err := in.ReadString('\n'); err == nil; line, err = in.ReadString('\n') {
		line = strings.TrimSpace(line)
		split := strings.Split(line, "\t", 2)
		name := split[1]
		var codepoint int
		_, err := fmt.Sscanf(split[0][2:], "%x", &codepoint)
		must(err)
		entities[codepoint] = name
		entityLookup[name] = codepoint
	}
}

