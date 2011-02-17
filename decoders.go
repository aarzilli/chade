package main

import (
	"iconv"
)

type Decoder struct {
	name string
	fn func([]byte) (bool, int, string)
}

var decoders []Decoder = []Decoder{
	Decoder{ "ASCII", DecASCII },

	// UTF-8
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
		out, err := iconv.Conv("UTF-8", charset, string(in))
		if err != nil { return false, -1, "Rejected by iconv" }
		if len(out) == 0 { return false, -1, "Rejected by iconv" }
		return true, []int(out)[0], ""
	}
}


