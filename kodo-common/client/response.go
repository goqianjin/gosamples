package client

type Resp struct {
	Body       []byte
	StatusCode int
	Err        error
	Headers    map[string][]string
}
