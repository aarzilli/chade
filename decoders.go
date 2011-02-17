package main

import (
	"iconv"
	"fmt"
)

type Decoder struct {
	name string
	fn func([]byte) (bool, int, string)
}

var decoders []Decoder = []Decoder{
	Decoder{ "ASCII", DecASCII },

	Decoder{ "UTF-8", DecUtf8 },

	// UTF-16LE
	// UTF-16BE
	
	Decoder{ "ISO-8859-1 (latin1)", MakeDecIconv("iso-8859-1") },
	Decoder{ "Windows-1251 (latin1 for windows)", MakeDecIconv("windows-1252") },
	Decoder{ "Windows-1256 (arab windows)", MakeDecIconv("windows-1256") },
	Decoder{ "ISO-8856-7 (greek)", MakeDecIconv("iso8859-7") },
	Decoder{ "Windows-1253 (greek windows)", MakeDecIconv("windows-1253") },
	Decoder{ "ISO-8859-8 (hebrew)", MakeDecIconv("iso-8859-8") },
	Decoder{ "Windows-1255", MakeDecIconv("windows-1255") },

	// Shift-JIS
	// ISO-2022-JP
	// EUC-JP
	// EUC-KR
	// ISO-2022-KR
	// KOI8-R (cyrillic)
	// Windows-1251 (russian windows)
	// Windows-874 (thai windows)
	// ISO-8859-11 (thai)
	// TIS-620 (thai)
	// Windows-1258 (vietnamese)
	// EUC-CN (chinese)
	// BIG5 (chinese)
	// GBK (chinese)
}

// byte -> uint8

func DecASCII(in []byte) (bool, int, string) {
	if len(in) > 1 { return false, -1, "Too many bytes" };
	if in[0] >= 128 { return false, -1, "MSB set" };
	return true, int(in[0]), ""
}

func MakeDecIconv(charset string) func([]byte) (bool, int, string) {
	return func(in []byte) (bool, int, string) {
		if len(in) != 1 { return false, -1, "More than one character" }
		out, err := iconv.Conv("UTF-8", charset, string(in))
		if err != nil { return false, -1, "Rejected by iconv" }
		if len(out) == 0 { return false, -1, "Rejected by iconv" }
		return true, []int(out)[0], ""
	}
}

func Utf8Char1Decode(in byte) (length int, code byte) {
	if in < 128 { return 1, in }
	if in < 192 { return -1, 0 }
	if in < 224 { return 2, in & 0x1F }
	if in < 240 { return 3, in & 0x0F }
	if in < 248 { return 4, in & 0x07 }
	if in < 252 { return 5, in & 0x03 }
	if in < 254 { return 6, in & 0x01 }
	return -1, 0
}

func AcceptUtf8SequenceByte(in byte) (bool, byte) {
	if in < 128 { return false, 0 }
	if in >= 192 { return false, 0 }
	return true, in & 0x3F
}

func DecUtf8(in []byte) (bool, int, string) {
	var length, acccode int
	
	length, subcode := Utf8Char1Decode(in[0])
	acccode = int(subcode)
	if length == -1 {
		return false, -1, "Invalid first byte of an utf8 sequence (FE or FF)"
	}
	if len(in) != length {
		return false, -1, fmt.Sprintf("First byte requires a sequence of %d characters but %d characters were provided", length, len(in))
	}

	for i, abyte := range in[1:len(in)] {
		if ok, subcode := AcceptUtf8SequenceByte(abyte); ok {
			acccode <<= 6
			acccode += int(subcode)
		} else {
			return false, -1, fmt.Sprintf("Character %d can not be part of an utf8 sequence", i)
		}
	}

	return true, acccode, ""
}