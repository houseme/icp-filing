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

// Package filling is the icp filling number
package filling

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/houseme/icp-filing/tld"
	"github.com/houseme/icp-filing/utility/logger"
	"github.com/houseme/icp-filing/utility/request"
)

const (
	authorizePath = "auth"
	queryPath     = "icpAbbreviateInfo/queryByCondition"

	authorizeContentType = "application/x-www-form-urlencoded;charset=UTF-8"
	queryContentType     = "application/json;charset=UTF-8"

	originAndReferer = "https://beian.miit.gov.cn/"

	defaultToken = "0"

	domainLevel = 0

	randomIP = "101.%d.%d.%d"

	maxRetry = 255
)

// Filling is the icp filling number object
type Filling struct {
	token   string
	ip      string
	request request.Request
	logger  logger.ILogger
}

type options struct {
	Request request.Request
	Logger  logger.ILogger
}

// Option is the option for logger.
type Option func(o *options)

// WithRequest is the option for request.
func WithRequest(req request.Request) Option {
	return func(o *options) {
		o.Request = req
	}
}

// WithLogger is the option for logger.
func WithLogger(logger logger.ILogger) Option {
	return func(o *options) {
		o.Logger = logger
	}
}

// New return a new filling number object
func New(ctx context.Context, opts ...Option) *Filling {
	var op = options{}
	for _, opt := range opts {
		opt(&op)
	}
	f := &Filling{
		token:   defaultToken,
		ip:      fmt.Sprintf(randomIP, rand.Intn(maxRetry), rand.Intn(maxRetry), rand.Intn(maxRetry)),
		logger:  op.Logger,
		request: op.Request,
	}
	return f
}

// doRequest execute request
func (i *Filling) doRequest(ctx context.Context, in *ParamInput) ([]byte, error) {
	headMap := map[string]string{
		"Content-Type":    in.ContentType,
		"Origin":          originAndReferer,
		"Referer":         originAndReferer,
		"token":           i.token,
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		"CLIENT_IP":       i.ip,
		"X-FORWARDED-FOR": i.ip,
	}
	url := "https://hlwicpfwc.miit.gov.cn/icpproject_query/api/" + in.Path
	i.logger.Debugf(ctx, "do request in url: %s, params: %s", url, in.String())
	var (
		resp []byte
		err  error
	)
	if in.Path != authorizePath {
		resp, err = i.request.PostJSON(ctx, url, in.QueryRequest, headMap)
	} else {
		var jsonByte []byte
		if jsonByte, err = json.Marshal(map[string]string{
			"authKey":   in.AuthorizeRequest.AuthKey,
			"timeStamp": in.AuthorizeRequest.Timestamp,
		}); err != nil {
			return nil, err
		}
		resp, err = i.request.Post(ctx, url, jsonByte, headMap)
	}

	return resp, err
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
	if err = json.Unmarshal(resp, &response); err != nil {
		return err
	}
	if response == nil {
		return errors.New("response is nil")
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
	if err = json.Unmarshal(resp, &queryResp); err != nil {
		return nil, err
	}
	return queryResp, nil
}

// md5 .
func (i *Filling) md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str))) // 将 []byte 转成 16 进制
}

// String return filling json string
func (i *Filling) String() string {
	return `{"ip":"` + i.ip + `","token":"` + i.token + `"}`
}

// DomainFilling query domain filling number
func (i *Filling) DomainFilling(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	if req.UnitName == "" && req.Link != "" {
		resp, err := tld.GetTLD(ctx, req.Link, domainLevel)
		if err != nil {
			return nil, err
		}
		i.logger.Debugf(ctx, "GetTld resp: %s", resp.String())
		req.UnitName = resp.Domain
	}

	if err := i.authorize(ctx); err != nil {
		return nil, err
	}

	return i.QueryFilling(ctx, req)
}
