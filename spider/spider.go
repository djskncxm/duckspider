package spider

import (
	"github.com/djskncxm/duckspider/httpio"
)

type BaseSpider interface {
	Name() string
	StartRequest() <-chan *httpio.Request
}

type MakeSpider struct {
	SpiderName string
	URLS       []string
	URL        string
	Callback   func(response string) <-chan interface{}
}

func (m MakeSpider) Name() string {
	return m.SpiderName
}

func (m MakeSpider) StartRequest() <-chan *httpio.Request {
	ch := make(chan *httpio.Request, 5)
	go func() {
		defer close(ch)
		if m.URL != "" {
			ch <- httpio.NewRequest(m.URL, httpio.WithCallback(m.Callback))
		} else {
			for _, u := range m.URLS {
				ch <- httpio.NewRequest(u, httpio.WithCallback(m.Callback))
			}
		}
	}()
	return ch
}
