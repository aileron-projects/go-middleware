package middleware

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestHandlerFunc(t *testing.T) {
	t.Parallel()
	called := false
	h := HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	h.ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, called)
}

func TestServerMiddlewareFunc(t *testing.T) {
	t.Parallel()
	called := false
	m := ServerMiddlewareFunc(func(next http.Handler) http.Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})
	})
	h := m.ServerMiddleware(nil)
	h.ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, called)
}

func TestRoundTripperFunc(t *testing.T) {
	t.Parallel()
	called := false
	h := RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		return nil, nil
	})
	h.RoundTrip(nil)
	tester.AssertEqual(t, true, called)
}

func TestClientMiddlewareFunc(t *testing.T) {
	t.Parallel()
	called := false
	m := ClientMiddlewareFunc(func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return nil, nil
		})
	})
	rt := m.ClientMiddleware(nil)
	rt.RoundTrip(nil)
	tester.AssertEqual(t, true, called)
}

func TestInsertAt(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		index  int
		target []string
		elems  []string
		want   []string
	}{
		{-1, []string{}, []string{}, []string{}},
		{0, []string{}, []string{}, []string{}},
		{1, []string{}, []string{}, []string{}},
		{-1, []string{"A"}, []string{"a"}, []string{"a", "A"}},
		{0, []string{"A"}, []string{"a"}, []string{"a", "A"}},
		{1, []string{"A"}, []string{"a"}, []string{"A", "a"}},
		{2, []string{"A"}, []string{"a"}, []string{"A", "a"}},
		{-1, []string{"A", "B"}, []string{"a"}, []string{"a", "A", "B"}},
		{0, []string{"A", "B"}, []string{"a"}, []string{"a", "A", "B"}},
		{1, []string{"A", "B"}, []string{"a"}, []string{"A", "a", "B"}},
		{2, []string{"A", "B"}, []string{"a"}, []string{"A", "B", "a"}},
		{3, []string{"A", "B"}, []string{"a"}, []string{"A", "B", "a"}},
		{-1, []string{"A", "B"}, []string{"a", "b"}, []string{"a", "b", "A", "B"}},
		{0, []string{"A", "B"}, []string{"a", "b"}, []string{"a", "b", "A", "B"}},
		{1, []string{"A", "B"}, []string{"a", "b"}, []string{"A", "a", "b", "B"}},
		{2, []string{"A", "B"}, []string{"a", "b"}, []string{"A", "B", "a", "b"}},
		{3, []string{"A", "B"}, []string{"a", "b"}, []string{"A", "B", "a", "b"}},
	}
	for i, tc := range testCases {
		t.Run("case:"+strconv.Itoa(i), func(t *testing.T) {
			got := insertAt(tc.index, tc.target, tc.elems)
			tester.AssertDeepEqual(t, tc.want, got)
		})
	}
}

func TestInsertAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		target []string
		elems  []string
		want   []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{}, []string{"a"}, []string{}},
		{[]string{"A"}, []string{}, []string{"A"}},
		{[]string{"A"}, []string{"a"}, []string{"a", "A", "a"}},
		{[]string{"A"}, []string{"a", "b"}, []string{"a", "b", "A", "a", "b"}},
		{[]string{"A", "B"}, []string{"a", "b"}, []string{"a", "b", "A", "a", "b", "B", "a", "b"}},
	}
	for i, tc := range testCases {
		t.Run("case:"+strconv.Itoa(i), func(t *testing.T) {
			got := insertAll(tc.target, tc.elems)
			tester.AssertDeepEqual(t, tc.want, got)
		})
	}
}

func TestInsertBeforeAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		target []string
		elems  []string
		want   []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{}, []string{"a"}, []string{}},
		{[]string{"A"}, []string{}, []string{"A"}},
		{[]string{"A"}, []string{"a"}, []string{"a", "A"}},
		{[]string{"A"}, []string{"a", "b"}, []string{"a", "b", "A"}},
		{[]string{"A", "B"}, []string{"a", "b"}, []string{"a", "b", "A", "a", "b", "B"}},
	}
	for i, tc := range testCases {
		t.Run("case:"+strconv.Itoa(i), func(t *testing.T) {
			got := insertBeforeAll(tc.target, tc.elems)
			tester.AssertDeepEqual(t, tc.want, got)
		})
	}
}

func TestInsertAfterAll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		target []string
		elems  []string
		want   []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{}, []string{"a"}, []string{}},
		{[]string{"A"}, []string{}, []string{"A"}},
		{[]string{"A"}, []string{"a"}, []string{"A", "a"}},
		{[]string{"A"}, []string{"a", "b"}, []string{"A", "a", "b"}},
		{[]string{"A", "B"}, []string{"a", "b"}, []string{"A", "a", "b", "B", "a", "b"}},
	}
	for i, tc := range testCases {
		t.Run("case:"+strconv.Itoa(i), func(t *testing.T) {
			got := insertAfterAll(tc.target, tc.elems)
			tester.AssertDeepEqual(t, tc.want, got)
		})
	}
}
