package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Important: allow >1 open connection so we can obtain two different connections.
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)

	// Get two distinct connections from the pool.
	c1, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c1.Close()

	c2, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c2.Close()

	log.Println("== Create table + insert on connection #1 ==")
	_, err = c1.ExecContext(ctx, `
		CREATE TABLE users(id INTEGER, name TEXT);
		INSERT INTO users(id, name) VALUES (1,'Blue'), (2,'Red');
	`)
	if err != nil {
		log.Fatalf("c1 exec failed: %v", err)
	}
	log.Println("c1: created table + inserted rows")

	log.Println("== Query on connection #2 (same DSN ':memory:') ==")
	var n2 int
	err = c2.QueryRowContext(ctx, `SELECT count(*) FROM users;`).Scan(&n2)
	if err != nil {
		log.Printf("c2: expected failure: %v", err)
	} else {
		log.Printf("c2: UNEXPECTED success, count=%d", n2)
	}

	log.Println("== Query on connection #1 (should work) ==")
	var n1 int
	err = c1.QueryRowContext(ctx, `SELECT count(*) FROM users;`).Scan(&n1)
	if err != nil {
		log.Fatalf("c1 query failed: %v", err)
	}
	log.Printf("c1: success, count=%d", n1)
}
