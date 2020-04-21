package errors

import (
	"fmt"
	"runtime"

	"errors"

	"go.opencensus.io/trace"
)

// Overwrite these values during build via -ldflags.
var (
	// VERSION is the app-global version.
	VERSION = "UNKNOWN"

	// COMMIT is the app-global commit id.
	COMMIT = "UNKNOWN"

	// BRANCH is the app-global git branch.
	BRANCH = "UNKNOWN"
)

const (
	wrappedFunctionCallDepth = 2
)

// Tracer represents an error that has TraceContext and SourceLocation.
type Tracer interface {
	SourceLocation() SourceLocation
	TraceContext() TraceContext
	SetTraceContext(trace.SpanContext)
	SetSourceLocation(depth int)
}

// TraceContext is used to provide a tracing context to an object for logging purposes.
// This is helpful for Developers to link Stackdriver traces to Stackdriver logs.
// See https://cloud.google.com/trace/docs/viewing-details.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry.
type TraceContext struct {
	TraceID string `json:"trace"`
	SpanID  string `json:"spanId"`
}

// SourceLocation provides the information where the actual error happened in the code.
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
type SourceLocation struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Version  string `json:"version,omitempty"`
	Commit   string `json:"commit,omitempty"`
	Branch   string `json:"branch,omitempty"`
}

// NewSourceLocation creates a SourceLocation using stdlib runtime.Caller.
func NewSourceLocation(depth int) SourceLocation {
	function, file, line, _ := runtime.Caller(depth)
	return SourceLocation{
		runtime.FuncForPC(function).Name(), file, line, VERSION, COMMIT, BRANCH,
	}
}

// As is a drop-in replacement for errors.As method.
func As(target error, dest interface{}) bool {
	return errors.As(target, dest)
}

// Is is a drop-in replacement for errors.Is method.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// ErrorTracer represents an error with tracing.
type ErrorTracer interface {
	error
	// Unwrap is used to implement the Go 1.13 error wrapping technique.
	// This makes the implementation of ErrorTracer to make use of errors.Is and errors.As.
	Unwrap() error
	Tracer
}

// Compile time implementation check.
var _ ErrorTracer = &errorContext{}

type errorContext struct {
	err            error
	sourceLocation SourceLocation
	traceContext   TraceContext
}

func (e *errorContext) Unwrap() error {
	return e.err
}

func (e *errorContext) Error() string {
	return e.err.Error()
}

func (e *errorContext) SourceLocation() SourceLocation {
	return e.sourceLocation
}

func (e *errorContext) SetSourceLocation(depth int) {
	e.sourceLocation = NewSourceLocation(depth)
}

func (e *errorContext) TraceContext() TraceContext {
	return e.traceContext
}

func (e *errorContext) SetTraceContext(t trace.SpanContext) {
	e.traceContext = TraceContext{
		TraceID: t.TraceID.String(),
		SpanID:  t.SpanID.String(),
	}
}

// NewCaller wraps errors.New with a specified caller depth.
func NewCaller(depth int, m string) error {
	err := &errorContext{
		err:            errors.New(m),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// NewCallerT wraps errors.New with a specified caller depth and a span trace context.
func NewCallerT(depth int, span *trace.Span, m string) error {
	err := &errorContext{
		err:            errors.New(m),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

// NewCallerf wraps fmt.Errorf with a specified caller depth.
func NewCallerf(depth int, m string, args ...interface{}) error {
	err := &errorContext{
		err:            fmt.Errorf(m, args...),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// NewCallerfT wraps fmt.Errorf with a specified caller depth and a span trace context.
func NewCallerfT(depth int, span *trace.Span, m string, args ...interface{}) error {
	err := &errorContext{
		err:            fmt.Errorf(m, args...),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

// WrapCaller wraps fmt.Errorf with a specified caller depth.
func WrapCaller(depth int, e error, m string) error {
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// WrapCallerT wraps fmt.Errorf with a specified caller depth with a span trace context.
func WrapCallerT(depth int, span *trace.Span, e error, m string) error {
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

// WrapCallerf wraps fmt.Errorf with a specified caller depth.
func WrapCallerf(depth int, e error, format string, args ...interface{}) error {
	m := fmt.Sprintf(format, args...)
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// WrapCallerfT wraps fmt.Errorf with a specified caller depth with a span trace context.
func WrapCallerfT(depth int, span *trace.Span, e error, format string, args ...interface{}) error {
	m := fmt.Sprintf(format, args...)
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

// New is the drop-in replacement for errors.New.
func New(m string) error {
	err := &errorContext{
		err:            errors.New(m),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return err
}

// NewT wraps errors.New with a span trace context.
func NewT(span *trace.Span, m string) error {
	err := &errorContext{
		err:            errors.New(m),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return annotate(err, span)
}

// Errorf wraps fmt.Errorf.
func Errorf(m string, args ...interface{}) error {
	err := &errorContext{
		err:            fmt.Errorf(m, args...),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return err
}

// ErrorfT wraps fmt.Errorf with a span trace context.
func ErrorfT(span *trace.Span, m string, args ...interface{}) error {
	err := &errorContext{
		err:            fmt.Errorf(m, args...),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return annotate(err, span)
}

// Wrap wraps an error fmt.Errorf with `%w` without formatting.
func Wrap(e error, m string) error {
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return err
}

// WrapT wraps an error with a span trace context.
func WrapT(span *trace.Span, e error, m string) error {
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return annotate(err, span)
}

// Wrap wraps fmt.Errorf with `%w` with formatting.
func Wrapf(e error, f string, args ...interface{}) error {
	m := fmt.Sprintf(f, args...)
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return err
}

// WrapfT is Wrapf with a trace context.
func WrapfT(span *trace.Span, e error, f string, args ...interface{}) error {
	m := fmt.Sprintf(f, args...)
	err := &errorContext{
		err:            fmt.Errorf("%s: %w", m, e),
		sourceLocation: NewSourceLocation(wrappedFunctionCallDepth),
	}
	return annotate(err, span)
}

func annotate(e *errorContext, span *trace.Span) error {
	if span == nil {
		return e
	}

	// Add the trace ID and span ID.
	ctx := span.SpanContext()
	e.traceContext = TraceContext{
		TraceID: ctx.TraceID.String(),
		SpanID:  ctx.SpanID.String(),
	}

	// Add OpenCensus span annotation.
	src := e.SourceLocation()
	span.Annotate(
		[]trace.Attribute{
			trace.StringAttribute("function", src.Function),
			trace.StringAttribute("file", src.File),
			trace.Int64Attribute("line", int64(src.Line)),
			trace.StringAttribute("version", src.Version),
			trace.StringAttribute("commit", src.Commit),
			trace.StringAttribute("branch", src.Branch),
		},
		"Error: "+e.Error(),
	)

	// Generic error
	span.SetStatus(trace.Status{
		Code: trace.StatusCodeUnknown,
	})
	return e
}
