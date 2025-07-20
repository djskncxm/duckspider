package core

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/djskncxm/duckspider/httpio"
)

type Downloader struct {
	URL string
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Fetch(request *httpio.Request) string {
	return d.DownloadTest(request)
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
