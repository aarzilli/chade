package main;

import (
	"fmt"
	"iconv"
)

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
