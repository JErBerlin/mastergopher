package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Open db: %s", err.Error())
	}

	// Load data
	if err := load(db); err != nil {
		log.Fatalf("Load data: %s", err.Error())
	}

	// Get users
	users, err := GetUsers(db)
	if err != nil {
		log.Fatalf("Get users: " + err.Error())
	}

	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}
}

// users
type User struct {
	ID   string
	Name string
}

func GetUsers(db *sql.DB) ([]User, error) {
	q := "select * from users;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func load(db *sql.DB) error {
	_, err := db.Exec(`
		create table users(id text, name text);
		insert into users(id, name) values ('abdc0001', 'Leon'), ('ghdj0002', 'Francis'), ('rlhy0003', 'Pizzi'), ('etyu0004', 'Tosca');
	`)
	return err
}
