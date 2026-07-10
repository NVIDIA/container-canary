/*
* SPDX-FileCopyrightText: Copyright (c) <2026> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package validator

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestHTTPGetCheckStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{name: "200 passes", statusCode: http.StatusOK, want: true},
		{name: "299 passes", statusCode: 299, want: true},
		{name: "301 passes", statusCode: http.StatusMovedPermanently, want: true},
		{name: "399 passes", statusCode: 399, want: true},
		{name: "400 fails", statusCode: http.StatusBadRequest, want: false},
		{name: "404 fails", statusCode: http.StatusNotFound, want: false},
		{name: "503 fails", statusCode: http.StatusServiceUnavailable, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			got, err := HTTPGetCheck(nil, httpGetProbe(server, nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHTTPGetCheckFollowsRedirect(t *testing.T) {
	tests := []struct {
		name         string
		target       string
		targetStatus int
		want         bool
	}{
		{name: "redirect to success passes", target: "/healthy", targetStatus: http.StatusOK, want: true},
		{name: "redirect to error fails", target: "/broken", targetStatus: http.StatusInternalServerError, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" {
					http.Redirect(w, r, tt.target, http.StatusFound)
					return
				}
				w.WriteHeader(tt.targetStatus)
			}))
			defer server.Close()

			got, err := HTTPGetCheck(nil, httpGetProbe(server, nil))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHTTPGetCheckResponseHeaders(t *testing.T) {
	tests := []struct {
		name          string
		responseValue string
		setHeader     bool
		want          bool
	}{
		{name: "matching header passes", responseValue: "expected", setHeader: true, want: true},
		{name: "wrong header value fails", responseValue: "unexpected", setHeader: true, want: false},
		{name: "missing header fails", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tt.setHeader {
					w.Header().Set("X-Canary-Test", tt.responseValue)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			headers := []v1.HTTPHeader{{Name: "X-Canary-Test", Value: "expected"}}
			got, err := HTTPGetCheck(nil, httpGetProbe(server, headers))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func httpGetProbe(server *httptest.Server, responseHeaders []v1.HTTPHeader) *canaryv1.Probe {
	return &canaryv1.Probe{
		HTTPGet: &canaryv1.HTTPGetAction{
			Path:                "/",
			Port:                server.Listener.Addr().(*net.TCPAddr).Port,
			ResponseHTTPHeaders: responseHeaders,
		},
	}
}
