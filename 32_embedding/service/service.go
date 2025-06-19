// This is a boilerplate code which can be used as a starting point for a service in a larger system.
// It uses dependency injection and has logging abstraction, a graceful shutdown and concurrency safety.
// I also makes a clear service lifecycle separation.
// Background work is performed inside the Run method's goroutine.
// Internal clean-up logic specific to the service loop should go right before returning from that goroutine.
// External or system shutdown logic should be inserted in the Stop method after closing the stop signal.

package main

import (
	"log"
	"os"
	"sync"
	"time"
)

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

type Service struct {
	Name string
	Logger
	stop chan struct{}
	mu   sync.Mutex
	once sync.Once
}

func NewService(n string, l Logger) *Service {
	stop := make(chan struct{})
	return &Service{
		Name:   n,
		Logger: l,
		stop:   stop,
	}
}

func (s *Service) Init() {
	s.Logger.Printf("initializing : %s\n", s.Name)
}

func (s *Service) Run() {
	go func() {
		for {
			select {
			case <-s.stop:
				s.Logger.Println("stopping service")
				s.mu.Lock() // locking here in case future cleanup accesses shared state
				// some internal, service-loop related clean up work
				s.mu.Unlock()
				return
			case <-time.After(time.Second):
				s.Logger.Println("service is running..")
			}
		}
	}()
}

func (s *Service) Stop() {
	s.mu.Lock() // protects potential future shared state modifications across Stop and Run
	s.once.Do(func() {
		close(s.stop)
	})
	s.mu.Unlock()
	s.Logger.Printf("shutting down: %s\n", s.Name)
	// some external/system-triggered shutdown logic her
	s.Logger.Printf("shutdown successful")
}

func main() {
	l := log.New(os.Stdout, "[service] ", log.Lshortfile)
	s := NewService("VIS service", l)
	s.Init()
	s.Run()
	// let the service run a bit before stopping it
	time.Sleep(3500 * time.Millisecond)
	s.Stop()
}
