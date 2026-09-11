package middleware

import (
	"net/http"
	"slices"
)

var (
	_ ClientMiddleware = &ClientMiddlewareChain{}
	_ ClientMiddleware = ClientMiddlewareFunc(nil)
)

// NewRoundTripper returns a http round tripper with client-side middlewares.
func NewRoundTripper(rt http.RoundTripper, ms ...ClientMiddleware) http.RoundTripper {
	return ClientMiddlewareChain(ms).Terminate(rt)
}

// ClientMiddlewareChain is the client-side middleware chain.
// ClientMiddlewareChain is a [ClientMiddleware] itself.
type ClientMiddlewareChain []ClientMiddleware

// ClientMiddleware implements [ClientMiddleware].
func (c ClientMiddlewareChain) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
	return c.Terminate(next)
}

// Terminate returns a round tripper with the middleware chain terminated by rt.
func (c ClientMiddlewareChain) Terminate(rt http.RoundTripper) http.RoundTripper {
	for _, m := range slices.Backward(c) {
		rt = m.ClientMiddleware(rt)
	}
	return rt
}

// TerminateFunc returns a round tripper with the middleware chain terminated by rt.
func (c ClientMiddlewareChain) TerminateFunc(rt RoundTripperFunc) http.RoundTripper {
	return c.Terminate(RoundTripperFunc(rt))
}

// Append appends middlewares to the chain.
// Given middlewares are added to the end of the chain.
func (c ClientMiddlewareChain) Append(ms ...ClientMiddleware) ClientMiddlewareChain {
	return append(c, ms...)
}

// AppendFunc appends middleware functions to the chain.
// Given middlewares are added to the end of the chain.
func (c ClientMiddlewareChain) AppendFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain {
	return append(c, toClientMiddleware(ms)...)
}

// Prepend prepends middlewares to the chain.
// Given middlewares are added to the front of the chain.
func (c ClientMiddlewareChain) Prepend(ms ...ClientMiddleware) ClientMiddlewareChain {
	return append(ms, c...)
}

// PrependFunc prepends middleware functions to the chain.
// Given middlewares are added to the front of the chain.
func (c ClientMiddlewareChain) PrependFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain {
	return append(toClientMiddleware(ms), c...)
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
func (c ClientMiddlewareChain) InsertAt(index int, ms ...ClientMiddleware) ClientMiddlewareChain {
	return insertAt(index, c, ms)
}

// InsertAtFunc inserts middleware functions at the position of index.
func (c ClientMiddlewareChain) InsertAtFunc(index int, ms ...ClientMiddlewareFunc) ClientMiddlewareChain {
	return insertAt(index, c, toClientMiddleware(ms))
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
func (c ClientMiddlewareChain) InsertAll(ms ...ClientMiddleware) ClientMiddlewareChain {
	return insertAll(c, ms)
}

// InsertAllFunc inserts given middleware functions between each middlewares in the chain.
func (c ClientMiddlewareChain) InsertAllFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain {
	return insertAll(c, toClientMiddleware(ms))
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
func (c ClientMiddlewareChain) InsertBeforeAll(ms ...ClientMiddleware) ClientMiddlewareChain {
	return insertBeforeAll(c, ms)
}

// InsertBeforeAllFunc inserts given middleware functions before all middlewares in the chain.
func (c ClientMiddlewareChain) InsertBeforeAllFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain {
	return insertBeforeAll(c, toClientMiddleware(ms))
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
func (c ClientMiddlewareChain) InsertAfterAll(ms ...ClientMiddleware) ClientMiddlewareChain {
	return insertAfterAll(c, ms)
}

// InsertAfterAllFunc inserts given middleware functions after all middlewares in the chain.
func (c ClientMiddlewareChain) InsertAfterAllFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain {
	return insertAfterAll(c, toClientMiddleware(ms))
}

func toClientMiddleware(fs []ClientMiddlewareFunc) []ClientMiddleware {
	ms := make([]ClientMiddleware, 0, len(fs))
	for _, f := range fs {
		ms = append(ms, f)
	}
	return ms
}
