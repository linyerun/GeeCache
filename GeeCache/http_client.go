package GeeCache

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/linyerun/GeeCache/protobuf"
	"github.com/linyerun/GeeCache/utils"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient struct {
	baseUrl string
}

func NewHttpClient(baseUrl string) INodeClient {
	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl = baseUrl + "/"
	}
	return &httpClient{baseUrl: baseUrl}
}

func (h *httpClient) Get(group string, key string) ([]byte, error) {
	serverURL := jointServerURL(h.baseUrl, group, key)
	res, err := http.Get(serverURL)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	defer func() { _ = res.Body.Close() }()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	resp := new(protobuf.Response)
	err = proto.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp.GetValue(), nil
}

func (h *httpClient) Post(group string, key string, value []byte) error {
	serverURL := jointServerURL(h.baseUrl, group, key)
	resp, err := http.Post(serverURL, binaryContentType, bytes.NewBuffer(value))

	// 错误处理
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.Logger().Errorln(err.Error())
			return err
		}
		return errors.New(string(data))
	}

	return nil
}

func (h *httpClient) Delete(group string, key string) error {
	serverURL := jointServerURL(h.baseUrl, group, key)

	// 创建HttpRequest
	request, err := http.NewRequest(http.MethodDelete, serverURL, nil)
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.Logger().Errorln(err.Error())
			return err
		}
		return errors.New(string(data))
	}

	return nil
}

func (h *httpClient) GetBaseUrl() string {
	return h.baseUrl
}

func jointServerURL(baseUrl, group, key string) string {
	return fmt.Sprintf(
		"http://%v%v%v/%v",
		baseUrl,
		globalBasePath,
		url.QueryEscape(group), //QueryEscape函数对s进行转码使之可以安全的用在URL查询里。有些特殊字符不能用可转的
		url.QueryEscape(key),
	)
}
