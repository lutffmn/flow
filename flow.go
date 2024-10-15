package flow

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type FlowHandler struct {
	Handler http.Handler
}

func Handle(handler func(http.ResponseWriter, *http.Request)) *FlowHandler {
	return &FlowHandler{
		Handler: http.HandlerFunc(handler),
	}
}

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

func useMiddleware(handler http.Handler, middleware Middleware) http.Handler {
	return middleware(handler)
}

func reverse(order []Middleware) []Middleware {
	reversed := make([]Middleware, len(order))
	for i := len(order) - 1; i >= 0; i-- {
		reversed[len(order)-i-1] = order[i]
	}

	return reversed
}
