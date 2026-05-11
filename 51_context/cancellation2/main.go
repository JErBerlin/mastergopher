package main

import (
	"context"
	"database/sql"
	"errors"
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
	timeout := 500 * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	time.Sleep(498 * time.Millisecond)

	rows, err := db.QueryContext(ctx, "select id, name from users;")
	if err != nil {
		// If ctx already expired/cancelled, prefer ctx error
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("get users: query timed out after %s: %w", timeout, ctxErr)
		}
		return nil, fmt.Errorf("get users: query failed: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, fmt.Errorf("get users: scan row: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		// rows.Err() is where drivers surface mid-iteration cancel/timeout
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("get users: timed out after %s: %w", timeout, err)
		}
		return nil, fmt.Errorf("get users: iteration failed: %w", err)
	}

	return users, nil
}

// section: users
