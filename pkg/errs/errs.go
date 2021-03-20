package errs

import (
	"fmt"
	"runtime"
)

// StackError wrap error with stack
type StackError struct {
	Err   error
	Stack Stack `json:"-"`
}

func (s StackError) Error() string {
	return fmt.Sprintf("%v ( %v )", s.Err, s.Stack)
}

// NewStack adds stacktrace to an error.
func NewStack(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*StackError); ok {
		return err
	}
	return &StackError{Err: err, Stack: StackTrace(3)}
}

// Newf returns new error with stack
func Newf(format string, args ...interface{}) error {
	return &StackError{
		Err:   fmt.Errorf(format, args...),
		Stack: StackTrace(3),
	}
}

// Cause extract error from stack
func Cause(err error) error {
	if err == nil {
		return nil
	}
	stackError, ok := err.(*StackError)
	if ok {
		return stackError.Err
	}
	return err
}

// Frame is a single line of executed Code in a Stack.
type Frame struct {
	Filename string `json:"filename"`
	Method   string `json:"method"`
	Line     int    `json:"lineno"`
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d [%s]", f.Filename, f.Line, f.Method)
}

// Stack represents a stacktrace as a slice of Frames.
type Stack []Frame

func (s Stack) String() string {
	if len(s) == 0 {
		return "no stack"
	}
	res := ""
	for _, value := range s {
		res += value.String() + "\n"
	}
	return res
}

// StackTrace builds a full stacktrace for the current execution location.
func StackTrace(skip int) Stack {
	pc := make([]uintptr, 50)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]

	frames := runtime.CallersFrames(pc)
	stack := make(Stack, 0, n)
	for {
		frame, more := frames.Next()
		stack = append(stack, Frame{Filename: frame.File, Method: frame.Function, Line: frame.Line})

		if !more {
			break
		}
	}

	return stack
}

type MultiError []error

func (m MultiError) Error() string {
	if len(m) == 0 {
		panic("empty MultiError call!")
	}
	res := "Errors list \n"
	for i, err := range m {
		res += fmt.Sprintf("%d) %s \n", i+1, err.Error())
	}
	return res
}
