package core

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/djskncxm/duckspider/httpio"
)

type Downloader struct {
	activeQueue *treeset.Set
	mu          sync.Mutex
}

func NewDownloader() *Downloader {
	return &Downloader{
		activeQueue: treeset.NewWithStringComparator(),
	}
}

func (d *Downloader) Fetch(request *httpio.Request) string {
	d.mu.Lock()
	d.activeQueue.Add(request.UUID)
	d.mu.Unlock()

	defer func() {
		d.mu.Lock()
		d.activeQueue.Remove(request.UUID)
		d.mu.Unlock()
	}()

	resp := d.DownloadTest(request)
	return resp
}

func (d *Downloader) Download(request *httpio.Request) string {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(request.URL)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)

	return string(data)
}

func (d *Downloader) DownloadTest(request *httpio.Request) string {
	time.Sleep(1 * time.Second)
	fmt.Println("Download Test Is => " + request.URL)
	return "response"
}
