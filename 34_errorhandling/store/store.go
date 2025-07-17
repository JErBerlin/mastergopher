// This program builds a simple key-value string store. In this case it is used to store pairs (color, player_name).
// It is an exercise to show the use of error types and error interfaces.

package main

import (
	"errors"
	"fmt"
	"log/slog"
)

type errNotFound struct {
	key string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("couldn't find any stored value for the key: %s", e.key)
}

func main() {
	s := NewStore()

	// check if a given key exists
	// and store a new pair, if the key is unknown
	// but produce a panic if getting an error other that not found
	_, err := s.Get("red")
	var errnf errNotFound
	if err != nil {
		if !errors.As(err, &errnf) {
			slog.Error("unexpected error", "err", err)
			panic(err)
		}
		slog.Error("first try", "err", errnf)
		s.Set("red", "big loser")
	}

	// check for the stored key
	// and produce a panic if the pair was not stored
	v, err := s.Get("red")
	if err != nil {
		slog.Error("retrieving sotred key failed", "err", err)
		panic(err)
	}
	slog.Info("value retrieved for red", "value", v)
}

type Store struct {
	keys map[string]string
}

func NewStore() *Store {
	return &Store{
		keys: make(map[string]string),
	}
}

func (s *Store) Set(k, v string) {
	s.keys[k] = v
}

func (s *Store) Get(k string) (string, error) {
	if v, ok := s.keys[k]; ok {
		return v, nil
	}
	return "", errNotFound{k}
}
