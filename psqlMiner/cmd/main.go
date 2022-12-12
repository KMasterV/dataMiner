package main

import (
	"fmt"
	. "psqlMiner"
)

var usr string
var passwd string
var host string

//var port int
var dbase string

func main() {
	pq, err := Psql(usr, passwd, host, dbase)
	if err != nil {
		fmt.Println(err)
	}
	defer pq.Db.Close()

	if err := Search(pq); err != nil {
		fmt.Println(err)
	}
}
