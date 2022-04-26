package client

import "io"

type Req struct {
	method   string
	path     string
	rawQuery string
	body     io.Reader
	bodyStr  string
	headers  map[string][]string
}

type option func(req *Req)

func WithBody(bodyStr string) option {
	return func(req *Req) {
		req.bodyStr = bodyStr
	}
}

func NewReq(method, path string, options ...option) *Req {
	r := &Req{method: method, path: path, headers: make(map[string][]string, 0)}
	for _, opt := range options {
		opt(r)
	}
	return r
}

func (r *Req) Path(path string) *Req {
	r.path = path
	return r
}

func (r *Req) RawQuery(rawQuery string) *Req {
	r.rawQuery = rawQuery
	return r
}
func (r *Req) GetRawQuery() string {
	return r.rawQuery
}

func (r *Req) BodyStr(bodyStr string) *Req {
	r.bodyStr = bodyStr
	return r
}

func (r *Req) Body(body io.Reader) *Req {
	r.body = body
	return r
}

func (r *Req) Headers(headers map[string][]string) *Req {
	r.headers = headers
	return r
}

func (r *Req) SetHeader(key, value string) *Req {
	r.headers[key] = []string{value}
	return r
}

func (r *Req) AddHeader(key, value string) *Req {
	values, ok := r.headers[key]
	if ok {
		r.headers[key] = append(values, value)
	} else {
		r.headers[key] = []string{value}
	}
	return r
}

func (r *Req) DeepClone() *Req {
	copiedReq := *r
	copiedHeaders := make(map[string][]string, len(r.headers))
	for key, values := range r.headers {
		copiedValues := make([]string, len(values))
		copy(copiedValues, values)
		copiedHeaders[key] = copiedValues
	}
	copiedReq.headers = copiedHeaders
	return &copiedReq
}
