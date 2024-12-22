package closer

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
	"sync"
)

type Closer struct {
	mu   sync.Mutex
	exec []func() error
}

func New() *Closer {
	return &Closer{
		mu:   sync.Mutex{},
		exec: make([]func() error, 0),
	}
}

func (c *Closer) Add(f func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.exec = append(c.exec, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		wg       = sync.WaitGroup{}
		mu       = sync.Mutex{}
		messages = make([]string, 0, len(c.exec))
		complete = make(chan struct{})
	)

	for _, exec := range c.exec {
		wg.Add(1)
		go func(exec func() error) {
			defer wg.Done()
			if err := exec(); err != nil {
				mu.Lock()
				defer mu.Unlock()
				messages = append(messages, err.Error())
			}
		}(exec)
	}

	go func() {
		wg.Wait()
		close(complete)
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return errors.New("closing process timed out or was canceled")
	}

	if len(messages) > 0 {
		return fmt.Errorf("closing finished with error(s):\n%s", strings.Join(messages, "\n"))
	}

	return nil
}
