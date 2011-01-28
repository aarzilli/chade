package main

type Interpreter struct {
	name string
	fn func(string) (bool, int, []byte);
}

var interpreters []Interpreter = []Interpreter{
	Interpreter{ "character", IntCharacter },
}

func IntCharacter (arg string) (bool, int, []byte) {
	varg := []int(arg)
	if len(varg) == 1 {
		return true, varg[0], nil
	}
	return false, -1, nil
}

