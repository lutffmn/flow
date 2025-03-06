package flow

import (
	"fmt"
	"net/http"
	"sort"
)

// This is a type for Middleware
type Middleware func(http.Handler) http.Handler

// This is a type for Flow Instance
type FlowHandler struct {
	Handler http.Handler
}

// This is a type for multiple middewares
type Streams []Middleware

// Initialize new Flow Instance
func New(middlewares ...Middleware) Streams {
	streams := make(Streams, 0)
	if len(middlewares) > 0 {
		for _, m := range middlewares {
			streams = append(streams, m)
		}
	}

	return streams
}

// Execute Flow Instance
func (s Streams) Flow(handler func(http.ResponseWriter, *http.Request), exclude ...int) http.Handler {
	fh := FlowHandler{
		Handler: http.HandlerFunc(handler),
	}

	if len(exclude) > 0 && len(s) > 1 {
		for i, middleware := range s {
			exc := false
			for _, n := range exclude {
				if i == n {
					exc = true
					break
				}
			}
			if !exc {
				fh.Handler = useMiddleware(fh.Handler, middleware)
			}
		}
	} else if len(exclude) > 0 && len(s) == 1 {
		fh.Handler = useMiddleware(fh.Handler, s[0])
	} else if len(exclude) <= 0 && len(s) > 1 {
		for _, m := range reverse(s) {
			fh.Handler = useMiddleware(fh.Handler, m)
		}
	} else if len(exclude) <= 0 && len(s) == 1 {
		fh.Handler = useMiddleware(fh.Handler, s[0])
	}

	return fh.Handler
}

// Extend new middleware(s) to existing Flow Instance
func (s *Streams) Extend(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		*s = append(*s, middleware)
	}
}

// Reduce middleware(s) from existing Flow Instance
func (s *Streams) Reduce(index ...int) {
	if len(index) == 0 || len(*s) == 0 {
		return
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i] > index[j]
	})

	for _, index := range index {
		if index >= 0 && index < len(*s) {
			*s = append((*s)[:index], (*s)[index+1:]...)
		}
	}
}

// Print middlewares length of existing Flow Instance
func (s *Streams) Show() {
	fmt.Println(len(*s))
}

// Function to specifiying handler to pass into middleware(s)
func Handle(handler func(http.ResponseWriter, *http.Request)) *FlowHandler {
	return &FlowHandler{
		Handler: http.HandlerFunc(handler),
	}
}

// Lists all middleware(s) to be applied
func (fh *FlowHandler) Thru(middlewares ...Middleware) http.Handler {
	if len(middlewares) > 1 {
		for _, middleware := range reverse(middlewares) {
			fh.Handler = useMiddleware(fh.Handler, middleware)
		}
	} else {
		fh.Handler = useMiddleware(fh.Handler, middlewares[0])
	}
	return fh.Handler
}

// Activating middleware
func useMiddleware(handler http.Handler, middleware Middleware) http.Handler {
	return middleware(handler)
}

// Reversing order of middlewares
func reverse(order []Middleware) []Middleware {
	reversed := make([]Middleware, len(order))
	for i := len(order) - 1; i >= 0; i-- {
		reversed[len(order)-i-1] = order[i]
	}

	return reversed
}
