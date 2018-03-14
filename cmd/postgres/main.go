package main

import (
	"database/sql"
	"log"
	"fmt"
	"strings"
	"encoding/csv"
	"errors"
	"strconv"
	"math/rand"

	_ "github.com/lib/pq"
	"github.com/kennygrant/sanitize"
)

func main() {


	connStr := "postgresql://localhost?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	

	// _, err = db.Query("CREATE TABLE users (id integer, name varchar(40));")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = db.Query("INSERT INTO users VALUES(1, 'samuel');")
	if err != nil {
		log.Fatal(err)
	}
	
	insertId := 1
	rows, err := db.Query("SELECT name FROM users WHERE id = $1", insertId)
	if err != nil {
		log.Fatal(err)
	}

	var name string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func createTable(db *sql.DB, schema string, tableName string, columns []string) error {
	columnTypes := make([]string, len(columns))
	for i, col := range columns {
		columnTypes[i] = fmt.Sprintf("%s TEXT", col)
	}
	columnDefinitions := strings.Join(columnTypes, ",")
	fullyQualifiedTable := fmt.Sprintf("%s.%s", schema, tableName)
	tableSchema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", fullyQualifiedTable, columnDefinitions)

	_, err := db.Query(tableSchema)
	return err
}

// Parse columns from first header row or from flags
func parseColumns(reader *csv.Reader, skipHeader bool, fields string) ([]string, error) {
	var err error
	var columns []string
	if fields != "" {
		columns = strings.Split(fields, ",")

		if skipHeader {
			reader.Read() //Force consume one row
		}
	} else {
		columns, err = reader.Read()
		if err != nil {
			return nil, err
		}
	}

	for _, col := range columns {
		if containsDelimiter(col) {
			return columns, errors.New("Please specify the correct delimiter with -d.\nHeader column contains a delimiter character: " + col)
		}
	}

	for i, col := range columns {
		columns[i] = postgresify(col)
	}

	return columns, nil
}

func containsDelimiter(col string) bool {
	return strings.Contains(col, ";") || strings.Contains(col, ",") ||
		strings.Contains(col, "|") || strings.Contains(col, "\t") ||
		strings.Contains(col, "^") || strings.Contains(col, "~")
}

//Makes sure that a string is a valid PostgreSQL identifier
func postgresify(identifier string) string {
	str := sanitize.BaseName(identifier)
	str = strings.ToLower(identifier)
	str = strings.TrimSpace(str)

	replacements := map[string]string{
		" ": "_",
		"/": "_",
		".": "_",
		":": "_",
		";": "_",
		"|": "_",
		"-": "_",
		",": "_",
		"#": "_",
		
		"[":  "",
		"]":  "",
		"{":  "",
		"}":  "",
		"(":  "",
		")":  "",
		"?":  "",
		"!":  "",
		"$":  "",
		"%":  "",
		"*":  "",
		"\"": "",
	}
	for oldString, newString := range replacements {
		str = strings.Replace(str, oldString, newString, -1)
	}

	if len(str) == 0 {
		str = fmt.Sprintf("_col%d", rand.Intn(10000))
	} else {
		firstLetter := string(str[0])
		if _, err := strconv.Atoi(firstLetter); err == nil {
			str = "_" + str
		}
	}

	return str
}