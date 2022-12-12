package mssqlMiner

import (
	"fmt"
	"testing"
)

func TestMssql(t *testing.T) {
	ms, err := Mssql("SA", "p@ssW0rd", "10.1.1.128", 1433)
	if err != nil {
		fmt.Println(err)
	}
	defer ms.Db.Close()

	if err := Search(ms); err != nil {
		fmt.Println(err)
	}
}
