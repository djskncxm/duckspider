package core

import (
	"fmt"
	"github.com/djskncxm/duckspider/spider"
	"github.com/djskncxm/duckspider/utlis"
	"github.com/emirpasic/gods/sets/treeset"
	"sync"
)

type Crawler struct {
	Setting  *utlis.SettingManager
	LogLevel string
	Spider   spider.BaseSpider
	Engine   *Engine
}

func NewCrawler(spider spider.BaseSpider) *Crawler {
	Setting := spider.LoadConfig()
	LogLevel, _ := Setting.GetSetting("Spider.LOGLEVEL")
	return &Crawler{
		Spider:   spider,
		Setting:  Setting,
		LogLevel: LogLevel,
	}
}
func (crawler *Crawler) Crawl() {
	crawler.Engine = NewEngine(crawler)
	crawler.Engine.Start(crawler.Spider, crawler.Setting)
	//crawler.Engine.Logger.Stats.AddString("Over", "正常结束")
	//crawler.Engine.Logger.Stats.OutTableInfo()
}

type CrawlerProcess struct {
	crawlerSet map[string]*Crawler
	nameSet    *treeset.Set
}

func NewProcess() *CrawlerProcess {
	return &CrawlerProcess{
		crawlerSet: make(map[string]*Crawler),
		nameSet:    treeset.NewWithStringComparator(),
	}
}

func (process *CrawlerProcess) AddSpider(spider spider.BaseSpider) {
	crawler := process.CreateCrawl(spider)
	Name, ok := crawler.Setting.GetSetting("Spider.Name")
	if !ok {
		panic("未获取到Spider.Name")
	}
	if process.nameSet.Contains(Name) {
		fmt.Printf("[警告] 爬虫 %s 已存在，跳过注册\n", Name)
		return
	}
	process.nameSet.Add(Name)
	process.crawlerSet[Name] = crawler
}

func (process *CrawlerProcess) CreateCrawl(spider spider.BaseSpider) *Crawler {
	return NewCrawler(spider)
}

func (process *CrawlerProcess) StartCrawlers() {
	process.Crawl()
}

func (process *CrawlerProcess) Crawl() {
	var wg sync.WaitGroup
	SpiderName := process.nameSet.Iterator()
	for SpiderName.Next() {
		Name, ok := SpiderName.Value().(string)
		if !ok {
			continue
		}
		crawler := process.crawlerSet[Name]
		if crawler == nil {
			continue
		}
		wg.Add(1)
		go func(c *Crawler, name string) {
			defer wg.Done()
			c.Crawl()
		}(crawler, Name)
	}
	wg.Wait()
}
