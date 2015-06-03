package muxer

import (
	"net/http"
)

type Handler func(Response, *Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h(NewResponse(w), NewRequest(req))
}
