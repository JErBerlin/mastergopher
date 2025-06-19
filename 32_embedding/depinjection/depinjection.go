// This program demonstrates dependency injection and decoupling through interface abstraction.
// The Service depends only on a Logger interface, not on log.Logger directly.
// This allows substituting the logging mechanism without changing the Service logic.

package main

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type Service struct {
	Name string
	Logger
}

func NewService(n string, l Logger) *Service {
	return &Service{
		Name:   n,
		Logger: l,
	}
}

func (s *Service) Init() {
	s.Logger.Printf("initializing : %s\n", s.Name)
}

func (s *Service) Stop() {
	s.Logger.Printf("shutting down: %s\n", s.Name)
}

func main() {
	l := log.New(os.Stdout, "[service] ", log.Lshortfile)
	s := NewService("example", l)
	s.Init()
	s.Stop()
}
