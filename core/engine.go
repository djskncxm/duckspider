package core

import (
	"fmt"
	"github.com/djskncxm/duckspider/httpio"
	"github.com/djskncxm/duckspider/items"
	"github.com/djskncxm/duckspider/spider"
	"github.com/djskncxm/duckspider/utlis"
	"sync"
)

type Engine struct {
	Downloader *Downloader
	Scheduler  *Scheduler
	Crawler    *Crawler

	Settings *utlis.SettingManager
	Spider   spider.BaseSpider
	Logger   *utlis.Logger
}

func NewEngine(crawler *Crawler) *Engine {
	return &Engine{
		Downloader: NewDownloader(),
		Scheduler:  NewScheduler(),
		Crawler:    crawler,
	}
}

func (engine *Engine) Start(spider spider.BaseSpider, Settings *utlis.SettingManager) {
	engine.Spider = spider
	engine.Settings = Settings
	engine.Logger = utlis.InitLogger(spider.Name(), engine.Crawler.LogLevel)
	engine.Logger.Infof("Starting spider %s", spider.Name())
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
		if engine.Scheduler.IsEmpty() && engine.Downloader.activeQueue.Empty() {
			break
		}
	}
}

func (engine *Engine) _Crawl() {
	var wg sync.WaitGroup
	requestCh := make(chan *httpio.Request, engine.Settings.GetInt("Spider.WorkerNumber", 8))

	for i := 0; i < engine.Settings.GetInt("Spider.WorkerNumber", 8); i++ {
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
	}
	return nil
}

func (engine *Engine) DiversionRI(data interface{}) {
	switch v := data.(type) {
	case *httpio.Request:
		engine.EnRequest(v)
	case *items.StrictItem:
		fmt.Println(v.All())
	default:
		fmt.Println("")
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
