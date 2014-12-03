package main

import (
	"fmt"
	"io/ioutil"
	"github.com/wwkeyboard/rubyparse"
)
func main() {
	file := LoadFile("test.rb")
	lexer := rubyparse.Lex(file)

	for item := range lexer.Items {
		fmt.Printf("%v, %v\n", item.Val, item.Typ)
	}
}

func LoadFile(filename string) (string) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error %s", err)
	}

	return string(body)
}
