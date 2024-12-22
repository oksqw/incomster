package closer

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

type Closer struct {
	mu   sync.Mutex
	exec []func(context.Context) error
}

func New() *Closer {
	return &Closer{
		mu:   sync.Mutex{},
		exec: make([]func(context.Context) error, 0),
	}
}

func (c *Closer) Add(f func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.exec = append(c.exec, f)
}

func (c *Closer) CloseConcurrently(ctx context.Context) error {
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
		go func(exec func(context.Context) error) {
			defer wg.Done()
			if err := exec(ctx); err != nil {
				mu.Lock()
				defer mu.Unlock()
				messages = append(messages, fmt.Sprintf("[!] %s", err.Error()))
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

func (c *Closer) CloseSequentially(ctx context.Context, timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		mu       = sync.Mutex{}
		messages = make([]string, 0, len(c.exec))
	)

	addMessage := func(message string) {
		mu.Lock()
		defer mu.Unlock()
		messages = append(messages, message)
	}

	for i := len(c.exec) - 1; i >= 0; i-- {
		fn := c.exec[i]
		fnctx, cancel := context.WithTimeout(ctx, timeout)
		done := make(chan error, 1)

		go func() {
			done <- fn(fnctx)
		}()

		select {
		case err := <-done:
			if err != nil {
				addMessage(fmt.Sprintf("[!] %s", err.Error()))
			}
		case <-fnctx.Done():
			addMessage(fmt.Sprintf("[!] function %d timed out", i))
		}

		cancel()
	}

	if len(messages) > 0 {
		return fmt.Errorf("closing finished with error(s):\n%s", strings.Join(messages, "\n"))
	}

	return nil
}
