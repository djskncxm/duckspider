package spider

import (
	"github.com/djskncxm/duckspider/httpio"
	"github.com/djskncxm/duckspider/utlis"
	"gopkg.in/yaml.v3"
	"os"
)

type BaseSpider interface {
	Name() string
	StartRequest() <-chan *httpio.Request
	LoadConfig() *utlis.SettingManager
}

type MakeSpider struct {
	SpiderName   string
	URLS         []string
	URL          string
	Callback     func(response string) <-chan interface{}
	SettingsPath string
}

func (m MakeSpider) Name() string {
	return m.SpiderName
}

func (m MakeSpider) LoadConfig() *utlis.SettingManager {
	defaultPath, err := os.Getwd()

	if m.SettingsPath == "" {
		m.SettingsPath = defaultPath + "/config/config.yaml"
	}

	data, err := os.ReadFile(m.SettingsPath)
	if err != nil {
		panic(err)
	}
	var cfg utlis.Setting
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	config := utlis.NewSettingManager()
	config.LoadFromSetting(cfg)
	return config
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
