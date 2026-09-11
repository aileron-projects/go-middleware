package middleware_test

import (
	"net/http"
	"testing"

	"github.com/aileron-projects/go-middleware"
	"github.com/aileron-projects/go-tester"
)

type testRoundTripper struct {
	called bool
}

func (rt *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.called = true
	return nil, nil
}

type testClientMiddleware struct {
	name string
	list *[]string
}

func (m *testClientMiddleware) ClientMiddleware(next http.RoundTripper) http.RoundTripper {
	return middleware.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		*m.list = append(*m.list, m.name)
		return next.RoundTrip(r)
	})
}

func TestNewRoundTripper(t *testing.T) {
	t.Parallel()
	t.Run("no middleware", func(t *testing.T) {
		r := &testRoundTripper{}
		got := middleware.NewRoundTripper(r)
		got.RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
	})
	t.Run("single middleware", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		got := middleware.NewRoundTripper(r, m1)
		got.RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1"}, list)
	})
	t.Run("double middleware", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		m2 := &testClientMiddleware{name: "m2", list: &list}
		got := middleware.NewRoundTripper(r, m1, m2)
		got.RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2"}, list)
	})
}

func TestClientMiddlewareChain_ClientMiddleware(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := &testClientMiddleware{name: "m2", list: &list}
	m3 := &testClientMiddleware{name: "m3", list: &list}
	m4 := &testClientMiddleware{name: "m4", list: &list}

	c1 := middleware.ClientMiddlewareChain{m1, m2}
	c2 := middleware.ClientMiddlewareChain{m3, m4}
	got := c1.Append(c2)
	got.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2", "m3", "m4"}, list)
}

func TestClientMiddlewareChain_Terminate(t *testing.T) {
	t.Parallel()
	t.Run("no middleware", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		c1 := middleware.ClientMiddlewareChain{}
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{}, list)
	})
	t.Run("single middleware", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		c1 := middleware.ClientMiddlewareChain{m1}
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1"}, list)
	})
	t.Run("multiple middleware", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		m2 := &testClientMiddleware{name: "m2", list: &list}
		c1 := middleware.ClientMiddlewareChain{m1, m2}
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2"}, list)
	})
}

func TestClientMiddlewareChain_TerminateFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := (&testClientMiddleware{name: "m1", list: &list})
	m2 := (&testClientMiddleware{name: "m2", list: &list})
	c1 := middleware.ClientMiddlewareChain{m1, m2}
	c1.TerminateFunc(r.RoundTrip).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2"}, list)
}

func TestClientMiddlewareChain_Append(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := &testClientMiddleware{name: "m2", list: &list}
	m3 := &testClientMiddleware{name: "m3", list: &list}
	c1 := middleware.ClientMiddlewareChain{m1}
	got := c1.Append(m2, m3)
	got.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2", "m3"}, list)
}

func TestClientMiddlewareChain_AppendFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := (&testClientMiddleware{name: "m2", list: &list}).ClientMiddleware
	m3 := (&testClientMiddleware{name: "m3", list: &list}).ClientMiddleware
	c1 := middleware.ClientMiddlewareChain{m1}
	got := c1.AppendFunc(m2, m3)
	got.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2", "m3"}, list)
}

func TestClientMiddlewareChain_Prepend(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := &testClientMiddleware{name: "m2", list: &list}
	m3 := &testClientMiddleware{name: "m3", list: &list}
	c1 := middleware.ClientMiddlewareChain{m1}
	got := c1.Prepend(m2, m3)
	got.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
}

func TestClientMiddlewareChain_PrependFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := (&testClientMiddleware{name: "m2", list: &list}).ClientMiddleware
	m3 := (&testClientMiddleware{name: "m3", list: &list}).ClientMiddleware
	c1 := middleware.ClientMiddlewareChain{m1}
	got := c1.PrependFunc(m2, m3)
	got.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
}

