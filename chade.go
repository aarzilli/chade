package main

import (
	"os"
	"strings"
	"fmt"
)

func must(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func interpretInput(argument string) (string, int, []byte) {
	for _, interpreter := range interpreters {
		ok, char, bytes := interpreter.fn(argument)
		if ok { return interpreter.name, char, bytes }
	}
	return "", -1, nil
}

func decodeInput(bytes []byte) map[string]int {
	r := make(map[string]int)
	for _, decoder := range decoders {
		ok, char := decoder.fn(bytes)
		if ok { r[decoder.name] = char }
	}
	return r
}

type EncodingResult struct {
	name string
	value string
}

func runEncoders(character int) []EncodingResult {
	r := make([]EncodingResult, 0)
	for _, encoder := range encoders {
		ok, value := encoder.fn(character)
		if ok { r = append(r, EncodingResult{ encoder.name, value }) }
	}
	return r
}

func runEncodersCL(character int, indent string) {
	encodings := runEncoders(character)
	for _, encoding := range encodings {
		fmt.Printf("%sEncoded as %s:\t%s\n", indent, encoding.name, encoding.value)
	}
}

func main() {
	argument :=  strings.Join(os.Args[1:], " ")
	fmt.Printf("Argument: [%s]\n", argument)

	name, character, bytes := interpretInput(argument)

	if name == "" {
		fmt.Printf("Could not understand input\n")
		return;
	}

	fmt.Printf("Interpreted as %s\n", name)

	if bytes == nil {
		fmt.Printf("\n")
		runEncodersCL(character, "")
	} else {
		characters := decodeInput(bytes)
		for decoderName, character := range characters {
			fmt.Printf("Decoded as %s:\n\n", decoderName)
			runEncodersCL(character, "\t")
			fmt.Printf("\n")
		}
	}
}