package ctrl

import (
	"fmt"
	"runtime"
	"strings"
)

type CallStack struct {
	frames []CallFrame
}

type CallFrame struct {
	file string
	line int
	name string
}

func StackTrace() CallStack {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs) // skip runtime.Callers + current func
	raw := runtime.CallersFrames(pcs[:n])
	frames := make([]CallFrame, 0, n)
	for {
		f, more := raw.Next()
		frames = append(frames, CallFrame{
			file: f.File,
			line: f.Line,
			name: f.Function,
		})
		if !more {
			break
		}
	}
	return CallStack{frames: frames}
}

func (stack CallStack) String() string {
	var sb strings.Builder
	for n, frame := range stack.frames {
		if n > 0 {
			sb.WriteRune('\n')
		}
		fmt.Fprintf(&sb, "->[#%02d][%s]@%s:%d", n, frame.name, frame.file, frame.line)
	}
	return sb.String()
}

type PanicError struct {
	stack  CallStack
	reason any
}

func NewPanicError(stack CallStack, reason any) *PanicError {
	return &PanicError{stack: stack, reason: reason}
}

func (info *PanicError) reasonS() (s string) {
	switch v := info.reason.(type) {
	case string:
		s = v
	case error:
		s = fmt.Sprintf("error: %s", v.Error())
	default:
		s = fmt.Sprintf("unexpected: %+v", v)
	}
	return
}

func (info *PanicError) String() string {
	var sb strings.Builder
	sb.Grow(1 << 10)
	sb.WriteString(info.reasonS())
	sb.WriteRune('\n')
	sb.WriteString("=== stack trace ===")
	sb.WriteRune('\n')
	sb.WriteString(info.stack.String())
	sb.WriteRune('\n')
	sb.WriteString("===================")
	return sb.String()
}

func (info *PanicError) Error() string {
	return info.String()
}

func RunAndCatch(exec func()) (e error) {
	defer func() {
		if rec := recover(); rec != nil {
			stack := StackTrace()
			e = NewPanicError(stack, rec)
		}
	}()
	exec()
	return
}
