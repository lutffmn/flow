package flow

import (
	"net/http"
)

/*TODO:
- Create Init function for Flow Instance
- Create Flow method for Flow Instance which accept a specified handler as the arg
- and then start flowing it thru all the middlewares registered in Init function
- Make the test
*/

type Middleware func(http.Handler) http.Handler

type FlowHandler struct {
	Handler http.Handler
}

type Streams []Middleware

func Init(middlewares ...Middleware) Streams {
	streams := make(Streams, 0)
	if len(middlewares) > 0 {
		for _, m := range middlewares {
			streams = append(streams, m)
		}
	}

	return streams
}

func (s Streams) Flow(handler func(http.ResponseWriter, *http.Request)) {
	fh := FlowHandler{
		Handler: http.HandlerFunc(handler),
	}
	if len(s) > 1 {
		for _, m := range reverse(s) {
			fh.Handler = useMiddleware(fh.Handler, m)
		}
	} else {
		fh.Handler = useMiddleware(fh.Handler, s[0])
	}
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
