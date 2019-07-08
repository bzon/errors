// Package errors provides a drop-in replacement for github.com/pkg/errors:
//
// New, Errorf, Wrap, Wrapf, and Cause
//
// The errors will be decorated with the Error interface.
//
// 	type Error interface {
// 		error
// 		Tracer
// 	}
//
// 	type Tracer interface {
// 		SourceLocation() SourceLocation
// 		TraceContext() TraceContext
// 		SetTraceContext(trace.SpanContext)
// 		SetSourceLocation(depth int)
// 	}
//
// To add tracing context to OpenCensus trace.Span use the trace wrappers:
// NewT, WrapT ...
//
// If you are creating an error wrapper on top of the New, NewT, etc,
// use the Caller and provide it with the specific caller depth:
// NewCaller, NewWrapper ...
package errors

import (
	"runtime"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
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
}

// NewSourceLocation creates a SourceLocation using stdlib runtime.Caller.
func NewSourceLocation(depth int) SourceLocation {
	function, file, line, _ := runtime.Caller(depth)
	return SourceLocation{runtime.FuncForPC(function).Name(), file, line}
}

// Error represents an error with tracing.
type Error interface {
	error
	Tracer
}

type errorContext struct {
	err            error
	sourceLocation SourceLocation
	traceContext   TraceContext
}

type customError struct {
	originalError error
	errorContext
}

func (e *errorContext) Error() string {
	return e.err.Error()
}

// TODO: adding this does not implement cause properly.
func (e *errorContext) Cause() error {
	return e.err
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

// NewCallerf wraps errors.Errorf with a specified caller depth.
func NewCallerf(depth int, m string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Errorf(m, args...),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// NewCallerfT wraps errors.Errorf with a specified caller depth and a span trace context.
func NewCallerfT(depth int, span *trace.Span, m string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Errorf(m, args...),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

// WrapCaller wraps errors.Wrap with a specified caller depth.
func WrapCaller(depth int, e error, m string) error {
	err := &errorContext{
		err:            errors.Wrap(e, m),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// WrapCallerT wraps errors.Wrap with a specified caller depth with a span trace context.
func WrapCallerT(depth int, span *trace.Span, e error, m string) error {
	err := &errorContext{
		err:            errors.Wrap(e, m),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

// WrapCallerf wraps errors.Wrapf with a specified caller depth.
func WrapCallerf(depth int, e error, m string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Wrapf(e, m, args...),
		sourceLocation: NewSourceLocation(depth),
	}
	return err
}

// WrapCallerfT wraps errors.Wrapf with a specified caller depth with a span trace context.
func WrapCallerfT(depth int, span *trace.Span, e error, m string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Wrapf(e, m, args...),
		sourceLocation: NewSourceLocation(depth),
	}
	return annotate(err, span)
}

/// Drop-in replacement for pkg/errors.

// New wraps errors.New.
func New(m string) error {
	err := &errorContext{
		err:            errors.New(m),
		sourceLocation: NewSourceLocation(2),
	}
	return err
}

// NewT wraps errors.New with a span trace context.
func NewT(span *trace.Span, m string) error {
	err := &errorContext{
		err:            errors.New(m),
		sourceLocation: NewSourceLocation(2),
	}
	return annotate(err, span)
}

// Errorf wraps errors.Errorf.
func Errorf(m string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Errorf(m, args...),
		sourceLocation: NewSourceLocation(2),
	}
	return err
}

// ErrorfT wraps errors.Errorf with a span trace context.
func ErrorfT(span *trace.Span, m string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Errorf(m, args...),
		sourceLocation: NewSourceLocation(2),
	}
	return annotate(err, span)
}

// Wrap wraps errors.Wrap.
func Wrap(e error, msg string) error {
	err := &errorContext{
		err:            errors.Wrap(e, msg),
		sourceLocation: NewSourceLocation(2),
	}
	return err
}

// WrapT wraps errors.Wrap with a span trace context.
func WrapT(span *trace.Span, e error, msg string) error {
	err := &errorContext{
		err:            errors.Wrap(e, msg),
		sourceLocation: NewSourceLocation(2),
	}
	return annotate(err, span)
}

// Wrapf wraps errors.Wrapf.
func Wrapf(e error, msg string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Wrapf(e, msg, args...),
		sourceLocation: NewSourceLocation(2),
	}
	return err
}

// WrapfT wraps errors.Wrapf with a trace context.
func WrapfT(span *trace.Span, e error, msg string, args ...interface{}) error {
	err := &errorContext{
		err:            errors.Wrapf(e, msg, args...),
		sourceLocation: NewSourceLocation(2),
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
		},
		"Error: "+e.Error(),
	)
	return e
}

// Cause wraps errors.Cause.
func Cause(e error) error {
	return errors.Cause(e)
}
