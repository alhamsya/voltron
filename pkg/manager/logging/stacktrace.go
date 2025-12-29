package logging

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

func MarshalStack(err error) any {
	const (
		LimitStack = 5
	)

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	st, ok := err.(stackTracer)
	if !ok {
		return nil
	}

	stack := st.StackTrace()

	limit := LimitStack
	if len(stack) < limit {
		limit = len(stack)
	}

	frames := make([]map[string]string, 0, limit)
	for i := 0; i < limit; i++ {
		f := stack[i] // errors.Frame

		file := fmt.Sprintf("%s", f) // full path file
		frames = append(frames, map[string]string{
			"func":   fmt.Sprintf("%n", f), // function name
			"source": filepath.Base(file),  // file name only (optional)
			"line":   fmt.Sprintf("%d", f), // line number
		})
	}

	return frames
}
