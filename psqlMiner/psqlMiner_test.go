package psqlMiner

import (
	"fmt"
	"testing"
)

func TestPsql(t *testing.T) {
	pq, err := Psql("postgres", "password", "10.1.1.128", "store")
	if err != nil {
		fmt.Println(err)
	}
	defer pq.Db.Close()

	if err := Search(pq); err != nil {
		fmt.Println(err)
	}
}
