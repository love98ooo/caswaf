// Copyright 2024 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rule

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/casbin/caswaf/object"
	"github.com/casbin/caswaf/util"
)

type IpRule struct{}

func (r *IpRule) checkRule(expressions []*object.Expression, req *http.Request) (bool, string, error) {
	clientIp := util.GetClientIp(req)
	for _, expression := range expressions {
		reason := fmt.Sprintf("expression matched: \"%s %s %s\"", clientIp, expression.Operator, expression.Value)
		ips := strings.Split(expression.Value, " ")
		for _, ip := range ips {
			if strings.Contains(ip, "/") {
				_, ipNet, err := net.ParseCIDR(ip)
				if err != nil {
					return false, "", err
				}
				if ipNet.Contains(net.ParseIP(clientIp)) == (expression.Operator == "is in") {
					return true, reason, nil
				}
			} else if strings.ContainsAny(ip, ".:") {
				if (ip == clientIp) == (expression.Operator == "is in") {
					return true, reason, nil
				}
			} else {
				return false, "", fmt.Errorf("invalid ip: %s", ip)
			}
		}
	}
	return false, "", nil
}