package psqlMiner

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
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

type PqMiner struct {
	Usr    string
	Passwd string
	Host   string
	Port   int
	Dbase  string
	Db     sql.DB
}

func Psql(usr, passwd, host, dbase string) (*PqMiner, error) {
	p := PqMiner{
		Usr:    usr,
		Passwd: passwd,
		Host:   host,
		//Port:   port,
		Dbase: dbase,
	}
	err := p.connect()
	if err != nil {
		log.Fatal(err)
	}
	return &p, nil
}

func (p *PqMiner) connect() error {
	// "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", p.Usr, p.Passwd, p.Host, p.Dbase))
	if err != nil {
		log.Fatal(err)
	}
	p.Db = *conn
	return nil
}

func (p *PqMiner) GetSchema() (*Schema, error) {
	var s = new(Schema)
	//q := `SELECT column_name FROM information_schema.columns WHERE table_name = 'transactions'`
	q := `SELECT table_name FROM information_schema.tables WHERE table_schema='public'`
	rows, err := p.Db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//for rows.Next() {
	//	var (
	//		id   int64
	//		name string
	//	)
	//	if err := rows.Scan(&id, &name); err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("id %d name is %s\n", id, name)
	//}
	//if !rows.NextResultSet() {
	//	log.Fatalf("expected more result sets: %v", rows.Err())
	//}
	//var roleMap = map[int64]string{
	//	1: "user",
	//	2: "admin",
	//	3: "gopher",
	//	4: "postgres",
	//}
	//for rows.Next() {
	//	var (
	//		id   int64
	//		role int64
	//	)
	//	if err := rows.Scan(&id, &role); err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("id %d has role %s\n", id, roleMap[role])
	//}
	//if err := rows.Err(); err != nil {
	//	log.Fatal(err)
	//}

	return s, err
}
