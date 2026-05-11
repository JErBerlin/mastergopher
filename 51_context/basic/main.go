package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory")
	if err != nil {
		log.Fatal(err)
	}

	// Load sample data
	var loaded bool
	if loaded, err = dataLoaded(db); err != nil {
		log.Fatal(err)
	}
	if !loaded {
		fmt.Println("data was not loaded.. will load now")

		if err := load(db); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("data was loaded.. will retrieve users now")
	}

	// Get all users
	users, err := GetUsers(db)
	if err != nil {
		log.Fatal(db)
	}
	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}

}

// dataLoaded check if the table 'users' already exists and assumes that if created, the table will contain data
func dataLoaded(db *sql.DB) (bool, error) {
	var name string
	err := db.QueryRow(`
		SELECT name FROM sqlite_master
		WHERE type='table' AND name='users';
	`).Scan(&name)

	if err == sql.ErrNoRows {
		// table does not exist
		return false, nil
	} else if err != nil {
		// real error
		return false, err
	}

	// table exists
	return true, nil
}

func load(db *sql.DB) error {
	_, err := db.Exec(`
		create table users(id int, name text);
		insert into users(id, name) values (1, 'Blue'), (2, 'Red'), (3, 'Green'), (4, 'Gold');
	`)
	return err
}

// section: users
type User struct {
	ID   int
	Name string
}

func GetUsers(db *sql.DB) ([]User, error) {
	// TODO: update function to accept context parameter
	ctx := context.TODO()

	rows, err := db.QueryContext(ctx, "select * from users;")
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

// section: users
