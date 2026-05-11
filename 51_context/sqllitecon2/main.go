package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string
}

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Allow the pool to use more than one connection.
	db.SetMaxOpenConns(5)

	ctx := context.Background()

	if err := load(ctx, db); err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 50; i++ {
		log.Printf("=== round %d ===", i)
		if err := round(ctx, db); err != nil {
			log.Printf("round %d: %v", i, err)
		}
	}
}

func round(ctx context.Context, db *sql.DB) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Occupy one pool connection with a slow query.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := holdConnectionWithSlowQuery(ctx, db); err != nil {
			errCh <- fmt.Errorf("slow query: %w", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	// In parallel, try to read users.
	wg.Add(1)
	go func() {
		defer wg.Done()
		users, err := GetUsers(ctx, db)
		if err != nil {
			errCh <- fmt.Errorf("GetUsers: %w", err)
			return
		}
		log.Printf("GetUsers ok: %d users", len(users))
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

func load(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE users(id INTEGER, name TEXT);
		INSERT INTO users(id, name) VALUES
			(1, 'Blue'), (2, 'Red'), (3, 'Green'), (4, 'Gold');
	`)
	return err
}

func GetUsers(parent context.Context, db *sql.DB) ([]User, error) {
	ctx, cancel := context.WithTimeout(parent, 200*time.Millisecond)
	defer cancel()

	rows, err := db.QueryContext(ctx, `SELECT id, name FROM users;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func holdConnectionWithSlowQuery(parent context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	const N = 20000000

	query := fmt.Sprintf(`
		WITH RECURSIVE cnt(x) AS (
			SELECT 1
			UNION ALL
			SELECT x + 1 FROM cnt WHERE x < %d
		)
		SELECT sum(x) FROM cnt;
	`, N)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sum int64
	if rows.Next() {
		if err := rows.Scan(&sum); err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return err
			}
			return err
		}
	}
	return rows.Err()
}
