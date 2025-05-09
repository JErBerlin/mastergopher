package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const queryTimeout = 1 * time.Microsecond // adjust to test timeout behavior

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("open db: %s", err.Error())
	}

	if err := load(db); err != nil {
		log.Fatalf("load data from db: %s", err.Error())
	}

	users, err := GetUsers(ctx, db)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Error: get users: Request timed out (context deadline exceeded)")
		} else {
			log.Fatalf("get users: %s", err.Error())
		}
		return
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

// GetUsers runs a select all users on the db but using a context for possible cancellation.
func GetUsers(db *sql.DB) ([]User, error) {
	ctx := context.TODO() // context is expected to be provided by the caller in a later implementation

	q := "select * from users;"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
        insert into users(id, name) values 
            ('abdc0001', 'Leon'), 
            ('ghdj0002', 'Francis'), 
            ('rlhy0003', 'Pizzi'), 
            ('etyu0004', 'Tosca');
    `)
	return err
}
