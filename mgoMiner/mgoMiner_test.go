package mgoMiner

import (
	"fmt"
	"testing"
)

func TestNewMgo(t *testing.T) {

	mm, err := Mongo("mongodb://10.1.1.128:27017/store")
	if err != nil {
		fmt.Println(err)
	}
	if err := Search(mm); err != nil {
		fmt.Println(err)
	}
}
