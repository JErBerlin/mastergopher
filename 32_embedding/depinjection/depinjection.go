// This program demonstrates decoupling external dependencies using dependency injection.
// A general Tracker interface is defined, and Trackerlog is a concrete implementation based on log.Logger,
// which works as an adapter (it adapts the external dependency log.Logger to the internal Tracker interface).
// The Service depends only on the Tracker interface, not on any specific logging implementation.

package main

import (
	"fmt"
	"log"
	"os"
)

type Tracker interface {
	Print(string)
}

type Trackerlog struct {
	logger *log.Logger
}

func (t Trackerlog) Print(s string) {
	t.logger.Print(s)
}

type Service struct {
	Name string
	Tracker
}

func NewService(n string, t Tracker) *Service {
	return &Service{
		Name:    n,
		Tracker: t,
	}
}

func (s *Service) Init() {
	s.Print(fmt.Sprintf("initializing : %s\n", s.Name))
}

func (s *Service) Stop() {
	s.Print(fmt.Sprintf("shutting down: %s\n", s.Name))
}

func main() {
	l := log.New(os.Stdout, "[service] ", log.Lshortfile)
	t := Trackerlog{l}

	s := NewService("example", t)
	s.Init()
	s.Stop()
}
