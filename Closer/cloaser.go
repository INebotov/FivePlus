package Closer

import (
	"Backend/Logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// TODO whait ws
type Closer struct {
	CloseTimeout time.Duration
	CloseFuncts  []Cloasble

	Logger Logger.Logger
}
type Cloasble struct {
	Name string
	F    func() error
}

func (c *Closer) Add(name string, f func() error) {
	c.CloseFuncts = append(c.CloseFuncts, Cloasble{
		Name: name,
		F:    f,
	})
}

func (c *Closer) Listen() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-s

	defer os.Exit(0)
	var wg sync.WaitGroup
	for _, f := range c.CloseFuncts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Logger.Info("Exiting function: %s", f.Name)
			err := f.F()
			if err != nil {
				c.Logger.Error("Error exiting in function: %s. Error: %s", f.Name, err)
			}
		}()
	}
	done := make(chan bool)
	go func(done chan bool) {
		wg.Wait()
		done <- true
	}(done)
	tic := time.NewTimer(c.CloseTimeout)
	select {
	case <-tic.C:
		c.Logger.Error("Timeout Closing Service!")
	case <-done:
		c.Logger.Info("Successfully exited!")
	}
}
