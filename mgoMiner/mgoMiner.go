package mgoMiner

import (
	"fmt"
	"regexp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseMiner interface {
	GetSchema() (*Schema, error)
}

type Schema struct {
	Databases []Database
}

type Database struct {
	Name   string
	Tables []Collection
}

type Collection struct {
	Name   string
	Fields []string
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
			for _, field := range table.Fields {
				for _, r := range re {
					if r.MatchString(field) {
						fmt.Printf("[+] HIT: %s\n", field)
					}
				}
			}
		}
	}
	return nil
}

func getRegex() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)social`),
		regexp.MustCompile(`(?i)ssn`),
		regexp.MustCompile(`(?i)pass(word)?`),
		regexp.MustCompile(`(?i)hash`),
		regexp.MustCompile(`(?i)ccnum`),
		regexp.MustCompile(`(?i)card`),
		regexp.MustCompile(`(?i)security`),
		regexp.MustCompile(`(?i)key`),
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
	ret := fmt.Sprintf("[DB] = %+s\n", d.Name)
	for _, table := range d.Tables {
		ret += table.String()
	}
	return ret
}

func (t Collection) String() string {
	ret := fmt.Sprintf("    [TABLE] = %+s\n", t.Name)
	for _, field := range t.Fields {
		ret += fmt.Sprintf("       [COL] = %+s\n", field)
	}
	return ret
}

/* Extranneous code omitted for brevity */

type MgoMiner struct {
	Host string
	//Port    int
	//Dbase   string
	session *mgo.Session
}

func Mongo(host string) (*MgoMiner, error) {
	m := MgoMiner{
		Host: host,
	}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MgoMiner) connect() error {
	s, err := mgo.Dial(m.Host)
	if err != nil {
		return err
	}
	m.session = s
	return nil
}

func (m *MgoMiner) GetSchema() (*Schema, error) {
	var s = new(Schema)

	dbnames, err := m.session.DatabaseNames()
	if err != nil {
		return nil, err
	}

	for _, dbname := range dbnames {
		db := Database{Name: dbname, Tables: []Collection{}}
		Tables, err := m.session.DB(dbname).CollectionNames()
		if err != nil {
			return nil, err
		}

		for _, collection := range Tables {
			table := Collection{Name: collection, Fields: []string{}}

			var docRaw bson.Raw
			err := m.session.DB(dbname).C(collection).Find(nil).One(&docRaw)
			if err != nil {
				return nil, err
			}

			var doc bson.RawD
			if err := docRaw.Unmarshal(&doc); err != nil {
				if err != nil {
					return nil, err
				}
			}

			for _, f := range doc {
				table.Fields = append(table.Fields, f.Name)
			}
			db.Tables = append(db.Tables, table)
		}
		s.Databases = append(s.Databases, db)
	}
	return s, nil
}
