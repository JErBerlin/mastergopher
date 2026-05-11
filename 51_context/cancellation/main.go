package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load sample data
	var loaded bool
	if loaded, err = dataLoaded(db); err != nil {
		log.Fatal(err)
	}
	if !loaded {
		fmt.Println("data was not loaded.. will load now")

		if err := load(ctx, db); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("data was loaded.. will retrieve users now")
	}

	// Get all users
	users, err := GetUsers(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}
    /*	
	for {
		time.Sleep(1*time.Second)
		fmt.Println("*")
	}
	*/
}

// section: load

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

func load(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`
			create table users(id int, name text);
			insert into users(id, name) values (1, 'Blue'), (2, 'Red'), (3, 'Green'), (4, 'Gold');
		`,
	)
	return err
}

// section:load

// section: users

type User struct {
	ID   int
	Name string
}

func GetUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	// Set a timeout on the query
	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// make error channel to return error from goroutine
	ec := make(chan error, 1)

	// make channel to retrieve result
	result := make(chan []User, 1)

	// launch the query in a goroutine
	go func() {
		// pass query bound context
		rows, err := db.QueryContext(ctx, "select * from users;")
		if err != nil {
			ec <- err
			return
		}

		time.Sleep(timeout)

		users := []User{}
		for rows.Next() {
				ec <- err
				return
			}
			users = append(users, u)
		}
		result <- users
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ec:
		return nil, err
	case u := <-result:
		return u, nil
	}
}

// section: users
