# 框架文档 V 0.0.1
### 关于创建爬虫
``` go
type MakeSpider struct {
	SpiderName string
	URLS       []string
	URL        string
}
```
你如果只有一个起始的URL你就只写URL就可以了，如果你有多个起始的URL你就选择填写URLS，我在这个内部进行了处理
