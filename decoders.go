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
	Decoder{ "UTF-16LE", DecUtf16LE },
	Decoder{ "UTF-16BE", DecUtf16BE },
	
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
	if in & 0x80 == 0x00 { return 1, in & 0x7F }
	if in & 0xE0 == 0xC0 { return 2, in & 0x1F }
	if in & 0xF0 == 0xE0 { return 3, in & 0x0F }
	if in & 0xF8 == 0xF0 { return 4, in & 0x07 }
	if in & 0xFC == 0xF8 { return 5, in & 0x03 }
	if in & 0xFE == 0xFC { return 6, in & 0x01 }

	return -1, 0
}

func AcceptUtf8SequenceByte(in byte) (bool, byte) {
	return (in & 0xC0 == 0x80), in & 0x3F
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

func DecUtf16LE(in []byte) (bool, int, string) {
	if (len(in) != 2) && (len(in) != 4) {
		return false, -1, fmt.Sprintf("Unacceptable number of bytes for an UTF-16 character (can be 2 or 4 was %d)", len(in))
	}

	ints := make([]uint16, len(in)/2)

	for i := 0; i < len(in); i += 2 {
		ints[i/2] = uint16(in[i]) + (uint16(in[i+1]) << 8)
	}

	return DecUtf16Common(ints)
}

func DecUtf16BE(in []byte) (bool, int, string) {
	if (len(in) != 2) && (len(in) != 4) {
		return false, -1, fmt.Sprintf("Unacceptable number of bytes for an UTF-16 character (can be 2 or 4 was %d)", len(in))
	}

	ints := make([]uint16, len(in)/2)

	for i := 0; i < len(in); i += 2 {
		ints[i/2] = (uint16(in[i]) << 8) + uint16(in[i+1])
	}

	return DecUtf16Common(ints)
}

func DecUtf16Common(ints []uint16) (bool, int, string) {
	if len(ints) == 1 {
		return true, int(ints[0]), ""
	}

	if (ints[0] < 0xd800) || (ints[0] > 0xdbff) {
		return false, -1, fmt.Sprintf("First element of the pair is not a high surrogate (%x)", ints[0])
	}

	hisur := ints[0] - 0xd800

	if (ints[1] < 0xdc00) || (ints[1] > 0xdfff) {
		return false, -1, "Second element of the pair is not a low surrogate"
	}

	lowsur := ints[1] - 0xdc00

	result := (uint32(hisur) << 10) + uint32(lowsur) + 0x10000

	return true, int(result), ""
}