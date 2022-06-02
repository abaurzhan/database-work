package main

import (
	"database/sql"
	"fmt"
	"log"

	// Since database driver doesn't used directly, it should be imported with underscore.
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var db *sql.DB
	var err error

	// db is represent database connection pool [conn1,conn2,..., connN]
	// Don't Open() database every time you want to send the query, use
	// *sql.DB object instead.
	db, err = sql.Open("sqlite3", "newdb.sqlite")
	defer db.Close()
	printIfErr(err)

	// Check if the database actually available.
	err = db.Ping()
	printIfErr(err)
	fmt.Println("Connected!")

	// Database is empty. No data, no tables created yet so nothing to query.
	// Let's create a new table called `TEST_TABLE` using the query below.
	sqlQeury := `CREATE TABLE TEST_TABLE(
		id INT PRIMARY KEY,
		fname varchar(25),
		lname varchar(25),
		address varchar(100),
		bio text
	);`
	// I need to create a table and i don't care anything else,
	// so i just use simple quiry execution using Exec() method.
	_, err = db.Exec(sqlQeury)
	printIfErr(err)
	fmt.Println("Table `TEST_TABLE` created.")

	// Let's see what tables present in the database. Note: `sqlite_master` is
	// the internal table created by database itself for service information purpose.
	// Since the Query() method return all the result rows, you need to iterate over it in cycle.
	showAllTables := `SELECT name FROM sqlite_master WHERE type='table';`
	rows, err := db.Query(showAllTables)
	defer rows.Close()
	printIfErr(err)

	// read result set to variable.
	var tableName string
	var resultSet []string

	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Println(err.Error())
			return
		}
		resultSet = append(resultSet, tableName)
	}

	// in some cases Scan() method can defer an error until the end of scanning,
	// so it must be checked at the end of retrieving results.
	if err = rows.Err(); err != nil {
		log.Println(err.Error())
		return
	}
	// Print existing tables.
	fmt.Println(resultSet)

	//Now let's write some data into the table using transaction mode.
	//Transaction guarantees that everything be fully executed or nothing will be changed.
	tx, err := db.Begin()
	printIfErr(err)
	_, err = tx.Exec("INSERT INTO TEST_TABLE (id, fname, lname, address, bio) VALUES (?,?,?,?,?)", 1, "Humfried", "Ritelli", "6 Center Road", "repurpose extensible systems")
	if err != nil {
		log.Println(err.Error())
		log.Println("transaction failed")
		// if any error rise transaction must be rolled back.
		tx.Rollback()
		return
	}

	// transaction will apply all changes if no error rise.
	if err := tx.Commit(); err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("transaction commited!")

}

func printIfErr(e error) {
	if e != nil {
		log.Println(e.Error())
	}
}
