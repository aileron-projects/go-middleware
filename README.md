<!-- markdownlint-disable MD033 MD041 -->

<div align="center">

[![Release](https://img.shields.io/github/v/release/aileron-projects/go-middleware?sort=semver)](https://github.com/aileron-projects/go-middleware/releases)
[![Reference](https://pkg.go.dev/badge/github.com/aileron-projects/go-middleware.svg)](https://pkg.go.dev/github.com/aileron-projects/go-middleware)
[![DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-middleware)
[![Test](https://github.com/aileron-projects/go-middleware/actions/workflows/test.yaml/badge.svg)](https://github.com/aileron-projects/go-middleware/actions/workflows/test.yaml)

[![Insights](https://badgen.net/badge/Insights/open%2Fsource%2Finsights/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-middleware)
[![Insights](https://badgen.net/badge/Insights/OSS%2FInsight/orange)](https://ossinsight.io/analyze/aileron-projects/go-middleware)

</div>

# go-middleware

**Server-side and client-side middleware chaining library for Go.**

## Features

- HTTP server-side middleware chain
- HTTP client-side middleware chain
- Powerful chaining
- Simple

**HTTP server-side middleware pattern in go:**

```go
// ServerMiddleware is the server side middleware.
type ServerMiddleware interface {
    ServerMiddleware(next http.Handler) http.Handler
}

// ServerMiddlewareFunc type is an adapter to allow the
// use of ordinary function as server middleware.
type ServerMiddlewareFunc func(next http.Handler) http.Handler

func (f ServerMiddlewareFunc) ServerMiddleware(next http.Handler) http.Handler {
    return f(next)
}
```

**HTTP client-side middleware pattern in go:**

```go
// ClientMiddleware is the client side middleware.
type ClientMiddleware interface {
    ClientMiddleware(next http.RoundTripper) http.RoundTripper
}

// ClientMiddlewareFunc type is an adapter to allow the
// use of ordinary function as client middleware.
type ClientMiddlewareFunc func(next http.RoundTripper) http.RoundTripper

func (f ClientMiddlewareFunc) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
    return f(next)
}
```

## Usages

### Server-side middleware chaining

This example applies `logging middleware` and `timeout middleware` to the server's `hello handler`.
Server-side middlewares can be easily appended to the chain.

```go
func main() {
    chain := middleware.ServerMiddlewareChain{} // Create an empty chain.
    chain = chain.AppendFunc(logging, timeout)  // Append middlewares.
    handler := chain.TerminateFunc(hello)       // Finalize with the handler.

    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        panic(err)
    }
}

// logging prints access info.
func logging(next http.Handler) http.Handler {
    return middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Proto, r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

// timeout applies request timeout of 2 seconds.
func timeout(next http.Handler) http.Handler {
    return http.TimeoutHandler(next, time.Second, "request timeout")
}

// hello says hello after sleep. Maybe timeout.
func hello(w http.ResponseWriter, r *http.Request) {
    sleep := time.Second * time.Duration(rand.IntN(5))
    log.Println("sleeping:", sleep)
    time.Sleep(sleep)
    _, _ = w.Write([]byte("ok"))
}
```

The following methods are defined for the middleware chain.
See the [godoc](https://pkg.go.dev/github.com/aileron-projects/go-middleware) for details.

```go
// The middleware chain is a middleware itself.
func ServerMiddleware(next http.Handler) http.Handler

// Terminate the chain with the handler.
func Terminate(h http.Handler) http.Handler
func TerminateFunc(h http.HandlerFunc) http.Handler

// Add middlewares to the chain.
func Append(ms ...ServerMiddleware) ServerMiddlewareChain
func AppendFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain
func Prepend(ms ...ServerMiddleware) ServerMiddlewareChain
func PrependFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain
func InsertAt(index int, ms ...ServerMiddleware) ServerMiddlewareChain
func InsertAtFunc(index int, ms ...ServerMiddlewareFunc) ServerMiddlewareChain
func InsertBeforeAll(ms ...ServerMiddleware) ServerMiddlewareChain
func InsertBeforeAllFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain
func InsertAfterAll(ms ...ServerMiddleware) ServerMiddlewareChain
func InsertAfterAllFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain
func InsertAll(ms ...ServerMiddleware) ServerMiddlewareChain
func InsertAllFunc(ms ...ServerMiddlewareFunc) ServerMiddlewareChain
```

### Client-side middleware chaining

This example applies `logging middleware` and `timeout middleware` to the http client.
Here a dummy round tripper is used as the final round tripper..
Client-side middlewares can be easily appended to the chain.

```go

func main() {
    chain := middleware.ClientMiddlewareChain{} // Create an empty chain.
    chain = chain.AppendFunc(logging, timeout)  // Append middlewares.
    roundTripper := chain.TerminateFunc(hello)  // Finalize with the handler.

    r, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/hello", nil)
    w, err := roundTripper.RoundTrip(r)
    if err != nil {
        log.Println(err)
    } else {
        body, _ := io.ReadAll(w.Body)
        log.Println(string(body))
    }
}

// logging prints request info.
func logging(next http.RoundTripper) http.RoundTripper {
    return middleware.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
        log.Println(r.Proto, r.Method, r.RequestURI)
        return next.RoundTrip(r)
    })
}

// timeout applies request timeout of 2 seconds.
func timeout(next http.RoundTripper) http.RoundTripper {
    return middleware.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()
        r = r.WithContext(ctx)
        return next.RoundTrip(r)
    })
}

// hello returns dummy response.
func hello(r *http.Request) (*http.Response, error) {
    sleep := time.Second * time.Duration(rand.IntN(5))
    log.Println("sleeping:", sleep)
    select {
    case <-r.Context().Done():
        return nil, r.Context().Err()
    case <-time.After(sleep):
        return &http.Response{
            StatusCode: http.StatusOK,
            Body:       io.NopCloser(strings.NewReader("hello")),
        }, nil
    }
}
```

The following methods are defined for the middleware chain.
See the [godoc](https://pkg.go.dev/github.com/aileron-projects/go-middleware) for details.

```go
// The middleware chain is a middleware itself.
func ClientMiddleware(next http.RoundTripper) http.RoundTripper

// Terminate the chain with the round tripper.
func Terminate(rt http.RoundTripper) http.RoundTripper
func TerminateFunc(rt RoundTripperFunc) http.RoundTripper

// Add middlewares to the chain.
func Append(ms ...ClientMiddleware) ClientMiddlewareChain
func AppendFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain
func Prepend(ms ...ClientMiddleware) ClientMiddlewareChain
func PrependFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain
func InsertAt(index int, ms ...ClientMiddleware) ClientMiddlewareChain
func InsertAtFunc(index int, ms ...ClientMiddlewareFunc) ClientMiddlewareChain
func InsertBeforeAll(ms ...ClientMiddleware) ClientMiddlewareChain
func InsertBeforeAllFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain
func InsertAfterAll(ms ...ClientMiddleware) ClientMiddlewareChain
func InsertAfterAllFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain
func InsertAll(ms ...ClientMiddleware) ClientMiddlewareChain
func InsertAllFunc(ms ...ClientMiddlewareFunc) ClientMiddlewareChain
```

## Docs & Examples

- GoDoc: <https://pkg.go.dev/github.com/aileron-projects/go-middleware>
- Examples:
  - Server-side middleware chaining: [examples/server-middleware/](./examples/server-middleware/)
  - Client-side middleware chaining: [examples/client-middleware/](./examples/client-middleware/)

## References

- <https://github.com/justinas/alice>
- <https://pkg.go.dev/golang.org/x/pkgsite/internal/middleware>
- [Go Wiki: LearnServerProgramming](https://go.dev/wiki/LearnServerProgramming)
