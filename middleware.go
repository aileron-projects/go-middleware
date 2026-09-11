package middleware

import (
	"net/http"
	"slices"
)

// HandlerFunc type is an adapter to allow the
// use of ordinary function as HTTP handlers.
// If f is a function with the appropriate signature,
// HandlerFunc(f) is [http.Handler] that calls f.
// HandlerFunc is the same as [http.HandlerFunc].
//
// Example:
//
//	var h http.Handler
//	h = HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("ok"))
//	})
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// ServerMiddleware is the server side middleware.
type ServerMiddleware interface {
	ServerMiddleware(next http.Handler) http.Handler
}

// ServerMiddlewareFunc type is an adapter to allow the
// use of ordinary function as server middleware.
// If f is a function with the appropriate signature,
// ServerMiddlewareFunc(f) is [ServerMiddleware] that calls f.
//
// Example:
//
//	var m ServerMiddleware
//	m = ServerMiddlewareFunc(func(next http.Handler) http.Handler {
//		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// Some process.
//			next.ServeHTTP(w, r)
//		})
//	})
type ServerMiddlewareFunc func(next http.Handler) http.Handler

func (f ServerMiddlewareFunc) ServerMiddleware(next http.Handler) http.Handler {
	return f(next)
}

// RoundTripperFunc type is an adapter to allow the
// use of ordinary function as HTTP round tripper.
// If f is a function with the appropriate signature,
// RoundTripperFunc(f) is [http.RoundTripper] that calls f.
//
// Example:
//
//	var r http.RoundTripper
//	r = RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
//		// Some process.
//		return http.DefaultTransport.RoundTrip(r)
//	})
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// ClientMiddleware is the client side middleware.
type ClientMiddleware interface {
	ClientMiddleware(next http.RoundTripper) http.RoundTripper
}

// ClientMiddlewareFunc type is an adapter to allow the
// use of ordinary function as client middleware.
// If f is a function with the appropriate signature,
// ClientMiddlewareFunc(f) is [ClientMiddleware] that calls f.
//
// Example:
//
//	var m ClientMiddleware
//	m = ClientMiddlewareFunc(func(next http.RoundTripper) http.RoundTripper {
//		return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
//			// Some process.
//			return next.RoundTrip(r)
//		})
//	})
type ClientMiddlewareFunc func(next http.RoundTripper) http.RoundTripper

func (f ClientMiddlewareFunc) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
	return f(next)
}

// insertAt inserts elems at the given index position of target.
// It returns target itself when the elems is empty.
func insertAt[S ~[]T, T any](index int, target, elems S) []T {
	if len(elems) == 0 {
		return target
	}
	index = min(max(0, index), len(target))
	return slices.Insert(target, index, elems...)
}

// insertAll inserts elems before and after each elements of the target.
// It returns target itself when the target or elems is empty.
func insertAll[S ~[]T, T any](target, elems S) []T {
	if len(target) == 0 || len(elems) == 0 {
		return target
	}
	n := len(target) + (len(target)+1)*len(elems)
	arr := make([]T, 0, n)
	for _, t := range target {
		arr = append(arr, elems...)
		arr = append(arr, t)
	}
	arr = append(arr, elems...)
	return arr
}

// insertBeforeAll inserts elems before each elements of the target.
// It returns target itself when the target or elems is empty.
func insertBeforeAll[S ~[]T, T any](target, elems S) []T {
	if len(target) == 0 || len(elems) == 0 {
		return target
	}
	n := len(target) * (1 + len(elems))
	arr := make([]T, 0, n)
	for _, t := range target {
		arr = append(arr, elems...)
		arr = append(arr, t)
	}
	return arr
}

// insertAfterAll inserts elems after each elements of the target.
// It returns target itself when the target or elems is empty.
func insertAfterAll[S ~[]T, T any](target, elems S) []T {
	if len(target) == 0 || len(elems) == 0 {
		return target
	}
	n := len(target) * (1 + len(elems))
	arr := make([]T, 0, n)
	for _, t := range target {
		arr = append(arr, t)
		arr = append(arr, elems...)
	}
	return arr
}
