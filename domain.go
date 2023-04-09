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

package filling

import (
	"strconv"

	"github.com/bytedance/sonic"
)

// QueryRequest query request
type QueryRequest struct {
	PageNum  string `json:"pageNum"`
	PageSize string `json:"pageSize"`
	UnitName string `json:"unitName" description:"unit name"`
	Link     string `json:"link" description:"link"`
}

// String return query request string
func (r *QueryRequest) String() string {
	return `{"pageNum": "` + r.PageNum + `", "pageSize": "` + r.PageSize + `", "unitName": "` + r.UnitName + `"}`
}

// QueryResponse query response
type QueryResponse struct {
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Success bool         `json:"success"`
	Params  *QueryParams `json:"params"`
}

// String return query response string
func (r *QueryResponse) String() string {
	return `{"code": ` + strconv.Itoa(r.Code) + `, "msg": "` + r.Msg + `", "success": ` + strconv.FormatBool(r.Success) + `, "params": ` + r.Params.String() + `}`
}

// AuthParams auth params
type AuthParams struct {
	Business string `json:"bussiness"`
	Expire   int64  `json:"expire"`
	Refresh  string `json:"refresh"`
}

// QueryParams query params
type QueryParams struct {
	EndRow           int           `json:"endRow"`
	FirstPage        int           `json:"firstPage"`
	HasNextPage      bool          `json:"hasNextPage"`
	HasPreviousPage  bool          `json:"hasPreviousPage"`
	IsFirstPage      bool          `json:"isFirstPage"`
	IsLastPage       bool          `json:"isLastPage"`
	LastPage         int           `json:"lastPage"`
	List             []*DomainInfo `json:"list"`
	NavigatePages    int           `json:"navigatePages"`
	NavigatePageNums []int         `json:"navigatepageNums"`
	NextPage         int           `json:"nextPage"`
	PageNum          int           `json:"pageNum"`
	PageSize         int           `json:"pageSize"`
	Pages            int           `json:"pages"`
	PrePage          int           `json:"prePage"`
	Size             int           `json:"size"`
	StartRow         int           `json:"startRow"`
	Total            int           `json:"total"`
}

// String return query params string
func (r *QueryParams) String() string {
	return `{"endRow": ` + strconv.Itoa(r.EndRow) + `, "firstPage": ` + strconv.Itoa(r.FirstPage) + `, "hasNextPage": ` + strconv.FormatBool(r.HasNextPage) + `, "hasPreviousPage": ` + strconv.FormatBool(r.HasPreviousPage) + `, "isFirstPage": ` + strconv.FormatBool(r.IsFirstPage) + `, "isLastPage": ` + strconv.FormatBool(r.IsLastPage) + `, "lastPage": ` + strconv.Itoa(r.LastPage) + `, "list": ` + r.ParamsListString() + `, "navigatePages": ` + strconv.Itoa(r.NavigatePages) + `, "navigatepageNums": ` + r.NavigatePageNumsString() + `, "nextPage": ` + strconv.Itoa(r.NextPage) + `, "pageNum": ` + strconv.Itoa(r.PageNum) + `, "pageSize": ` + strconv.Itoa(r.PageSize) + `, "pages": ` + strconv.Itoa(r.Pages) + `, "prePage": ` + strconv.Itoa(r.PrePage) + `, "size": ` + strconv.Itoa(r.Size) + `, "startRow": ` + strconv.Itoa(r.StartRow) + `, "total": ` + strconv.Itoa(r.Total) + `}`
}

// NavigatePageNumsString Navigate Page Nums String
func (r *QueryParams) NavigatePageNumsString() string {
	if r.NavigatePageNums == nil || len(r.NavigatePageNums) < 1 {
		return ""
	}
	output, err := sonic.MarshalString(r.NavigatePageNums)
	if err != nil {
		return ""
	}
	return output
}

// ParamsListString return params list to string
func (r *QueryParams) ParamsListString() string {
	if r.List == nil || len(r.List) < 1 {
		return ""
	}
	output, err := sonic.MarshalString(r.List)
	if err != nil {
		return ""
	}
	return output
}

// DomainInfo domain info
type DomainInfo struct {
	ContentTypeName  string `json:"contentTypeName"`
	Domain           string `json:"domain"`
	DomainID         int64  `json:"domainId"`
	HomeURL          string `json:"homeUrl"`
	LeaderName       string `json:"leaderName"`
	LimitAccess      string `json:"limitAccess"`
	MainID           int64  `json:"mainId"`
	MainLicence      string `json:"mainLicence"`
	NatureName       string `json:"natureName"`
	ServiceID        int64  `json:"serviceId"`
	ServiceLicence   string `json:"serviceLicence"`
	ServiceName      string `json:"serviceName"`
	UnitName         string `json:"unitName"`
	UpdateRecordTime string `json:"updateRecordTime"`
}

// String return domain info string
func (r *DomainInfo) String() string {
	return `{"contentTypeName": "` + r.ContentTypeName + `", "domain": "` + r.Domain + `", "domainId": ` + strconv.FormatInt(r.DomainID, 10) + `, "homeUrl": "` + r.HomeURL + `", "leaderName": "` + r.LeaderName + `", "limitAccess": "` + r.LimitAccess + `", "mainId": ` + strconv.FormatInt(r.MainID, 10) + `, "mainLicence": "` + r.MainLicence + `", "natureName": "` + r.NatureName + `", "serviceId": ` + strconv.FormatInt(r.ServiceID, 10) + `, "serviceLicence": "` + r.ServiceLicence + `", "serviceName": "` + r.ServiceName + `", "unitName": "` + r.UnitName + `", "updateRecordTime": "` + r.UpdateRecordTime + `}`
}

// AuthorizeRequest authorize request
type AuthorizeRequest struct {
	AuthKey   string `json:"authKey"`
	Timestamp string `json:"timeStamp"`
}

// String return authorizes request string
func (r *AuthorizeRequest) String() string {
	return `{"authKey": "` + r.AuthKey + `", "timeStamp": "` + r.Timestamp + `"}`
}

// AuthorizeResponse authorize response
type AuthorizeResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Params  *AuthParams `json:"params"`
}

// QueryResp is a struct for icp data
type QueryResp struct {
	IcpNumber string `json:"icp_number"`
	IcpName   string `json:"icp_name"`
	Attr      string `json:"attr"`
	Date      string `json:"date"`
}

// ParamInput request params
type ParamInput struct {
	AuthorizeRequest *AuthorizeRequest
	QueryRequest     *QueryRequest
	Path             string
	ContentType      string
}

// String request params string
func (r *ParamInput) String() string {
	if r.QueryRequest == nil {
		return `{"authorizeRequest": ` + r.AuthorizeRequest.String() + `, "path": "` + r.Path + `", "contentType": "` + r.ContentType + `"}`
	}
	if r.AuthorizeRequest == nil {
		return `{"queryRequest": ` + r.QueryRequest.String() + `, "path": "` + r.Path + `", "contentType": "` + r.ContentType + `"}`
	}

	return `{"authorizeRequest": ` + r.AuthorizeRequest.String() + `, "queryRequest": ` + r.QueryRequest.String() + `, "path": "` + r.Path + `", "contentType": "` + r.ContentType + `"}`
}
