package core

import (
	"fmt"
	"github.com/djskncxm/duckspider/httpio"
	"github.com/djskncxm/duckspider/spider"
	"sync"
)

type Engine struct {
	Downloader *Downloader
	Scheduler  *Scheduler
	Spider     spider.BaseSpider
}

func NewEngine() *Engine {
	return &Engine{
		Downloader: NewDownloader(),
		Scheduler:  NewScheduler(),
	}
}

func (engine *Engine) Start(spider spider.BaseSpider) {
	engine.Spider = spider
	engine.OpenSpider(spider)
}

func (engine *Engine) OpenSpider(spider spider.BaseSpider) {
	var wg sync.WaitGroup

	wg.Add(1)
	go engine.Crawl(spider, &wg)
	wg.Wait()
}

func (engine *Engine) Crawl(spider spider.BaseSpider, wg *sync.WaitGroup) {
	defer wg.Done()
	request := spider.StartRequest()
	for {
		if ok := engine.Scheduler.IsEmpty(); !ok {
			engine._Crawl()
			continue
		}
		QueueRequest, more := <-request
		if !more {
			break
		}
		engine.EnRequest(QueueRequest)
	}
}

func (engine *Engine) _Crawl() {
	var wg sync.WaitGroup
	requestCh := make(chan *httpio.Request, 8)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for req := range requestCh {
				engine.worker(req)
			}
		}()
	}

	go func() {
		for {
			req, ok := engine.GetNextRequest()
			if !ok {
				break
			}
			requestCh <- req
		}
		close(requestCh)
	}()

	wg.Wait()
}

func (engine *Engine) worker(request *httpio.Request) {
	outputs := engine.Fetch(request)
	if outputs != nil {
		for output := range outputs {
			engine.DiversionRI(output)
		}
	}
}

func (engine *Engine) Fetch(request *httpio.Request) <-chan interface{} {
	response := engine.Downloader.Fetch(request)
	if request.Callback != nil {
		return request.Callback(response)
	} else {
		fmt.Println("此次request无回调函数")
	}
	return nil
}

func (engine *Engine) DiversionRI(data interface{}) {
	switch v := data.(type) {
	case *httpio.Request:
		engine.EnRequest(v)
	}
}

func (engine *Engine) EnRequest(request *httpio.Request) {
	engine.ScheduleRequest(request)
}

func (engine *Engine) ScheduleRequest(request *httpio.Request) {
	// TODO 日后去重
	engine.Scheduler.EnRequest(request)
}

func (engine *Engine) GetNextRequest() (*httpio.Request, bool) {
	return engine.Scheduler.NextRequest()
}
