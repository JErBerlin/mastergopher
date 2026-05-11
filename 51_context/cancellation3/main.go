// main.go
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
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SQLite ":memory:" databases are per-connection.
	// database/sql may use more than one connection over time, which would mean "a different in-memory DB".
	// We set this to 1 so the demo is deterministic: the table we create is guaranteed to be visible to the query.
	db.SetMaxOpenConns(1)

	ctx := context.Background()

	loaded, err := dataLoaded(db)
	if err != nil {
		log.Fatal(err)
	}
	if !loaded {
		log.Println("data was not loaded.. will load now")
		if err := load(ctx, db); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("data was loaded.. will retrieve users now")
	}

	// Run the "slow query with short timeout" demo via GetUsers
	users, err := GetUsers(ctx, db)
	if err != nil {
		log.Printf("GetUsers returned error: %v", err)
		return
	}

	log.Printf("GetUsers succeeded, rows=%d", len(users))
	for _, u := range users {
		log.Printf("user: id=%d name=%s", u.ID, u.Name)
	}
}

// section: load

func dataLoaded(db *sql.DB) (bool, error) {
	var name string
	err := db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type='table' AND name='users'
		LIMIT 1;
	`).Scan(&name)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func load(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE users(id INTEGER, name TEXT);
		INSERT INTO users(id, name) VALUES
			(1, 'Blue'), (2, 'Red'), (3, 'Green'), (4, 'Gold');
	`)
	return err
}

// section: users

type User struct {
	ID   int
	Name string
}

func GetUsers(parent context.Context, db *sql.DB) ([]User, error) {
	timeout := 50 * time.Millisecond
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	// This is intentionally big so the query is foreseeably slow. About half of the times will time out.
	// Tune it:
	// - If the query finishes always before the timeout, increase N.
	// - If it is always too slow, decrease N.
	const N = 140000

	query := fmt.Sprintf(`
		WITH RECURSIVE cnt(x) AS (
			SELECT 1
			UNION ALL
			SELECT x + 1 FROM cnt WHERE x < %d
		)
		SELECT sum(x) FROM cnt;
	`, N)

	start := time.Now()
	log.Println("--- GetUsers slow query test ---")
	log.Printf("timeout: %s", timeout)
	log.Printf("N: %d", N)
	log.Printf("start: %s", start.Format(time.RFC3339Nano))

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		elapsed := time.Since(start)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("timed out after %s (driver likely honored ctx): %w", elapsed, err)
		}
		return nil, fmt.Errorf("query failed after %s: %w", elapsed, err)
	}
	defer rows.Close()

	// Consume the single-row result. Cancellation can surface via Next/Scan/rows.Err.
	var sum int64
	if rows.Next() {
		if err := rows.Scan(&sum); err != nil {
			elapsed := time.Since(start)
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, fmt.Errorf("timed out during scan after %s (driver likely honored ctx): %w", elapsed, err)
			}
			return nil, fmt.Errorf("scan failed after %s: %w", elapsed, err)
		}
	}
	if err := rows.Err(); err != nil {
		elapsed := time.Since(start)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("timed out during iteration after %s (driver likely honored ctx): %w", elapsed, err)
		}
		return nil, fmt.Errorf("iteration failed after %s: %w", elapsed, err)
	}

	elapsed := time.Since(start)
	log.Printf("end: %s", time.Now().Format(time.RFC3339Nano))
	log.Printf("elapsed: %s", elapsed)
	log.Printf("sum result: %d", sum)
	log.Println("--- end GetUsers slow query test ---")

	// For the demo, we return the real users too (fast query) so the function keeps its original intent.
	// You can remove this part if you want GetUsers to *only* be the slow-query test.
	rows2, err := db.QueryContext(parent, "SELECT id, name FROM users;")
	if err != nil {
		return nil, fmt.Errorf("get users: query failed: %w", err)
	}
	defer rows2.Close()

	var users []User
	for rows2.Next() {
		var u User
		if err := rows2.Scan(&u.ID, &u.Name); err != nil {
			return nil, fmt.Errorf("get users: scan row: %w", err)
		}
		users = append(users, u)
	}
	if err := rows2.Err(); err != nil {
		return nil, fmt.Errorf("get users: iteration failed: %w", err)
	}

	return users, nil
}
