/*
 *  Copyright icp-filing Author(https://houseme.github.io/icp-filing/). All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  You can obtain one at https://github.com/houseme/icp-filing.
 */

// Package filling .
package filling

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	authorizePath = "auth"
	queryPath     = "icpAbbreviateInfo/queryByCondition"

	authorizeContentType = "application/x-www-form-urlencoded;charset=UTF-8"
	queryContentType     = "application/json;charset=UTF-8"

	defaultToken = "0"
)

// Filling is the icp filling number object
type Filling struct {
	token string
	ip    string
}

type options struct {
	// LogPath is the path of log file.
	LogPath string
	// LogLevel is the level of log.
	LogLevel hlog.Level
}

// Option is the option for logger.
type Option func(o *options)

// WithLogPath is the option for log path.
func WithLogPath(path string) Option {
	return func(o *options) {
		o.LogPath = path
	}
}

// WithLogLevel is the option for log level.
func WithLogLevel(level hlog.Level) Option {
	return func(o *options) {
		o.LogLevel = level
	}
}

// New return a new filling number object
func New(_ context.Context, opts ...Option) *Filling {
	var op = options{
		LogPath:  os.TempDir(),
		LogLevel: hlog.LevelDebug,
	}
	for _, opt := range opts {
		opt(&op)
	}
	initLog(op.LogPath, op.LogLevel)
	return &Filling{
		token: defaultToken,
		ip:    "101." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)),
	}
}

// doRequest execute request
func (i *Filling) doRequest(ctx context.Context, in *ParamInput) ([]byte, error) {
	hc, err := client.NewClient(client.WithTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	}), client.WithDialTimeout(30*time.Second))
	if err != nil {
		return nil, err
	}

	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.Set("Content-Type", in.ContentType)
	req.Header.Set("Origin", "https://beian.miit.gov.cn/")
	req.Header.Set("Referer", "https://beian.miit.gov.cn/")
	req.Header.Set("token", i.token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36")
	req.Header.Set("CLIENT_IP", i.ip)
	req.Header.Set("X-FORWARDED-FOR", i.ip)
	req.Header.SetMethod(consts.MethodPost)
	req.SetRequestURI("https://hlwicpfwc.miit.gov.cn/icpproject_query/api/" + in.Path)

	hlog.CtxDebugf(ctx, "do request in url: %s, params: %s", req.RequestURI(), in.String())
	if in.Path != authorizePath {
		jsonByte, err := json.Marshal(in.QueryRequest)
		if err != nil {
			return nil, err
		}
		req.SetBody(jsonByte)
		hlog.CtxDebugf(ctx, "do request in json body: %s", string(jsonByte))
	} else {
		req.SetFormData(map[string]string{
			"authKey":   in.AuthorizeRequest.AuthKey,
			"timeStamp": in.AuthorizeRequest.Timestamp,
		})
		hlog.CtxDebugf(ctx, "do request in form data body: %s", string(req.PostArgString()))
	}

	if err = hc.Do(ctx, req, res); err != nil {
		return nil, err
	}
	hlog.CtxDebugf(ctx, "do request out status code: %d , body: %s", res.StatusCode(), string(res.Body()))
	if res.StatusCode() >= http.StatusMultipleChoices {
		return nil, errors.New(`请求接口 ` + in.Path + ` 失败! ,返回状态码: ` + strconv.Itoa(res.StatusCode()) + ` 返回内容: ` + string(res.Body()))
	}

	return res.Body(), nil
}

// authorize .
func (i *Filling) authorize(ctx context.Context) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	resp, err := i.doRequest(ctx,
		&ParamInput{
			AuthorizeRequest: &AuthorizeRequest{
				AuthKey:   i.md5("testtest" + timestamp),
				Timestamp: timestamp,
			},
			QueryRequest: nil,
			Path:         authorizePath,
			ContentType:  authorizeContentType,
		})
	if err != nil {
		return err
	}
	var response *AuthorizeResponse
	if err = sonic.Unmarshal(resp, &response); err != nil {
		return err
	}
	if response == nil {
		return fmt.Errorf("response is nil")
	}
	if !response.Success {
		return errors.New("code: " + strconv.Itoa(response.Code) + " errMsg: " + response.Msg)
	}
	i.token = response.Params.Business
	return nil
}

// QueryFilling query domain filling number
func (i *Filling) QueryFilling(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	resp, err := i.doRequest(ctx, &ParamInput{
		QueryRequest:     req,
		AuthorizeRequest: nil,
		ContentType:      queryContentType,
		Path:             queryPath,
	})
	if err != nil {
		return nil, err
	}
	var queryResp *QueryResponse
	if err = sonic.Unmarshal(resp, &queryResp); err != nil {
		return nil, err
	}
	return queryResp, nil
}

// md5 .
func (i *Filling) md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str))) // 将[]byte转成16进制
}

// String return filling json string
func (i *Filling) String() string {
	return `{"ip":"` + i.ip + `","token":"` + i.token + `"}`
}

// DomainFilling query domain filling number
func (i *Filling) DomainFilling(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	if err := i.authorize(ctx); err != nil {
		return nil, err
	}
	return i.QueryFilling(ctx, req)
}
