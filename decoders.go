package main

type Decoder struct {
	name string
	fn func([]byte) (bool, int)
}

var decoders []Decoder = []Decoder{
	Decoder{ "ASCII", DecASCII },
}

// byte -> uint8

func DecASCII(in []byte) (bool, int) {
	if (len(in) == 1) && (in[0] < 128) {
		return true, int(in[0])
	}
	return false, -1
}




