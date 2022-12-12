/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"log"
	"regexp"
)

// SqlMinerCmd represents the SqlMiner command
var SqlMinerCmd = &cobra.Command{
	Use:   "SqlMiner",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SqlMiner called")
		sqlMiner()
	},
}

var usr string
var passwd string
var host string
var port int
var dbase string
var Keyword []string

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
	Name    string
	Columns []string
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
			for _, field := range table.Columns {
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

//func getRegex() []*regexp.Regexp {
//	return []*regexp.Regexp{
//		//regexp.MustCompile(`(?i)social`),
//		//regexp.MustCompile(`(?i)ssn`),
//		//regexp.MustCompile(`(?i)pass(word)?`),
//		//regexp.MustCompile(`(?i)hash`),
//		//regexp.MustCompile(`(?i)ccnum`),
//		//regexp.MustCompile(`(?i)card`),
//		//regexp.MustCompile(`(?i)security`),
//		//regexp.MustCompile(`(?i)key`),
//		regexp.MustCompile(keyword),
//	}
//}

func getRegex() []*regexp.Regexp {
	//Keyword = os.Args[1]
	var re []*regexp.Regexp
	for _, i := range Keyword {
		re = append(re, regexp.MustCompile(i))
	}
	return re
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

func (t Table) String() string {
	ret := fmt.Sprintf("    [TABLE] = %+s\n", t.Name)
	for _, field := range t.Columns {
		ret += fmt.Sprintf("       [COL] = %+s\n", field)
	}
	return ret
}

type MySQLMiner struct {
	Usr      string
	Passwd   string
	Host     string
	Port     int
	Database string
	Db       sql.DB
}

func Newsql(usr string, passwd string, host string, port int) (*MySQLMiner, error) {
	m := MySQLMiner{
		Usr:    usr,
		Passwd: passwd,
		Host:   host,
		Port:   port,
		//Database: database,
	}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MySQLMiner) connect() error {

	//db, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(%s:3306)/information_schema", m.Host))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema", m.Usr, m.Passwd, m.Host, m.Port))
	if err != nil {
		log.Panicln(err)
	}
	m.Db = *db
	return nil
}

func (m *MySQLMiner) GetSchema() (*Schema, error) {
	var s = new(Schema)

	sql := `SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM columns
    WHERE TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
    ORDER BY TABLE_SCHEMA, TABLE_NAME`
	schemarows, err := m.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer schemarows.Close()

	var prevschema, prevtable string
	var db Database
	var table Table
	for schemarows.Next() {
		var currschema, currtable, currcol string
		if err := schemarows.Scan(&currschema, &currtable, &currcol); err != nil {
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
			table = Table{Name: currtable, Columns: []string{}}
			prevtable = currtable
		}
		table.Columns = append(table.Columns, currcol)
	}
	db.Tables = append(db.Tables, table)
	s.Databases = append(s.Databases, db)
	if err := schemarows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func sqlMiner() {
	mm, err := Newsql(usr, passwd, host, port)
	//mm, err := Newsql("root", "password", "10.1.1.128",3306 , "store")
	if err != nil {
		panic(err)
	}
	defer mm.Db.Close()

	if err := Search(mm); err != nil {
		fmt.Println(err)
		//getRegex()
	}
}

func init() {
	rootCmd.AddCommand(SqlMinerCmd)
	rootCmd.Flags().StringVarP(&host, "host", "H", "", "database address")
	SqlMinerCmd.Flags().StringVarP(&usr, "usr", "U", "", "database use")
	SqlMinerCmd.Flags().StringVarP(&passwd, "passwd", "P", "", "database password")
	SqlMinerCmd.Flags().StringVarP(&host, "host", "H", "", "database address")
	SqlMinerCmd.Flags().IntVarP(&port, "port", "p", 0, "database service number")
	SqlMinerCmd.Flags().StringSliceVarP(&Keyword, "keyword", "k", []string{}, "type keywords to quarry databases.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// SqlMinerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// SqlMinerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
