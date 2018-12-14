package main

import (
	"database/sql"
	"fmt"
	//added here because the postgres driver needs this
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "lenslocked_dev"
)

func main() {
	plsqlString := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	fmt.Print(plsqlString)

	db, err := sql.Open("postgres", plsqlString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO users(name, email) VALUES ($1, $2)`,
		"Mihir Kelkar",
		"mihir@mihirkelkar.co")

	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to database successful")
	db.Close()
}
