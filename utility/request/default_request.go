/*
 * Copyright icp-filing Author(https://houseme.github.io/icp-filing/). All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * You can obtain one at https://github.com/houseme/icp-filing.
 *
 */

package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	headerContentType      = "Content-Type"
	headerContentTypeValue = "application/json;charset=utf-8"
	headerUserAgent        = "User-Agent"
	headerUserAgentValue   = `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36`
)

// DefaultRequest 默认请求
type DefaultRequest struct {
}

// NewDefaultRequest 实例化
func NewDefaultRequest() *DefaultRequest {
	return &DefaultRequest{}
}

// Get HTTP get request
func (srv *DefaultRequest) Get(ctx context.Context, url string, headMap map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if headMap != nil {
		for key, value := range headMap {
			if strings.TrimSpace(value) != "" {
				req.Header.Set(key, value)
			}
		}
	}
	req.Header.Set(headerUserAgent, headerUserAgentValue)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// Post HTTP post request
func (srv *DefaultRequest) Post(ctx context.Context, url string, data []byte, headMap map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if headMap != nil {
		for key, value := range headMap {
			if strings.TrimSpace(value) != "" {
				req.Header.Set(key, value)
			}
		}
	}
	req.Header.Set(headerUserAgent, headerUserAgentValue)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// PostJSON HTTP post JSON request
func (srv *DefaultRequest) PostJSON(ctx context.Context, url string, data any, headMap map[string]string) ([]byte, error) {
	var (
		jsonBuf = new(bytes.Buffer)
		enc     = json.NewEncoder(jsonBuf)
	)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, jsonBuf)
	if err != nil {
		return nil, err
	}
	if headMap != nil {
		for key, value := range headMap {
			if strings.TrimSpace(value) != "" {
				req.Header.Set(key, value)
			}
		}
	}
	req.Header.Set(headerContentType, headerContentTypeValue)
	req.Header.Set(headerUserAgent, headerUserAgentValue)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// PostJSONWithRespContentType HTTP post JSON request with the response content type
func (srv *DefaultRequest) PostJSONWithRespContentType(ctx context.Context, url string, data any) ([]byte, string, error) {
	var (
		jsonBuf = new(bytes.Buffer)
		enc     = json.NewEncoder(jsonBuf)
	)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, jsonBuf)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set(headerContentType, headerContentTypeValue)
	req.Header.Set(headerUserAgent, headerUserAgentValue)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("http post error : uri=%v , statusCode=%v", url, resp.StatusCode)
	}
	res, err := io.ReadAll(resp.Body)
	contentType := resp.Header.Get(headerContentType)
	return res, contentType, err
}

// PostFile HTTP post file request
func (srv *DefaultRequest) PostFile(ctx context.Context, url string, files []MultipartFormField) ([]byte, error) {
	return srv.PostMultipartForm(ctx, url, files)
}

// PostMultipartForm HTTP post multipart form request
func (srv *DefaultRequest) PostMultipartForm(ctx context.Context, url string, files []MultipartFormField) (resp []byte, err error) {
	var (
		bodyBuf    = &bytes.Buffer{}
		bodyWriter = multipart.NewWriter(bodyBuf)
	)
	for _, field := range files {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.FieldName, field.FileName)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%w", e)
				return
			}

			fh, e := os.Open(field.FileName)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%w", e)
				return
			}

			if _, err = io.Copy(fileWriter, fh); err != nil {
				_ = fh.Close()
				return
			}
			_ = fh.Close()
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.FieldName)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%w", e)
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set(headerContentType, contentType)
	req.Header.Set(headerUserAgent, headerUserAgentValue)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v", url, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

// PostXML perform the HTTP/POST request with XML body
func (srv *DefaultRequest) PostXML(ctx context.Context, url string, data any) ([]byte, error) {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(xmlData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(headerContentType, "application/xml;charset=utf-8")
	req.Header.Set(headerUserAgent, headerUserAgentValue)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", url, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

// PostXMLWithTLS perform the HTTP/POST request with XML body and TLS
func (srv *DefaultRequest) PostXMLWithTLS(ctx context.Context, url string, data any, ca, key string) ([]byte, error) {
	return nil, nil
}
