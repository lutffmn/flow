package flow

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitAndFlow(t *testing.T) {
	var middleware1Called, middleware2Called bool

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware1Called = true
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware2Called = true
			next.ServeHTTP(w, r)
		})
	}

	streams := Init(middleware1)
	streams.Extend(middleware2, middleware1)
	streams.Reduce(2)

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	streams.Flow(handlerFunc, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	streams[0](streams[1](http.HandlerFunc(handlerFunc))).ServeHTTP(w, req)

	if !middleware1Called {
		t.Error("Middleware 1 was not called")
	}

	if !middleware2Called {
		t.Error("Middleware 2 was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleAndThru(t *testing.T) {
	var middleware1Called, middleware2Called bool

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware1Called = true
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware2Called = true
			next.ServeHTTP(w, r)
		})
	}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	flowHandler := Handle(handlerFunc).Thru(middleware1, middleware2)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	flowHandler.ServeHTTP(w, req)

	if !middleware1Called {
		t.Error("Middleware 1 was not called")
	}

	if !middleware2Called {
		t.Error("Middleware 2 was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReverse(t *testing.T) {
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	right := []Middleware{middleware1, middleware2}
	left := reverse([]Middleware{middleware2, middleware1})

	for i, middleware := range left {
		if fmt.Sprintf("%v", middleware) != fmt.Sprintf("%v", right[i]) {
			t.Errorf("left: %v\tright: %v\n", left, right)
		}
	}
}

func TestSingleMiddlewareFlow(t *testing.T) {
	var middlewareCalled bool

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	streams := Init(middleware)

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	streams.Flow(handlerFunc, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	streams[0](http.HandlerFunc(handlerFunc)).ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSingleMiddlewareThru(t *testing.T) {
	var middlewareCalled bool

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	flowHandler := Handle(handlerFunc).Thru(middleware)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	flowHandler.ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
