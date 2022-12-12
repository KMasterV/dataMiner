package main

import (
	"fmt"
	. "mssqlMiner"
)

var usr string
var passwd string
var host string
var dbase string
var port int

func main() {
	ms, err := Mssql(usr, passwd, host, port)
	if err != nil {
		fmt.Println(err)
	}
	defer ms.Db.Close()

	if err := Search(ms); err != nil {
		fmt.Println(err)
	}
}
