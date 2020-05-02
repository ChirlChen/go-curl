/**
 * golang版本的curl请求库
 * Request构造类，用于设置请求参数，发起http请求
 * @author mike <mikemintang@126.com>
 * @blog http://idoubi.cc
 */

package curl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"runtime/debug"
)

// Request构造类
type Request struct {
	cli           *http.Client
	req           *http.Request
	Raw           *http.Request
	Method        string
	Url           string
	Headers       map[string]string
	Cookies       map[string]string
	Queries       map[string]string
	PostData      interface{}
	PostDataQuery map[string]interface{}
}

// 创建一个Request实例
func NewRequest() *Request {
	r := &Request{}
	r.req = &http.Request{}
	return r
}

// 设置请求方法
func (this *Request) SetMethod(method string) *Request {
	this.Method = method
	return this
}

// 设置请求地址
func (this *Request) SetUrl(url string) *Request {
	this.Url = url
	return this
}

// 设置请求头
func (this *Request) SetHeaders(headers map[string]string) *Request {
	this.Headers = headers
	return this
}

// 将用户自定义请求头添加到http.Request实例上
func (this *Request) setHeaders() error {
	for k, v := range this.Headers {
		this.req.Header.Set(k, v)
	}
	return nil
}

// 设置请求cookies
func (this *Request) SetCookies(cookies map[string]string) *Request {
	this.Cookies = cookies
	return this
}

// 将用户自定义cookies添加到http.Request实例上
func (this *Request) setCookies() error {
	for k, v := range this.Cookies {
		this.req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return nil
}

// 设置url查询参数
func (this *Request) SetQueries(queries map[string]string) *Request {
	this.Queries = queries
	return this
}

// 将用户自定义url查询参数添加到http.Request上
func (this *Request) setQueries() error {
	q := this.req.URL.Query()
	for k, v := range this.Queries {
		q.Add(k, v)
	}
	this.req.URL.RawQuery = q.Encode()
	return nil
}

// 设置post请求的提交数据
func (this *Request) SetPostData(postData interface{}) *Request {
	this.PostData = postData
	return this
}

// 设置post请求的提交数据
func (this *Request) SetPostDataUrlEncode(postData map[string]interface{}) *Request {
	this.PostDataQuery = postData
	return this
}

// 发起get请求
func (this *Request) Get() (*Response, error) {
	return this.Send(this.Url, http.MethodGet)
}

// 发起Delete请求
func (this *Request) Delete() (*Response, error) {
	return this.Send(this.Url, http.MethodDelete)
}

// 发起Delete请求
func (this *Request) Put() (*Response, error) {
	return this.Send(this.Url, http.MethodPut)
}

// 发起post请求
func (this *Request) Post() (*Response, error) {
	return this.Send(this.Url, http.MethodPost)
}

// 发起put请求
func (this *Request) PUT() (*Response, error) {
	return this.Send(this.Url, http.MethodPut)
}

// 发起put请求
func (this *Request) PATCH() (*Response, error) {
	return this.Send(this.Url, http.MethodPatch)
}

// 发起请求
func (this *Request) Send(url string, method string) (*Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover found：%s\n", r)
			debug.PrintStack()
		}
	}()
	// 检测请求url是否填了
	if url == "" {
		return nil, errors.New("Lack of request url")
	}
	// 检测请求方式是否填了
	if method == "" {
		return nil, errors.New("Lack of request method")
	}
	// 初始化Response对象
	response := NewResponse()
	// 初始化http.Client对象
	this.cli = &http.Client{}
	// 加载用户自定义的post数据到http.Request
	var payload io.Reader
	if (method == "POST" || method == "PUT") && this.PostData != nil {
		jData, err := json.Marshal(this.PostData)
		if err != nil {
			return nil, err
		} 
		fmt.Printf("xxxxxxxx: %v", this.PostData)
		payload = bytes.NewReader(jData)
		
	} else {
		fmt.Printf("xxxxxxxx: %v", this.PostData)
		payload = nil
	}
	if (method == "POST" || method == "PUT") && this.PostDataQuery != nil {
		this.req.PostForm = url2.Values{}
		for k, v := range this.PostDataQuery {
			this.req.PostForm.Add(k, fmt.Sprint(v))
		}
	} else {
		payload = nil
	}

	if req, err := http.NewRequest(method, url, payload); err != nil {
		return nil, err
	} else {
		this.req = req
	}

	this.setHeaders()
	this.setCookies()
	this.setQueries()

	this.Raw = this.req

	if resp, err := this.cli.Do(this.req); err != nil {
		return nil, err
	} else {
		response.Raw = resp
	}

	defer response.Raw.Body.Close()

	response.parseHeaders()
	response.parseBody()

	return response, nil
}
