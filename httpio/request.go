package httpio

import (
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

type Request struct {
	URL      string
	Header   http.Header
	Body     []byte
	Method   string
	Cookies  []*http.Cookie
	Proxy    map[string]*url.URL
	Meat     map[string]interface{}
	Priority int
	Callback func(str string) <-chan interface{}
	UUID     string
}

func NewRequest(URL string, opts ...func(*Request)) *Request {
	r := &Request{
		URL:    URL,
		Header: http.Header{},
		UUID:   uuid.New().String(),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithCallback(callback func(str string) <-chan interface{}) func(*Request) {
	return func(r *Request) {
		r.Callback = callback
	}
}

func WithMeat(meat map[string]interface{}) func(*Request) {
	return func(r *Request) {
		r.Meat = meat
	}
}
