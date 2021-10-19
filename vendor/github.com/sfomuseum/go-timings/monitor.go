package timings

import (
	"context"
	"io"
)

// type Monitor provides a common interface for timings-based monitors
type Monitor interface {
	// The Start method starts the monitor dispatching notifications to an io.Writer instance
	Start(context.Context, io.Writer) error
	// The Stop method will stop monitoring
	Stop(context.Context) error
	// The Signal method will dispatch messages to the monitoring process
	Signal(context.Context, ...interface{}) error
}
