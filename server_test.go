package middleware_test

import (
	"net/http"
	"testing"

	"github.com/aileron-projects/go-middleware"
	"github.com/aileron-projects/go-tester"
)

type testHandler struct {
	called bool
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
}

type testServerMiddleware struct {
	name string
	list *[]string
}

func (m *testServerMiddleware) ServerMiddleware(next http.Handler) http.Handler {
	return middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*m.list = append(*m.list, m.name)
		next.ServeHTTP(w, r)
	})
}

func TestNewHandler(t *testing.T) {
	t.Parallel()
	t.Run("no middleware", func(t *testing.T) {
		h := &testHandler{}
		got := middleware.NewHandler(h)
		got.ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
	})
	t.Run("single middleware", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		got := middleware.NewHandler(h, m1)
		got.ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1"}, list)
	})
	t.Run("double middleware", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		m2 := &testServerMiddleware{name: "m2", list: &list}
		got := middleware.NewHandler(h, m1, m2)
		got.ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2"}, list)
	})
}

func TestServerMiddlewareChain_ServerMiddleware(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	m3 := &testServerMiddleware{name: "m3", list: &list}
	m4 := &testServerMiddleware{name: "m4", list: &list}

	c1 := middleware.ServerMiddlewareChain{m1, m2}
	c2 := middleware.ServerMiddlewareChain{m3, m4}
	got := c1.Append(c2)
	got.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2", "m3", "m4"}, list)
}

func TestServerMiddlewareChain_Terminate(t *testing.T) {
	t.Parallel()
	t.Run("no middleware", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		c1 := middleware.ServerMiddlewareChain{}
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{}, list)
	})
	t.Run("single middleware", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		c1 := middleware.ServerMiddlewareChain{m1}
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1"}, list)
	})
	t.Run("multiple middleware", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		m2 := &testServerMiddleware{name: "m2", list: &list}
		c1 := middleware.ServerMiddlewareChain{m1, m2}
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2"}, list)
	})
}

func TestServerMiddlewareChain_TerminateFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	c1 := middleware.ServerMiddlewareChain{m1, m2}
	c1.TerminateFunc(h.ServeHTTP).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2"}, list)
}

func TestServerMiddlewareChain_Append(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	m3 := &testServerMiddleware{name: "m3", list: &list}
	c1 := middleware.ServerMiddlewareChain{m1}
	got := c1.Append(m2, m3)
	got.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2", "m3"}, list)
}

func TestServerMiddlewareChain_AppendFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := (&testServerMiddleware{name: "m2", list: &list}).ServerMiddleware
	m3 := (&testServerMiddleware{name: "m3", list: &list}).ServerMiddleware
	c1 := middleware.ServerMiddlewareChain{m1}
	got := c1.AppendFunc(m2, m3)
	got.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m1", "m2", "m3"}, list)
}

func TestServerMiddlewareChain_Prepend(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	m3 := &testServerMiddleware{name: "m3", list: &list}
	c1 := middleware.ServerMiddlewareChain{m1}
	got := c1.Prepend(m2, m3)
	got.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
}

func TestServerMiddlewareChain_PrependFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := (&testServerMiddleware{name: "m2", list: &list}).ServerMiddleware
	m3 := (&testServerMiddleware{name: "m3", list: &list}).ServerMiddleware
	c1 := middleware.ServerMiddlewareChain{m1}
	got := c1.PrependFunc(m2, m3)
	got.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
}

func TestServerMiddlewareChain_InsertAt(t *testing.T) {
	t.Parallel()
	t.Run("-1", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		m2 := &testServerMiddleware{name: "m2", list: &list}
		m3 := &testServerMiddleware{name: "m3", list: &list}
		c1 := middleware.ServerMiddlewareChain{}
		c1 = c1.InsertAt(-1, m1)
		c1 = c1.InsertAt(-1, m2, m3)
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
	})
	t.Run("0", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		m2 := &testServerMiddleware{name: "m2", list: &list}
		m3 := &testServerMiddleware{name: "m3", list: &list}
		c1 := middleware.ServerMiddlewareChain{}
		c1 = c1.InsertAt(0, m1)
		c1 = c1.InsertAt(0, m2, m3)
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m2", "m3", "m1"}, list)
	})
	t.Run("1", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		m2 := &testServerMiddleware{name: "m2", list: &list}
		m3 := &testServerMiddleware{name: "m3", list: &list}
		m4 := &testServerMiddleware{name: "m4", list: &list}
		c1 := middleware.ServerMiddlewareChain{}
		c1 = c1.InsertAt(1, m1)
		c1 = c1.InsertAt(1, m2, m3)
		c1 = c1.InsertAt(1, m4)
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1", "m4", "m2", "m3"}, list)
	})
	t.Run("2", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := &testServerMiddleware{name: "m1", list: &list}
		m2 := &testServerMiddleware{name: "m2", list: &list}
		m3 := &testServerMiddleware{name: "m3", list: &list}
		m4 := &testServerMiddleware{name: "m4", list: &list}
		c1 := middleware.ServerMiddlewareChain{}
		c1 = c1.InsertAt(2, m1)
		c1 = c1.InsertAt(2, m2, m3)
		c1 = c1.InsertAt(2, m4)
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2", "m4", "m3"}, list)
	})
	t.Run("func", func(t *testing.T) {
		list := []string{}
		h := &testHandler{}
		m1 := (&testServerMiddleware{name: "m1", list: &list}).ServerMiddleware
		m2 := (&testServerMiddleware{name: "m2", list: &list}).ServerMiddleware
		m3 := (&testServerMiddleware{name: "m3", list: &list}).ServerMiddleware
		m4 := (&testServerMiddleware{name: "m4", list: &list}).ServerMiddleware
		c1 := middleware.ServerMiddlewareChain{}
		c1 = c1.InsertAtFunc(2, m1)
		c1 = c1.InsertAtFunc(2, m2, m3)
		c1 = c1.InsertAtFunc(2, m4)
		c1.Terminate(h).ServeHTTP(nil, nil)
		tester.AssertEqual(t, true, h.called)
		tester.AssertDeepEqual(t, []string{"m1", "m2", "m4", "m3"}, list)
	})
}

