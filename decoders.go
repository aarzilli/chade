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
	Decoder{ "KOI8-R", MakeDecIconv("koi8-r") },
	Decoder{ "Windows-1251 (russian windows)", MakeDecIconv("windows-1251") },
	Decoder{ "Windows-874 (thai windows)", MakeDecIconv("windows-874") },
	Decoder{ "ISO-8859-11 (thai)", MakeDecIconv("iso-8859-11") },
	Decoder{ "TIS-620 (thai)", MakeDecIconv("tis-620") },
	Decoder{ "Windows-1258 (vietnamese)", MakeDecIconv("windows-1258") },

	Decoder{ "Shift-JIS", ShiftJISDecoder },
	

	// ISO-2022-JS
	// EUC-JP
	// EUC-KR
	// ISO-2022-KR
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

func IconvDecoder(in []byte, charset string) (bool, int, string) {
	out, err := iconv.Conv("UTF-8", charset, string(in))
	if err != nil { return false, -1, "Rejected by iconv" }
	if len(out) == 0 { return false, -1, "Rejected by iconv" }
	return true, []int(out)[0], ""
}

func MakeDecIconv(charset string) func([]byte) (bool, int, string) {
	return func(in []byte) (bool, int, string) {
		if len(in) != 1 { return false, -1, "More than one character" }
		return IconvDecoder(in, charset)
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

const (
	SINGLE_BYTE int = iota
	FORBIDDEN_FIRST_BYTE
	FIRST_BYTE
	NONSTANDARD_FIRST_BYTE
)

const (
	FORBIDDEN_SECOND_BYTE int = iota
	SECOND_BYTE_EVEN
	SECOND_BYTE_ODD
)

func ClassifyShiftJISByte1(in byte) int {
	switch in & 0xF0 {
	case 0x00, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70: return SINGLE_BYTE
	case 0x80:
		if in == 0x80 { return FORBIDDEN_FIRST_BYTE }
		return FIRST_BYTE
	case 0x90: return FIRST_BYTE
	case 0xA0, 0xB0, 0xC0, 0xD0:
		if in == 0xA0 { return FORBIDDEN_FIRST_BYTE }
		return SINGLE_BYTE
	case 0xE0: return FIRST_BYTE
	case 0xF0:
		if (in >= 0xF3) && (in <= 0xF9) { return NONSTANDARD_FIRST_BYTE }
		return FORBIDDEN_FIRST_BYTE
	}

	return -1
}

func ClassifyShiftJISByte2(in byte) int {
	switch in & 0xF0 {
	case 0x00, 0x10, 0x20, 0x30: return FORBIDDEN_SECOND_BYTE
	case 0x40, 0x50, 0x60, 0x70, 0x80, 0x90:
		if in == 0x7F { return FORBIDDEN_SECOND_BYTE }
		if in == 0x9F { return SECOND_BYTE_EVEN }
		return SECOND_BYTE_ODD
	case 0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0:
		if in <= 0xFD {
			return SECOND_BYTE_EVEN }
		return FORBIDDEN_SECOND_BYTE
	}

	return FORBIDDEN_SECOND_BYTE
}



func ShiftJISCheckByte2(in []byte) string {
	if len(in) > 2 { return "Too many bytes (never more than 2 bytes in a Shift-JIS character)" }
	if len(in) < 2 { return fmt.Sprintf("Not enought bytes for Shift-JIS sequence starting with: %x", in[0]) }
	switch ClassifyShiftJISByte2(in[1]) {
	case FORBIDDEN_SECOND_BYTE:
		return fmt.Sprintf("Unacceptable second byte: %x", in[1]);
	case SECOND_BYTE_ODD:
		fallthrough
	case SECOND_BYTE_EVEN:
		return ""
	}

	return ""
}


func ShiftJISDecoder(in []byte) (bool, int, string) {
	switch ClassifyShiftJISByte1(in[0]) {
	case SINGLE_BYTE:
		if len(in) > 1 { return false, -1, "Too many bytes, the first byte indicates only one is needed" }
		return IconvDecoder(in, "shift_jis")
	case FORBIDDEN_FIRST_BYTE:
		return false, -1, "First (only?) byte is forbidden in Shift-JIS"
	case FIRST_BYTE:
		if r := ShiftJISCheckByte2(in); r != "" {
			return false, -1, r
		}
		return IconvDecoder(in, "shift_jis")
	case NONSTANDARD_FIRST_BYTE:
		if r := ShiftJISCheckByte2(in); r != "" {
			return false, -1, r
		}
		return false, -1, "Non-standard first byte used (but everything else is ok, so this could be an emoji)"
	}

	return false, -1, fmt.Sprintf("This shouldn't happen %d", ClassifyShiftJISByte1(in[0]))
}