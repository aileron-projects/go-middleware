package middleware

import (
	"net/http"
	"slices"
)

var (
	_ ServerMiddleware = &ServerMiddlewareChain{}
	_ ServerMiddleware = ServerMiddlewareFunc(nil)
)

// NewHandler returns a http handler with server-side middlewares.
func NewHandler(h http.Handler, ms ...ServerMiddleware) http.Handler {
	return ServerMiddlewareChain(ms).Terminate(h)
}

// ServerMiddlewareChain is the server-side middleware chain.
// ServerMiddlewareChain is a [ServerMiddleware] itself.
type ServerMiddlewareChain []ServerMiddleware

// ServerMiddleware implements [ServerMiddleware].
func (c ServerMiddlewareChain) ServerMiddleware(next http.Handler) http.Handler {
	return c.Terminate(next)
}

// Terminate returns a handler with the middleware chain terminated by h.
func (c ServerMiddlewareChain) Terminate(h http.Handler) http.Handler {
	for _, m := range slices.Backward(c) {
		h = m.ServerMiddleware(h)
	}
	return h
}

// TerminateFunc returns a handler with the middleware chain terminated by h.
func (c ServerMiddlewareChain) TerminateFunc(h http.HandlerFunc) http.Handler {
	return c.Terminate(HandlerFunc(h))
}

// Append appends middlewares to the chain.
// Given middlewares are added to the end of the chain.
func (c ServerMiddlewareChain) Append(ms ...ServerMiddleware) ServerMiddlewareChain {
	return append(c, ms...)
}

// AppendFunc appends middleware functions to the chain.
// Given middlewares are added to the end of the chain.
func (c ServerMiddlewareChain) AppendFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain {
	return append(c, toServerMiddleware(ms)...)
}

// Prepend prepends middlewares to the chain.
// Given middlewares are added to the front of the chain.
func (c ServerMiddlewareChain) Prepend(ms ...ServerMiddleware) ServerMiddlewareChain {
	return append(ms, c...)
}

// PrependFunc prepends middleware functions to the chain.
// Given middlewares are added to the front of the chain.
func (c ServerMiddlewareChain) PrependFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain {
	return append(toServerMiddleware(ms), c...)
}

// InsertAt inserts middlewares at the position of index.
// For index<=0, middlewares are added at the beginning of the chain.
// For index>=len(c),  middlewares are added at the end of the chain.
//
//	       ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	m.w.   │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	       ├──────┼──────┼──────┼──────┼─────┼────────┤
//	       ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index  0      1      2      3      4    n-1       n
func (c ServerMiddlewareChain) InsertAt(index int, ms ...ServerMiddleware) ServerMiddlewareChain {
	return insertAt(index, c, ms)
}

// InsertAtFunc inserts middleware functions at the position of index.
func (c ServerMiddlewareChain) InsertAtFunc(index int, ms ...ServerMiddlewareFunc) ServerMiddlewareChain {
	return insertAt(index, c, toServerMiddleware(ms))
}

// InsertAll inserts given middlewares between each middlewares in the chain.
// Middlewares are inserted at the index of 0 to n.
// Calling InsertAll to an empty chain returns an empty chain also.
//
//	insert ○      ○      ○      ○      ○      ○       ○
//	       ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	m.w.   │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	       ├──────┼──────┼──────┼──────┼─────┼────────┤
//	       ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index  0      1      2      3      4    n-1       n
func (c ServerMiddlewareChain) InsertAll(ms ...ServerMiddleware) ServerMiddlewareChain {
	return insertAll(c, ms)
}

// InsertAllFunc inserts given middleware functions between each middlewares in the chain.
func (c ServerMiddlewareChain) InsertAllFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain {
	return insertAll(c, toServerMiddleware(ms))
}

// InsertBeforeAll inserts given middlewares before all middlewares in the chain.
// Middlewares are inserted at the index of 0 to n-1.
// Calling InsertBeforeAll to an empty chain returns an empty chain also.
//
//	insert ○      ○      ○      ○      ○      ○       ×
//	       ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	m.w.   │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	       ├──────┼──────┼──────┼──────┼─────┼────────┤
//	       ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index  0      1      2      3      4    n-1       n
func (c ServerMiddlewareChain) InsertBeforeAll(ms ...ServerMiddleware) ServerMiddlewareChain {
	return insertBeforeAll(c, ms)
}

// InsertBeforeAllFunc inserts given middleware functions before all middlewares in the chain.
func (c ServerMiddlewareChain) InsertBeforeAllFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain {
	return insertBeforeAll(c, toServerMiddleware(ms))
}

// InsertAfterAll inserts given middlewares after all middlewares in the chain.
// Middlewares are inserted at the index of 1 to n.
// Calling InsertAfterAll to an empty chain returns an empty chain also.
//
//	insert ×      ○      ○      ○      ○      ○       ○
//	       ┌──────┬──────┬──────┬──────┬─────┬────────┐
//	m.w.   │ m[0] │ m[1] │ m[2] │ m[3] │ ... │ m[n-1] │
//	       ├──────┼──────┼──────┼──────┼─────┼────────┤
//	       ↑      ↑      ↑      ↑      ↑     ↑        ↑
//	index  0      1      2      3      4    n-1       n
func (c ServerMiddlewareChain) InsertAfterAll(ms ...ServerMiddleware) ServerMiddlewareChain {
	return insertAfterAll(c, ms)
}

// InsertAfterAllFunc inserts given middleware functions after all middlewares in the chain.
func (c ServerMiddlewareChain) InsertAfterAllFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain {
	return insertAfterAll(c, toServerMiddleware(ms))
}

func toServerMiddleware(fs []ServerMiddlewareFunc) []ServerMiddleware {
	ms := make([]ServerMiddleware, 0, len(fs))
	for _, f := range fs {
		ms = append(ms, f)
	}
	return ms
}
