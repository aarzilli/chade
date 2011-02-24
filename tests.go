package main;

import (
	"fmt"
	"iconv"
)

func testJIS() {
	count := 0
	
	for i := 0; i <= 0x10FFFF; i++ {
		if (i & 0xFFFF) == 0 {
			fmt.Printf("Examining %x\n", i)
		}

		// workaround for bug in glibc iconv implementation of shift_jis encoder
		if i == 0x5c { continue }
		if i == 0x7e { continue }

		if shiftJISStr, err := iconv.Conv("shift_jis", "UTF-8", string(i)); err == nil {
			count++
			ok, out, reason := ShiftJISDecoder([]byte(shiftJISStr))
			if !ok {
				panic(fmt.Sprintf("Error decoding encoded shift jis character at codepoint %d: %s", i, reason))
			}
			if out != i {
				panic(fmt.Sprintf("Decoding mismatch for character at codepoint %x, returned %x", i, out))
			}
		}
	}

	fmt.Printf("Examined %d characters\n")
}

func testUnidecode() {
	for i := 0; i <= 0x10FFFF; i++ {
		if (i & 0xFFFF) == 0 {
			fmt.Printf("Examining %x\n", i)
		}
		
		if okUtf8, redec, explanation := DecUtf8([]byte(string(i))); !okUtf8 {
			panic(fmt.Sprintf("Decoding of utf-8 encoded %d failed with reason: %s", i, explanation))
		} else {
			if redec != i {
				panic(fmt.Sprintf("Decoding of utf-8 encoded %d erroneous, returned %d", i, redec))
			}
		}

		if (i >= 0xd800) && (i <= 0xdfff) { continue }

		if utf16LEStr, err := iconv.Conv("UTF-16LE", "UTF-8", string(i)); err != nil {
			panic(fmt.Sprintf("Iconv error at %x: %s", i, err))
		} else {
			if okUtf16, redec, explanation := DecUtf16LE([]byte(utf16LEStr)); !okUtf16 {
				panic(fmt.Sprintf("Decoding of utf-16le encoded %d failed with reason: %s", i, explanation))
			} else {
				if redec != i {
					panic(fmt.Sprintf("Decoding of utf-16le encoded %d erroneous, returned %d", i, redec))
				}
			}
		}
		

		if utf16BEStr, err := iconv.Conv("UTF-16BE", "UTF-8", string(i)); err != nil {
			panic(fmt.Sprintf("Iconv error at %x: %s", i, err))
		} else {
			if okUtf16, redec, explanation := DecUtf16BE([]byte(utf16BEStr)); !okUtf16 {
				panic(fmt.Sprintf("Decoding of utf-16be encoded %d failed with reason: %s", i, explanation))
			} else {
				if redec != i {
					panic(fmt.Sprintf("Decoding of utf-16be encoded %d erroneous, returned %d", i, redec))
				}
			}
		}
	}
}
