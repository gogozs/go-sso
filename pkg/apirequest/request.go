package apirequest

import (
	"bytes"
	"fmt"
	"go-weixin/config"
	"go-weixin/pkg/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Get(url string, headers map[string]string, params map[string]string) ([]byte, io.ReadCloser) {
	httpRequest, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}

	// add headers
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	} else {
		httpRequest.Header.Add("Content-Type", "application/json")
		httpRequest.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) 
AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`)
	}

	// add url params
	q := httpRequest.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	httpRequest.URL.RawQuery = q.Encode()

	// run http request
	res, err := client.Do(httpRequest)
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		res.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // io.ReadWriter can only read once
		return data, res.Body
	}
}

func Post(url string, params map[string]string, body interface{}) *http.Response {
	c := config.GetConfig()
	baseUrl := c.Ops.Url
	token := c.Ops.Token
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	httpRequest, _ := http.NewRequest("POST", baseUrl+url, strings.NewReader(string(b)))
	client := &http.Client{}

	// add headers
	httpRequest.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	httpRequest.Header.Add("Content-Type", "application/json")

	// add url params
	q := httpRequest.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	httpRequest.URL.RawQuery = q.Encode()

	// run http request
	res, err := client.Do(httpRequest)
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		return res
	}
}
