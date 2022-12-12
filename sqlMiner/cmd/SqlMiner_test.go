package cmd

import "testing"

func TestSqlMiner(t *testing.T) {
	//mm, err := Newsql(usr, passwd, host, port)
	mm, err := Newsql("root", "password", "10.1.1.128",3306)
	if err != nil {
		panic(err)
	}
	defer mm.Db.Close()

	if err := Search(mm); err != nil {
		panic(err)
	}
}