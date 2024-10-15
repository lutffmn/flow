package flow

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestHandle(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}

	fh := Handle(handler)

	if fh == nil {
		t.Errorf("Handle returned nil")
	}

	if fh.Handler == nil {
		t.Errorf("Handler field is nil")
	}
}

func TestThru(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}

	fh := Handle(handler)

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

	result := fh.Thru(middleware1, middleware2)

	left := reflect.ValueOf(result)
	right := reflect.ValueOf(fh.Handler)

	if !reflect.DeepEqual(left.Pointer(), right.Pointer()) {
		t.Errorf("result %v\twant %v", result, fh.Handler)
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