func TestServerMiddlewareChain_InsertAll(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	m3 := &testServerMiddleware{name: "m3", list: &list}
	m4 := &testServerMiddleware{name: "m4", list: &list}
	c1 := middleware.ServerMiddlewareChain{m1}
	c1 = c1.InsertAll(m2)
	c1 = c1.InsertAll(m3, m4)
	c1.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1", "m3", "m4", "m2", "m3", "m4"}, list)
}

func TestServerMiddlewareChain_InsertAllFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := (&testServerMiddleware{name: "m2", list: &list}).ServerMiddleware
	m3 := (&testServerMiddleware{name: "m3", list: &list}).ServerMiddleware
	m4 := (&testServerMiddleware{name: "m4", list: &list}).ServerMiddleware
	c1 := middleware.ServerMiddlewareChain{m1}
	c1 = c1.InsertAllFunc(m2)
	c1 = c1.InsertAllFunc(m3, m4)
	c1.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1", "m3", "m4", "m2", "m3", "m4"}, list)
}

func TestServerMiddlewareChain_InsertBeforeAll(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	m3 := &testServerMiddleware{name: "m3", list: &list}
	m4 := &testServerMiddleware{name: "m4", list: &list}
	c1 := middleware.ServerMiddlewareChain{m1}
	c1 = c1.InsertBeforeAll(m2)
	c1 = c1.InsertBeforeAll(m3, m4)
	c1.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1"}, list)
}

func TestServerMiddlewareChain_InsertBeforeAllFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := (&testServerMiddleware{name: "m2", list: &list}).ServerMiddleware
	m3 := (&testServerMiddleware{name: "m3", list: &list}).ServerMiddleware
	m4 := (&testServerMiddleware{name: "m4", list: &list}).ServerMiddleware
	c1 := middleware.ServerMiddlewareChain{m1}
	c1 = c1.InsertBeforeAllFunc(m2)
	c1 = c1.InsertBeforeAllFunc(m3, m4)
	c1.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m3", "m4", "m2", "m3", "m4", "m1"}, list)
}

func TestServerMiddlewareChain_InsertAfterAll(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := &testServerMiddleware{name: "m2", list: &list}
	m3 := &testServerMiddleware{name: "m3", list: &list}
	m4 := &testServerMiddleware{name: "m4", list: &list}
	c1 := middleware.ServerMiddlewareChain{m1}
	c1 = c1.InsertAfterAll(m2)
	c1 = c1.InsertAfterAll(m3, m4)
	c1.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m1", "m3", "m4", "m2", "m3", "m4"}, list)
}

func TestServerMiddlewareChain_InsertAfterAllFunc(t *testing.T) {
	t.Parallel()
	list := []string{}
	h := &testHandler{}
	m1 := &testServerMiddleware{name: "m1", list: &list}
	m2 := (&testServerMiddleware{name: "m2", list: &list}).ServerMiddleware
	m3 := (&testServerMiddleware{name: "m3", list: &list}).ServerMiddleware
	m4 := (&testServerMiddleware{name: "m4", list: &list}).ServerMiddleware
	c1 := middleware.ServerMiddlewareChain{m1}
	c1 = c1.InsertAfterAllFunc(m2)
	c1 = c1.InsertAfterAllFunc(m3, m4)
	c1.Terminate(h).ServeHTTP(nil, nil)
	tester.AssertEqual(t, true, h.called)
	tester.AssertDeepEqual(t, []string{"m1", "m3", "m4", "m2", "m3", "m4"}, list)
}
