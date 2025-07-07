// This program builds a simple key-value string store. In this case it is used to store pairs (color, player_name).
// It is an exercise to show the use of error types and error interfaces.

package main

import (
	"fmt"
	"log"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const errNotFound = Error("couldn't find any stored value for the given key")

func main() {
	s := NewStore()

	// check if a given key exists
	// and store a new pair, if the key is unknown
	// but produce a panic if getting an error other that not found
	_, err := s.Get("red")
	if err != nil {
		if err == errNotFound {
			s.Set("red", "big loser")
		} else {
			log.Fatal(err)
		}
	}

	// check for the stored key
	// and produce a panic if the pair was not stored
	v, err := s.Get("red")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", v)
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
	return "", errNotFound
}
