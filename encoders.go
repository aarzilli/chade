package main

import (
	"iconv"
	"fmt"
	"strings"
	"unicode"
)

type Encoder struct {
	name string
	fn func(int) (bool, string)
}

var encoders []Encoder = []Encoder{
	Encoder{ "Character", EncCharacter },
	
	Encoder{ "Codepoint", EncCodepoint },
	Encoder{ "Unicode Informations", EncUnicodeInfo },

	Encoder{ "Java String Literal", EncJava },
	Encoder{ "HTML Entity", EncHTML },
	
	Encoder{ "ASCII", EncASCII },
	
	Encoder{ "UTF-8", EncUtf8 },
	Encoder{ "UTF-16LE", MakeEncIconv("UTF-16LE", false) },
	Encoder{ "UTF-16BE", MakeEncIconv("UTF-16BE", false) },
	
	Encoder{ "ISO-8859-1 (latin1)", MakeEncIconv("iso-8859-1", true) },
	Encoder{ "Windows-1252 (latin1 for windows)", MakeEncIconv("windows-1252", true) },
	
	Encoder{ "Windows-1256 (arab windows)", MakeEncIconv("windows-1256", true) },
	Encoder{ "ISO-8859-7 (greek)", MakeEncIconv("iso-8859-7", true) },
	Encoder{ "Windows-1253 (greek windows)", MakeEncIconv("windows-1253", true) },
	Encoder{ "ISO-8859-8 (hebrew)", MakeEncIconv("iso-8859-8", true) },
	Encoder{ "Windows-1255 (hebrew windows)", MakeEncIconv("windows-1255", true) },
	Encoder{ "Shift-JIS", MakeEncIconv("shift_jis", true) },
	Encoder{ "ISO-2022-JP", MakeEncIconv("iso-2022-jp", true) },
	Encoder{ "EUC-JP", MakeEncIconv("euc-jp", true) },
	Encoder{ "EUC-KR", MakeEncIconv("euc-kr", true) },
	Encoder{ "ISO-2022-KR", MakeEncIconv("iso-2022-kr", true) },
	Encoder{ "KOI8-R (cyrillic)", MakeEncIconv("koi8-r", true) },
	Encoder{ "Windows-1251 (russian windows)", MakeEncIconv("windows-1251", true) },
	Encoder{ "Windows-874 (thai windows)", MakeEncIconv("windows-874", true) },
	Encoder{ "ISO-8859-11 (thai)", MakeEncIconv("iso-8859-11", true) },
	Encoder{ "TIS-620 (thai)", MakeEncIconv("tis-620", true) },
	Encoder{ "Windows-1258 (vietnamese)", MakeEncIconv("windows-1258", true) },
	Encoder{ "EUC-CN (chinese)", MakeEncIconv("euc-cn", true) },
	Encoder{ "BIG5 (chinese)", MakeEncIconv("big5", true) },
	Encoder{ "GBK (chinese)", MakeEncIconv("gbk", true) },
}

func encodeBytes(s string) string {
	r := make([]string, 0)
	r = append(r, "(hex)")
	for i := 0; i < len(s); i++ {
		r = append(r, fmt.Sprintf("%02X", s[i]))
	}
	return strings.Join(r, " ")
}

func EncCharacter(char int) (bool, string) {

	if unicode.Is(unicode.Cc, char) || unicode.Is(unicode.Cf, char) || unicode.Is(unicode.Co, char) || unicode.Is(unicode.Cs, char) || unicode.Is(unicode.Zl, char) || unicode.Is(unicode.Zp, char) || unicode.Is(unicode.Zs, char) { return false, "" }
	s := string(char)

	return true, s
}

func EncCodepoint(char int) (bool, string) {
	return true, fmt.Sprintf("%X (decimal: %d)", char, char)
}

func EncUnicodeInfo(char int) (bool, string) {
	return true, "\n"+UnicodeDataFile[char].String()
}

func EncUtf8(char int) (bool, string) {
	s := string(char)
	return true, encodeBytes(s)
}

func EncASCII(char int) (bool, string) {
	if char < 128 {
		return true, encodeBytes(string(char))
	}
	return false, ""
}

func EncJava(char int) (bool, string) {
	return true, fmt.Sprintf("\\u%04X", char)
}

func MakeEncIconv(charset string, excludeAscii bool) func(char int) (bool, string) {
	return func(char int) (bool, string) {
		s := string(char)
		if excludeAscii && (char < 128) { return false, "" }
		out, err := iconv.Conv(charset, "UTF-8", s)
		if err != nil { return false, "" }
		if len(out) == 0 { return false, "" }
		return true, encodeBytes(out)
	}
}

func EncHTML(char int) (bool, string) {
	symb, ok := entities[char]
	symbStr := ""
	if ok {
		symbStr = fmt.Sprintf(" entity: &%s;", symb)
	}
	return true, fmt.Sprintf("decimal: &#%d; hexadecimal: &#%X;%s", char, char, symbStr)
}