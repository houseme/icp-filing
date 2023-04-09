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

package tld

import (
	"context"
	"fmt"
	"testing"
)

var testUrls = []string{"www.google.com.hk", "www.discuz.net", "com",
	"www.discuz.vip", "www.ritto.shiga.jp", "ritto.shiga.jp", "mp.weixin.qq.com", "jonsen.yang.cn"}

func TestGetTld(t *testing.T) {
	// url := "www.google.com.hk"
	ctx := context.Background()
	for _, url := range testUrls {
		ss, dd, tld := GetSubdomain(ctx, url, 2)
		fmt.Println("GetSubdomain:", ss, dd, tld)
		resp, err := GetTLD(ctx, url, 0)
		fmt.Println("GetTld: ", resp, err)
		if nil != err {
			t.Error("Failed get TLD:" + err.Error())
			return
		}
		t.Logf("resp：%s: %v, %s\n", url, resp.Tld, resp.Domain)
	}

	// t.Fail()
}

func BenchmarkGetTld(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		GetTLD(ctx, "www.aaa.bbb.ccc.ddd.forease.com.cn", 0)
	}
}

func BenchmarkGetSubdomain(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		GetSubdomain(ctx, "www.aaa.bbb.ccc.ddd.forease.com.cn", 0)
	}
}
