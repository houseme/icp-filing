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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func TestICP_Md5(t *testing.T) {
	type fields struct {
		token string
		ip    string
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test",
			fields: fields{
				token: "0",
				ip:    "127.0.0.1",
			},
			args: args{
				str: "test",
			},
			want: "098f6bcd4621d373cade4e832627b4f6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Filling{
				token: tt.fields.token,
				ip:    tt.fields.ip,
			}
			if got := i.md5(tt.args.str); got != tt.want {
				t.Errorf("md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestICP_String(t *testing.T) {
	type fields struct {
		token string
		ip    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test",
			fields: fields{token: defaultToken, ip: "101,110,123,124"},
			want:   `{"ip":"101,110,123,124","token":"0"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Filling{
				token: tt.fields.token,
				ip:    tt.fields.ip,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestICP_authorize(t *testing.T) {
	type fields struct {
		token string
		ip    string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				token: "0",
				ip:    "101.123.124.119",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Filling{
				token: tt.fields.token,
				ip:    tt.fields.ip,
			}
			initLog(os.TempDir(), hlog.LevelInfo)
			fmt.Println("icp:", i)
			if err := i.authorize(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
