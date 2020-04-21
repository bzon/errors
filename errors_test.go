package errors_test

import (
	"context"
	"fmt"

	stderr "errors"

	"github.com/bzon/errors"
	"go.opencensus.io/trace"
)

type customError struct {
	Err string
}

func (e *customError) Error() string {
	return e.Err
}

var errSentinel = stderr.New("sentinel error")

func ExampleAs() {
	err := &customError{Err: "err"}
	wrappedErr := errors.Wrap(err, "wrapped")

	var e *customError

	if errors.As(wrappedErr, &e) {
		fmt.Println(e.Err)
	}

	// Output:
	// err
}

func ExampleIs() {
	err := errSentinel
	wrappedErr := errors.Wrap(err, "wrapped")
	if errors.Is(wrappedErr, errSentinel) {
		fmt.Println("wrapped error is a sentinel error")
	}

	// Output:
	// wrapped error is a sentinel error
}

func ExampleNew() {
	err := errors.New("a")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// a
	// github.com/bzon/errors_test.ExampleNew
}

func ExampleErrorf() {
	err := errors.Errorf("a")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// a
	// github.com/bzon/errors_test.ExampleErrorf
}

func ExampleWrap() {
	err := errors.New("a")
	err = errors.Wrap(err, "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// b: a
	// github.com/bzon/errors_test.ExampleWrap
}

func ExampleWrapf() {
	err := errors.New("a")
	err = errors.Wrapf(err, "test %s", "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// test b: a
	// github.com/bzon/errors_test.ExampleWrapf
}

// func ExampleCause() {
// 	err := errors.New("a")
// 	err = errors.Wrap(err, "b")
// 	cause := errors.Cause(err)
// 	fmt.Println(cause)

// 	// Output: a
// }

func ExampleNewT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.NewT(span, "a")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// a
	// github.com/bzon/errors_test.ExampleNewT
}

func ExampleErrorfT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.ErrorfT(span, "a")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// a
	// github.com/bzon/errors_test.ExampleErrorfT
}

func ExampleWrapT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.New("a")
	err = errors.WrapT(span, err, "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// b: a
	// github.com/bzon/errors_test.ExampleWrapT
}

func ExampleWrapfT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.New("a")
	err = errors.WrapfT(span, err, "test %s", "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// test b: a
	// github.com/bzon/errors_test.ExampleWrapfT
}

func callFoo() error {
	return errors.NewCaller(2, "a")
}

func ExampleNewCaller() {
	// func callFoo() error {
	// 	return errors.NewCaller(2, "a")
	// }
	err := callFoo()
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// github.com/bzon/errors_test.callFoo
}

func ExampleNewCallerf() {
	err := errors.NewCallerf(2, "test %s", "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// test b
	// github.com/bzon/errors_test.ExampleNewCallerf
}

func ExampleNewCallerT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.NewCallerT(2, span, "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// b
	// github.com/bzon/errors_test.ExampleNewCallerT
}

func ExampleNewCallerfT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.NewCallerfT(2, span, "test %s", "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// test b
	// github.com/bzon/errors_test.ExampleNewCallerfT
}

func ExampleWrapCaller() {
	err := errors.New("a")
	err = errors.WrapCaller(2, err, "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// b: a
	// github.com/bzon/errors_test.ExampleWrapCaller
}

func ExampleWrapCallerf() {
	err := errors.New("a")
	err = errors.WrapCallerf(2, err, "test %s", "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// test b: a
	// github.com/bzon/errors_test.ExampleWrapCallerf
}

func ExampleWrapCallerT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.New("a")
	err = errors.WrapCallerT(2, span, err, "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// b: a
	// github.com/bzon/errors_test.ExampleWrapCallerT
}

func ExampleWrapCallerfT() {
	_, span := trace.StartSpan(context.Background(), "foo")
	defer span.End()

	err := errors.New("a")
	err = errors.WrapCallerfT(2, span, err, "test %s", "b")
	fmt.Println(err)
	e := err.(errors.ErrorTracer)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// test b: a
	// github.com/bzon/errors_test.ExampleWrapCallerfT
}

func ExampleSourceLocation() {
	err := errors.New("a")
	e := err.(errors.ErrorTracer)
	e.SetSourceLocation(1)
	fmt.Println(e.SourceLocation().Function)

	// Output:
	// github.com/bzon/errors.(*errorContext).SetSourceLocation
}

func ExampleTraceContext() {
	err := errors.New("a")
	e := err.(errors.ErrorTracer)
	e.SetTraceContext(trace.SpanContext{
		TraceID: [16]byte{'a', 'b', 'c'},
		SpanID:  [8]byte{'d', 'e', 'f'},
	})
	fmt.Println(e.TraceContext().TraceID)
	fmt.Println(e.TraceContext().SpanID)

	// Output:
	// 61626300000000000000000000000000
	// 6465660000000000
}
