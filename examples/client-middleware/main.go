package main

import (
	"context"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/aileron-projects/go-middleware"
)

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
