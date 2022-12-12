package mssqlMiner

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"regexp"
)

type DatabaseMiner interface {
	GetSchema() (*Schema, error)
}

type Schema struct {
	Databases []Database
}

type Database struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name   string
	Column []string
}

func Search(m DatabaseMiner) error {
	s, err := m.GetSchema()
	if err != nil {
		return err
	}

	re := getRegex()
	for _, database := range s.Databases {
		fmt.Println(database)
		for _, table := range database.Tables {
			for _, filed := range table.Column {
				for _, r := range re {
					if r.MatchString(filed) {
						fmt.Printf("[+] HIT: %s\n", filed)
					}
				}
			}
		}
	}
	return nil
}

func getRegex() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)ccnum`),
		regexp.MustCompile(`(?i)pass(word)?`),
	}
}

func (s Schema) String() string {
	var ret string
	for _, database := range s.Databases {
		ret += fmt.Sprint(database.String() + "\n")
	}
	return ret
}

func (d Database) String() string {
	ret := fmt.Sprintf("[DB] = %s\n", d.Name)
	for _, table := range d.Tables {
		ret += table.String()
	}
	return ret
}

func (t Table) String() string {
	ret := fmt.Sprintf("		[TABLE] = %s\n", t.Name)
	for _, field := range t.Column {
		ret += fmt.Sprintf("			[COL] = %s\n", field)
	}
	return ret
}

type MsMiner struct {
	Usr    string
	Passwd string
	Host   string
	Port   int
	//Dbase  string
	Db sql.DB
}

func Mssql(usr, passwd, host string, port int) (*MsMiner, error) {
	ms := MsMiner{
		Usr:    usr,
		Passwd: passwd,
		Host:   host,
		Port:   port,
		//Dbase:  dbase,
	}
	err := ms.connect()
	if err != nil {
		fmt.Println(err)
	}
	return &ms, nil
}

func (ms *MsMiner) connect() error {
	conn, err := sql.Open("sqlserver", fmt.Sprintf("sqlserver://%s:%s@%s:%d/instance?database=TestDB&encrypt=disable", ms.Usr, ms.Passwd, ms.Host, ms.Port))
	if err != nil {
		fmt.Println(err)
	}
	ms.Db = *conn
	return err
}

func (ms *MsMiner) GetSchema() (*Schema, error) {
	var s = new(Schema)
	//select name from syscolumns where id=object_id('table_name')
	q := `select name from syscolumns where id=object_id('table_name')`
	rows, err := ms.Db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prevschema, prevtable string
	var db Database
	var table Table
	for rows.Next() {
		var currschema, currtable, currcol string
		if err := rows.Scan(&currschema, &currtable, &currcol); err != nil {
			return nil, err
		}

		if currschema != prevschema {
			if prevschema != "" {
				db.Tables = append(db.Tables, table)
				s.Databases = append(s.Databases, db)
			}
			db = Database{Name: currschema, Tables: []Table{}}
			prevschema = currschema
			prevtable = ""
		}

		if currtable != prevtable {
			if prevtable != "" {
				db.Tables = append(db.Tables, table)
			}
			table = Table{Name: currtable, Column: []string{}}
			prevtable = currtable
		}
		table.Column = append(table.Column, currcol)
	}
	db.Tables = append(db.Tables, table)
	s.Databases = append(s.Databases, db)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return s, err
}
