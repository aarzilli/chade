package main;

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"strconv"
)

type UnicodeData struct {
	Name string
	Block string
	GeneralCategory string
	CanonicalCombiningClass string
	BidiClass string
	DecompositionType string
	DecompositionMapping string
	NumericType string
	NumericValue string
	BidiMirrored string
	Unicode1Name string
	ISOComment string
	SimpleUppercaseMapping string
	SimpleLowercaseMapping string
	SimpleTitlecaseMapping string
}

func (ud *UnicodeData) String() string {
	r := ""
	
	r += "Name: " + ud.Name + "\n"
	r += "Block: " + ud.Block + "\n"
	if ud.GeneralCategory != "" { r += "General Category: " + ud.GeneralCategory + "\n" }
	if ud.CanonicalCombiningClass != "" { r += "Canonical Combining Class: " + ud.CanonicalCombiningClass + "\n" }
	if ud.BidiClass != "" { r += "Bidi Class: " + ud.BidiClass + "\n" }
	if ud.DecompositionType != "" { r += "Decomposition Type: " + ud.DecompositionType + "\n" }
	if ud.DecompositionMapping != "" { r += "Decomposition Mapping: " + ud.DecompositionMapping + "\n" }
	if ud.NumericType != "" { r += "Numeric Type: " + ud.NumericType + "\n" }
	if ud.NumericValue != "" { r += "Numeric Value: " + ud.NumericValue + "\n" }
	if ud.BidiMirrored != "" { r += "Bidi Mirrored: " + ud.BidiMirrored + "\n" }
	if ud.Unicode1Name != "" { r += "Unicode 1 Name: " + ud.Unicode1Name + "\n" }
	if ud.ISOComment != "" { r += "ISO Comment: " + ud.ISOComment + "\n" }
	if ud.SimpleUppercaseMapping != "" { r += "Simple Uppercase Mapping: " + ud.SimpleUppercaseMapping + "\n" }
	if ud.SimpleLowercaseMapping != "" { r += "Simple Lowercase Mapping: " + ud.SimpleLowercaseMapping + "\n" }
	if ud.SimpleTitlecaseMapping != "" { r += "Simple Titlecase Mapping: " + ud.SimpleTitlecaseMapping + "\n" }

	return r
}

func MakeFromUnicodeDataLine(line string) (int, *UnicodeData) {
	fields := strings.Split(strings.TrimSpace(line), ";", -1)
	var n int
	_, err := fmt.Sscanf(fields[0], "%x", &n)
	must(err)
	return n, &UnicodeData{
		fields[1],
		"No_Block",
		fields[2], // general category
		fields[3], // canonical combining class
		fields[4], // bidi class
		fields[5], // decomposition
		fields[6],
		fields[7], // numeric
		fields[8],
		fields[9], // bidi mirrored
		fields[10], // compat
		fields[11],
		fields[12], // case mappings
		fields[13],
		fields[14],
	}
}

var UnicodeDataFile [0x10FFFF]*UnicodeData

func InitUnicodeDataUnicodeData() {
	file, err := os.Open("UnicodeData.txt", os.O_RDONLY, 0)
	must(err)
	defer file.Close()
	in := bufio.NewReader(file)

	lastId := 0
	
	for line, err := in.ReadString('\n'); err == nil; line, err = in.ReadString('\n') {
		id, ud := MakeFromUnicodeDataLine(line)
		for skipped := lastId; skipped < id; skipped++ {
			UnicodeDataFile[skipped] = UnicodeDataFile[lastId]
		}
		UnicodeDataFile[id] = ud
		lastId = id
	}
}

func InitUnicodeDataBlocks() {
	file, err := os.Open("Blocks.txt", os.O_RDONLY, 0)
	must(err)
	defer file.Close()
	in := bufio.NewReader(file)

	for line, err := in.ReadString('\n'); err == nil; line, err = in.ReadString('\n') {
		line = strings.TrimSpace(line)
		if len(line) == 0 { continue }
		if line[0] == '#' { continue }
		
		split := strings.Split(line, ";", 2)
		if len(split) != 2 { continue }

		therange, block := split[0], split[1]

		block = strings.TrimSpace(block)

		split = strings.Split(therange, "..", 2)

		if len(split) != 2 { continue }

		start, _ := strconv.Atoi(split[0])
		end, _ := strconv.Atoi(split[1])

		for i := start; i <= end; i++ {
			UnicodeDataFile[i].Block = block
		}
	}
}

func InitUnicodeData() {
	InitUnicodeDataUnicodeData()
	InitUnicodeDataBlocks()
}