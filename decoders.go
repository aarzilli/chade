package main

type Decoder struct {
	name string
	fn func([]byte) (bool, int)
}

var decoders []Decoder

