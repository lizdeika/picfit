package muxer

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Request struct {
	Request     *http.Request
	Params      map[string]string
	QueryString map[string]string
}

// CurrentRoute allows to retrieve the current
// from a request
func (r *Request) CurrentRoute() *mux.Route {
	return mux.CurrentRoute(r.Request)
}

func (r *Request) GetHost() string {
	return getHost(r.Request)
}

func (r *Request) IsAjax() bool {
	return matchMap(map[string]string{
		"HTTP_X_REQUESTED_WITH": "XMLHttpRequest",
	}, r.Request.Header, true)
}

func NewRequest(req *http.Request) *Request {
	req.ParseForm()

	queryString := make(map[string]string)

	if len(req.Form) > 0 {
		for k, v := range req.Form {
			queryString[k] = v[0]
		}
	}

	return &Request{
		Request:     req,
		QueryString: queryString,
		Params:      mux.Vars(req),
	}
}
