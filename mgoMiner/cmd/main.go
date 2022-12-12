package main

import (
	"fmt"
	. "mgoMiner"
)

var Host string

func main() {
	mm, err := Mongo("mongodb://10.1.1.128")
	if err != nil {
		fmt.Println(err)
	}
	if err := Search(mm); err != nil {
		fmt.Println(err)
	}
}
