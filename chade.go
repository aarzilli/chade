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

func decodeInput(bytes []byte) (map[int][]string, map[string]string) {
	r := make(map[int][]string)
	reasons := make(map[string]string)
	for _, decoder := range decoders {
		ok, char, reason := decoder.fn(bytes)
		if ok {
			r[char] = append(r[char], decoder.name)
		} else {
			reasons[decoder.name] = reason
		}
	}
	return r, reasons
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
	switch os.Args[1] {
	case "test-unidecode":
		testUnidecode()
		return
	case "test-jis":
		testJIS()
		return
	}
	
	InitUnicodeData()
	InitHTMLEntities()
	
	argument :=  strings.TrimSpace(strings.Join(os.Args[1:], " "))
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
		characters, reasons := decodeInput(bytes)
		for character, decoderNames := range characters {
			fmt.Printf("Decoded as %v:\n\n", decoderNames)
			runEncodersCL(character, "\t")
			fmt.Printf("\n")
		}
		for decoderName, reason := range reasons {
			fmt.Printf("Can not be decoded as %s because %s\n", decoderName, reason);
		}
	}
}