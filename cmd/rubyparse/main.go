package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func loadFile(filename string) (string) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error %s", err)
	}

	return string(body)
}

func tokenize(body string) ([]string) {
	lines := strings.Split(body, "\n")

	tokens := make([]string, 0)

	for _,line := range lines {
		line_tokens := strings.Fields(line)
		partOfString := false
		// whitespace
		// start of a string
		// in a string
		// end of a string
		// token

		for _,token := range line_tokens {
			fmt.Printf("--- %v \t", token)
			switch {
			case partOfstring:
				if strings.HasSuffix(token, "\"") || strings.HasSuffix(token, "'") {
					fmt.Printf("end of string ")
				} else {
					fmt.Printf("part of string ")
				}

			case strings.HasPrefix(token, "\"") || strings.HasPrefix(token, "'"):
				fmt.Printf("start of string ")
				partOfString = true
				tokens = append(tokens, token)
			case true:
				fmt.Printf("default ")
				tokens = append(tokens, token)
			}
			fmt.Printf("\n")
		}
	}
	return tokens
}

func main() {
	file := loadFile("test.rb")
	tokenize(file)

//	for _,token := range tokens {
//		fmt.Printf("%T, \t %v\n", token, token)
//	}
}
