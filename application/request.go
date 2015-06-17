package application

import (
	"net/http"
	"net/url"

	"github.com/jmoiron/jsonq"
	"github.com/lizdeika/picfit/engines"
	"github.com/lizdeika/picfit/hash"
	"github.com/lizdeika/picfit/signature"
	"github.com/lizdeika/picfit/util"
	"github.com/thoas/gokvstores"
	"github.com/thoas/muxer"
)

type Request struct {
	*muxer.Request
	Operation  *engines.Operation
	Connection gokvstores.KVStoreConnection
	Key        string
	URL        *url.URL
	Filepath   string
}

const SIG_PARAM_NAME = "sig"

func NewRequest(req *http.Request, con gokvstores.KVStoreConnection, jq *jsonq.JsonQuery) (*Request, error) {
	request := muxer.NewRequest(req)

	for k, v := range request.Params {
		request.QueryString[k] = v
	}

	if request.QueryString["op"] == "" {
		op, err := jq.String("defaults", "operation")
		if err != nil {
			request.QueryString["op"] = "fit"
		} else {
			request.QueryString["op"] = op
		}
	}

	if request.QueryString["w"] == "" {
		w, err := jq.String("defaults", "width")
		if err != nil {
			request.QueryString["w"] = "2000"
		} else {
			request.QueryString["w"] = w
		}
	}

	if request.QueryString["h"] == "" {
		h, err := jq.String("defaults", "height")
		if err != nil {
			request.QueryString["h"] = "2000"
		} else {
			request.QueryString["h"] = h
		}
	}

	extracted := map[string]interface{}{}

	for key, extractor := range Extractors {
		result, err := extractor(key, request)

		if err != nil {
			return nil, err
		}

		extracted[key] = result
	}

	sorted := util.SortMapString(request.QueryString)

	delete(sorted, SIG_PARAM_NAME)

	serialized := hash.Serialize(sorted)

	key := hash.Tokey(serialized)

	var u *url.URL
	var path string

	value, ok := extracted["url"]

	if ok && value != nil {
		u = value.(*url.URL)
	}

	value, ok = extracted["path"]

	if ok && value != nil {
		path = value.(string)
	}

	return &Request{
		request,
		extracted["op"].(*engines.Operation),
		con,
		key,
		u,
		path,
	}, nil
}

func (r *Request) IsAuthorized(key string) bool {
	params := url.Values{}
	for k, v := range util.SortMapString(r.QueryString) {
		params.Set(k, v)
	}

	return signature.VerifySign(key, params.Encode())
}
