package Closer

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO whait ws (for chat)
type Closer struct {
	CloseTimeout time.Duration
	CloseFuncts  []Cloasble

	Logger *zap.Logger
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
	c.Logger.Info("Stopping application ....")
	done := make(chan bool)
	go c.Stop(done)

	timer := time.NewTimer(c.CloseTimeout)

	select {
	case <-timer.C:
		fmt.Printf("Timeout stopping application\n")
	case <-done:
		fmt.Printf("Successfully stopped.\n")
	}
}

func (c *Closer) Stop(done chan bool) {
	for _, f := range c.CloseFuncts {
		fmt.Printf("Exiting function. Function name: %s\n", f.Name)

		err := f.F()

		if err != nil {
			fmt.Printf("Error in %s! Error %s\n", f.Name, err)
		}
	}
	done <- true
}
