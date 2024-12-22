package closer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CloseConcurrently(t *testing.T) {
	var (
		timeout = 3 * time.Second

		errorTimeout   = errors.New("closing process timed out or was canceled")
		errorAggregate = errors.New("closing finished with error(s):\n[!] error 1\n[!] error 2")

		closersWithTimeout = []func(context.Context) error{
			func(context.Context) error { time.Sleep(1 * time.Second); return nil },
			func(context.Context) error { time.Sleep(4 * time.Second); return nil },
		}
		closersWithErrors = []func(context.Context) error{
			func(context.Context) error { time.Sleep(1 * time.Second); return errors.New("error 1") },
			func(context.Context) error { time.Sleep(2 * time.Second); return errors.New("error 2") },
		}
		closersWithoutErrors = []func(context.Context) error{
			func(context.Context) error { return nil },
			func(context.Context) error { return nil },
		}
	)

	tests := []struct {
		name     string
		closers  []func(context.Context) error
		timeout  time.Duration
		expected error
	}{
		{
			name:     "completed without errors",
			closers:  closersWithoutErrors,
			timeout:  timeout,
			expected: nil,
		},
		{
			name:     "completed with errors",
			closers:  closersWithErrors,
			timeout:  timeout,
			expected: errorAggregate,
		},
		{
			name:     "timeout occurs",
			closers:  closersWithTimeout,
			timeout:  timeout,
			expected: errorTimeout,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New()
			for _, closer := range test.closers {
				c.Add(closer)
			}

			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			err := c.CloseConcurrently(ctx)
			assert.Equal(t, test.expected, err)
		})
	}
}

func Test_CloseSequentially(t *testing.T) {
	var (
		timeout = 3 * time.Second

		errorTimeout   = errors.New("closing finished with error(s):\n[!] function 1 timed out")
		errorAggregate = errors.New("closing finished with error(s):\n[!] error 2\n[!] error 1")

		closersWithTimeout = []func(context.Context) error{
			func(context.Context) error { time.Sleep(1 * time.Second); return nil },
			func(context.Context) error { time.Sleep(4 * time.Second); return nil },
		}
		closersWithErrors = []func(context.Context) error{
			func(context.Context) error { return errors.New("error 1") },
			func(context.Context) error { return errors.New("error 2") },
		}
		closersWithoutErrors = []func(context.Context) error{
			func(context.Context) error { return nil },
			func(context.Context) error { return nil },
		}
	)

	tests := []struct {
		name     string
		closers  []func(context.Context) error
		timeout  time.Duration
		expected error
	}{
		{
			name:     "completed without errors",
			closers:  closersWithoutErrors,
			timeout:  timeout,
			expected: nil,
		},
		{
			name:     "completed with errors",
			closers:  closersWithErrors,
			timeout:  timeout,
			expected: errorAggregate,
		},
		{
			name:     "timeout occurs",
			closers:  closersWithTimeout,
			timeout:  timeout,
			expected: errorTimeout,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := New()
			for _, closer := range test.closers {
				c.Add(closer)
			}

			ctx := context.Background()
			err := c.CloseSequentially(ctx, test.timeout)
			assert.Equal(t, test.expected, err)
		})
	}
}
