package ctrl

import (
	"errors"
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
	sb.Grow(1 << 10)
	for n, frame := range stack.frames {
		if n > 0 {
			sb.WriteRune('\n')
		}
		fmt.Fprintf(&sb, "[#%02d][%s]@%s:%d", n, frame.name, frame.file, frame.line)
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

func (info *PanicError) Unwrap() error {
	if e, ok := info.reason.(error); ok {
		return e
	} else {
		return nil
	}
}

func (info *PanicError) reasonS() (s string) {
	switch v := info.reason.(type) {
	case string:
		s = v
	case error:
		s = v.Error()
	default:
		s = fmt.Sprintf("unexpected: %+v", v)
	}
	return
}

func (info *PanicError) Error() string {
	return info.reasonS()
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

func StackUnroll(e error) (string, bool) {
	var info *PanicError
	var stacks []CallStack
	for errors.As(e, &info) {
		stacks = append(stacks, info.stack)
		e = info.Unwrap()
	}

	if len(stacks) == 0 {
		return "no stack trace available", false
	}

	var sb strings.Builder
	sb.Grow(4 << 10)
	for i := 0; i < len(stacks); i++ {
		stack := stacks[len(stacks)-1-i]
		fmt.Fprintf(&sb, "===== callstack %d =====\n", i+1)
		sb.WriteString(stack.String())
		sb.WriteRune('\n')
	}
	sb.WriteString("=======================")
	return sb.String(), true
}