func TestClientMiddlewareChain_InsertAt(t *testing.T) {
	t.Parallel()
	t.Run("-1", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		m2 := &testClientMiddleware{name: "m2", list: &list}
		m3 := &testClientMiddleware{name: "m3", list: &list}
		c1 := middleware.ClientMiddlewareChain{}
		c1 = c1.InsertAt(-1, m1)
		c1 = c1.InsertAt(-1, m2, m3)
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
	})
	t.Run("0", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		m2 := &testClientMiddleware{name: "m2", list: &list}
		m3 := &testClientMiddleware{name: "m3", list: &list}
		c1 := middleware.ClientMiddlewareChain{}
		c1 = c1.InsertAt(0, m1)
		c1 = c1.InsertAt(0, m2, m3)
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
	})
	t.Run("1", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		m2 := &testClientMiddleware{name: "m2", list: &list}
		m3 := &testClientMiddleware{name: "m3", list: &list}
		m4 := &testClientMiddleware{name: "m4", list: &list}
		c1 := middleware.ClientMiddlewareChain{}
		c1 = c1.InsertAt(1, m1)
		c1 = c1.InsertAt(1, m2, m3)
		c1 = c1.InsertAt(1, m4)
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1", "m4", "m2", "m3"}, list)
	})
	t.Run("2", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := &testClientMiddleware{name: "m1", list: &list}
		m2 := &testClientMiddleware{name: "m2", list: &list}
		m3 := &testClientMiddleware{name: "m3", list: &list}
		m4 := &testClientMiddleware{name: "m4", list: &list}
		c1 := middleware.ClientMiddlewareChain{}
		c1 = c1.InsertAt(2, m1)
		c1 = c1.InsertAt(2, m2, m3)
		c1 = c1.InsertAt(2, m4)
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2", "m4", "m3"}, list)
	})
	t.Run("func", func(t *testing.T) {
		list := []string{}
		r := &testRoundTripper{}
		m1 := (&testClientMiddleware{name: "m1", list: &list}).ClientMiddleware
		m2 := (&testClientMiddleware{name: "m2", list: &list}).ClientMiddleware
		m3 := (&testClientMiddleware{name: "m3", list: &list}).ClientMiddleware
		m4 := (&testClientMiddleware{name: "m4", list: &list}).ClientMiddleware
		c1 := middleware.ClientMiddlewareChain{}
		c1 = c1.InsertAtFunc(2, m1)
		c1 = c1.InsertAtFunc(2, m2, m3)
		c1 = c1.InsertAtFunc(2, m4)
		c1.Terminate(r).RoundTrip(nil)
		tester.AssertEqual(t, true, r.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2", "m4", "m3"}, list)
	})
}

func TestClientMiddlewareChain_InsertAll(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := &testClientMiddleware{name: "m2", list: &list}
	m3 := &testClientMiddleware{name: "m3", list: &list}
	m4 := &testClientMiddleware{name: "m4", list: &list}
	c1 := middleware.ClientMiddlewareChain{m1}
	c1 = c1.InsertAll(m2)
	c1 = c1.InsertAll(m3, m4)
	c1.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1", "m3", "m4", "m2", "m3", "m4"}, list)
}

func TestClientMiddlewareChain_InsertAllFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := (&testClientMiddleware{name: "m2", list: &list}).ClientMiddleware
	m3 := (&testClientMiddleware{name: "m3", list: &list}).ClientMiddleware
	m4 := (&testClientMiddleware{name: "m4", list: &list}).ClientMiddleware
	c1 := middleware.ClientMiddlewareChain{m1}
	c1 = c1.InsertAllFunc(m2)
	c1 = c1.InsertAllFunc(m3, m4)
	c1.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1", "m3", "m4", "m2", "m3", "m4"}, list)
}

func TestClientMiddlewareChain_InsertBeforeAll(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := &testClientMiddleware{name: "m2", list: &list}
	m3 := &testClientMiddleware{name: "m3", list: &list}
	m4 := &testClientMiddleware{name: "m4", list: &list}
	c1 := middleware.ClientMiddlewareChain{m1}
	c1 = c1.InsertBeforeAll(m2)
	c1 = c1.InsertBeforeAll(m3, m4)
	c1.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1"}, list)
}

func TestClientMiddlewareChain_InsertBeforeAllFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := (&testClientMiddleware{name: "m2", list: &list}).ClientMiddleware
	m3 := (&testClientMiddleware{name: "m3", list: &list}).ClientMiddleware
	m4 := (&testClientMiddleware{name: "m4", list: &list}).ClientMiddleware
	c1 := middleware.ClientMiddlewareChain{m1}
	c1 = c1.InsertBeforeAllFunc(m2)
	c1 = c1.InsertBeforeAllFunc(m3, m4)
	c1.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1"}, list)
}

func TestClientMiddlewareChain_InsertAfterAll(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := &testClientMiddleware{name: "m2", list: &list}
	m3 := &testClientMiddleware{name: "m3", list: &list}
	m4 := &testClientMiddleware{name: "m4", list: &list}
	c1 := middleware.ClientMiddlewareChain{m1}
	c1 = c1.InsertAfterAll(m2)
	c1 = c1.InsertAfterAll(m3, m4)
	c1.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m1", "m3", "m4", "m2", "m3", "m4"}, list)
}

func TestClientMiddlewareChain_InsertAfterAllFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	r := &testRoundTripper{}
	m1 := &testClientMiddleware{name: "m1", list: &list}
	m2 := (&testClientMiddleware{name: "m2", list: &list}).ClientMiddleware
	m3 := (&testClientMiddleware{name: "m3", list: &list}).ClientMiddleware
	m4 := (&testClientMiddleware{name: "m4", list: &list}).ClientMiddleware
	c1 := middleware.ClientMiddlewareChain{m1}
	c1 = c1.InsertAfterAllFunc(m2)
	c1 = c1.InsertAfterAllFunc(m3, m4)
	c1.Terminate(r).RoundTrip(nil)
	tester.AssertEqual(t, true, r.called)
	tester.AssertDeepEqual(t, []string{"m1", "m3", "m4", "m2", "m3", "m4"}, list)
}
