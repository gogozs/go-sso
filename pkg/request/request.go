package request

import (
	"bytes"
	"fmt"
	"go-sso/pkg/json"
	"go-sso/pkg/log"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string, headers map[string]interface{}, params map[string]interface{}) ([]byte, int, io.ReadCloser, error) {
	httpRequest, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}

	// add headers
	httpRequest.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		httpRequest.Header.Add(k, v.(string))
	}

	// add url params
	q := httpRequest.URL.Query()
	for k, v := range params {
		q.Add(k, v.(string))
	}
	httpRequest.URL.RawQuery = q.Encode()

	// run http request
	res, err := client.Do(httpRequest)
	if err != nil {
		log.Error(err)
		return nil, 0, nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // io.ReadWriter can only read once
	return data, res.StatusCode, res.Body, nil
}

func Post(url, baseUrl, token string, params map[string]string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
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
		log.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	return res, nil
}
