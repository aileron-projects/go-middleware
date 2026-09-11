package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/aileron-projects/go-middleware"
)

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
