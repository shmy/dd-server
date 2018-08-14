package util

import (
	"net/http"
	"io"
	"net/url"
)

var client *http.Client

func init()  {
	// 代理
	proxy := func(*http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1087")
	}
	transport := &http.Transport{Proxy: proxy}
	client = &http.Client{
		Transport: transport,
	}
}
// get 方法 ua
func HttpGet (url string) io.ReadCloser {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1")
	request.Header.Add("Host", "m.iqiyi.com")
	request.Header.Add("Upgrade-Insecure-Requests", "1")
	if err != nil {
		return nil
	}
	response, err := client.Do(request)
	//defer response.Body.Close()
	if err != nil {
		return nil
	}
	return response.Body
}