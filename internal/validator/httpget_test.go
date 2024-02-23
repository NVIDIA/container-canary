/*
* SPDX-FileCopyrightText: Copyright (c) <2024> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
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
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func serveHTTP(t *testing.T) *http.Server {
	require := require.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/success", func(writer http.ResponseWriter, request *http.Request) {
		require.Equal("ping", request.Header.Get("X-Container-Canary-Request"))
		writer.Header().Add("X-Container-Canary-Response", "pong")
		writer.Header().Add("X-Container-Canary-Extra", "extra")
		writer.WriteHeader(200)
	})
	mux.HandleFunc("/failure/status-code", func(writer http.ResponseWriter, request *http.Request) {
		require.Equal("ping", request.Header.Get("X-Container-Canary-Request"))
		writer.Header().Add("X-Container-Canary-Response", "pong")
		writer.Header().Add("X-Container-Canary-Extra", "extra")
		writer.WriteHeader(403)
	})
	mux.HandleFunc("/failure/wrong-header", func(writer http.ResponseWriter, request *http.Request) {
		require.Equal("ping", request.Header.Get("X-Container-Canary-Request"))
		writer.Header().Add("X-Container-Canary-Response", "no")
		writer.Header().Add("X-Container-Canary-Extra", "extra")
		writer.WriteHeader(200)
	})
	mux.HandleFunc("/failure/no-header", func(writer http.ResponseWriter, request *http.Request) {
		require.Equal("ping", request.Header.Get("X-Container-Canary-Request"))
		writer.Header().Add("X-Container-Canary-Extra", "extra")
		writer.WriteHeader(200)
	})

	server := &http.Server{
		Handler: mux,
	}
	l, err := net.Listen("tcp", "localhost:54321")
	require.NoError(err)
	go func() {
		err := server.Serve(l)
		require.IsType(http.ErrServerClosed, err)
	}()

	return server
}

func shutdownHTTP(t *testing.T, server *http.Server) {
	require := require.New(t)

	err := server.Shutdown(context.Background())
	require.NoError(err)
}

func TestHTTPGetSuccess(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		HTTPGet: &canaryv1.HTTPGetAction{
			Path: "/success",
			Port: 54321,
			HTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Request",
					Value: "ping",
				},
			},
			ResponseHTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Response",
					Value: "pong",
				},
			},
		},
	}

	server := serveHTTP(t)
	defer shutdownHTTP(t, server)

	logger, buf := logger()
	e := logger.Info()
	result, err := HTTPGetCheck(c, probe, e)
	e.Send()

	require.True(result)
	require.NoError(err)
	d := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	var je map[string]interface{}
	require.NoError(d.Decode(&je))
	require.Equal(200.0, je["status"])
	require.Equal("pong", je["headers"].(map[string]interface{})["X-Container-Canary-Response"])
	require.Equal("extra", je["headers"].(map[string]interface{})["X-Container-Canary-Extra"])
}

func TestHTTPGetStatusCodeFailure(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		HTTPGet: &canaryv1.HTTPGetAction{
			Path: "/failure/status-code",
			Port: 54321,
			HTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Request",
					Value: "ping",
				},
			},
			ResponseHTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Response",
					Value: "pong",
				},
			},
		},
	}

	server := serveHTTP(t)
	defer shutdownHTTP(t, server)

	logger, buf := logger()
	e := logger.Info()
	result, err := HTTPGetCheck(c, probe, e)
	e.Send()

	require.False(result)
	require.NoError(err)
	d := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	var je map[string]interface{}
	require.NoError(d.Decode(&je))
	require.Equal(403.0, je["status"])
	require.Equal("pong", je["headers"].(map[string]interface{})["X-Container-Canary-Response"])
	require.Equal("extra", je["headers"].(map[string]interface{})["X-Container-Canary-Extra"])
}

func TestHTTPGetWrongHeaderFailure(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		HTTPGet: &canaryv1.HTTPGetAction{
			Path: "/failure/wrong-header",
			Port: 54321,
			HTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Request",
					Value: "ping",
				},
			},
			ResponseHTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Response",
					Value: "pong",
				},
			},
		},
	}

	server := serveHTTP(t)
	defer shutdownHTTP(t, server)

	logger, buf := logger()
	e := logger.Info()
	result, err := HTTPGetCheck(c, probe, e)
	e.Send()

	require.False(result)
	require.NoError(err)
	d := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	var je map[string]interface{}
	require.NoError(d.Decode(&je))
	require.Equal(200.0, je["status"])
	require.Equal("no", je["headers"].(map[string]interface{})["X-Container-Canary-Response"])
	require.Equal("extra", je["headers"].(map[string]interface{})["X-Container-Canary-Extra"])
}

func TestHTTPGetNoHeaderFailure(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		HTTPGet: &canaryv1.HTTPGetAction{
			Path: "/failure/no-header",
			Port: 54321,
			HTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Request",
					Value: "ping",
				},
			},
			ResponseHTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Response",
					Value: "pong",
				},
			},
		},
	}

	server := serveHTTP(t)
	defer shutdownHTTP(t, server)

	logger, buf := logger()
	e := logger.Info()
	result, err := HTTPGetCheck(c, probe, e)
	e.Send()

	require.False(result)
	require.NoError(err)
	d := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	var je map[string]interface{}
	require.NoError(d.Decode(&je))
	require.Equal(200.0, je["status"])
	_, ok := je["headers"].(map[string]interface{})["X-Container-Canary-Response"]
	require.False(ok)
	require.Equal("extra", je["headers"].(map[string]interface{})["X-Container-Canary-Extra"])
}

func TestHTTPGetNoServerFailure(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		HTTPGet: &canaryv1.HTTPGetAction{
			Path: "/success",
			Port: 54321,
			HTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Request",
					Value: "ping",
				},
			},
			ResponseHTTPHeaders: []v1.HTTPHeader{
				{
					Name:  "X-Container-Canary-Response",
					Value: "pong",
				},
			},
		},
	}

	logger, buf := logger()
	e := logger.Info()
	result, err := HTTPGetCheck(c, probe, e)
	e.Send()

	require.False(result)
	require.NoError(err)
	require.Equal("{\"level\":\"info\"}\n", buf.String())
}
