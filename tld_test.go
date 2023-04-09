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
	"os"
	"reflect"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/houseme/icp-filing/tld"
)

func TestFilling_DomainTLD(t *testing.T) {
	type args struct {
		ctx   context.Context
		link  string
		level int
	}

	var (
		ctx = context.Background()
		f   = New(ctx, WithLogLevel(hlog.LevelInfo), WithLogPath(os.TempDir()))
	)

	tests := []struct {
		name     string
		args     args
		wantResp *tld.DomainTLDResp
		wantErr  bool
	}{
		{
			name:     "TestFilling_DomainTLD",
			args:     args{ctx: ctx, link: "https://www.baidu.com", level: 0},
			wantErr:  false,
			wantResp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := f.DomainTLD(tt.args.ctx, tt.args.link, tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainTLD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("DomainTLD() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
