package closer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloser(t *testing.T) {
	var (
		errorTimeout   = errors.New("closing process timed out or was canceled")
		errorAggregate = errors.New("closing finished with error(s):\nerror 1\nerror 2")

		closersWithSleep = []func() error{
			func() error { time.Sleep(1 * time.Second); return nil },
			func() error { time.Sleep(4 * time.Second); return nil },
		}
		closersWithErrors = []func() error{
			func() error { time.Sleep(1 * time.Second); return errors.New("error 1") },
			func() error { time.Sleep(2 * time.Second); return errors.New("error 2") },
		}
		closersWithoutErrors = []func() error{
			func() error { return nil },
			func() error { return nil },
		}
	)

	tests := []struct {
		name     string
		closers  []func() error
		timeout  time.Duration
		expected error
	}{
		{
			name:     "completed without errors",
			closers:  closersWithoutErrors,
			timeout:  3 * time.Second,
			expected: nil,
		},
		{
			name:     "completed with errors",
			closers:  closersWithErrors,
			timeout:  3 * time.Second,
			expected: errorAggregate,
		},
		{
			name:     "timeout occurs",
			closers:  closersWithSleep,
			timeout:  3 * time.Second,
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

			err := c.Close(ctx)

			assert.Equal(t, test.expected, err)
		})
	}
}
