package closer

import (
	"os"
	"os/signal"
	"sync"

	"github.com/maisiq/go-words-jar/internal/logger"
)

type CloseFunc func() error

type Closer struct {
	funcs []CloseFunc
	once  sync.Once
	mu    sync.Mutex
	wait  chan interface{}
}

func New(sig ...os.Signal) *Closer {

	closer := &Closer{
		wait: make(chan interface{}, 1),
	}

	notifyChan := make(chan os.Signal, len(sig))
	signal.Notify(notifyChan, sig...)

	go func() {
		<-notifyChan
		signal.Stop(notifyChan)
		closer.CloseAll()
	}()
	return closer

}

func (c *Closer) Add(fn CloseFunc) {
	c.mu.Lock()
	c.funcs = append(c.funcs, fn)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.wait
}

func (c *Closer) CloseAll() {
	logger.Debugw("starting to close all services")

	c.once.Do(func() {
		defer close(c.wait)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errch := make(chan error, len(funcs))

		for _, fn := range funcs {
			go func(fn CloseFunc) {
				errch <- fn()
			}(fn)
		}

		for i := 0; i < cap(errch); i++ {
			if err := <-errch; err != nil {
				logger.Errorw("failed to close service", "error", err.Error())
			}
		}

	})
}
