/*
* SPDX-FileCopyrightText: Copyright (c) <2022> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
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
	"fmt"
	"net/http"
	"strings"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/nvidia/container-canary/internal/container"
)

func HTTPGetCheck(c container.ContainerInterface, probe *canaryv1.Probe) (bool, error) {
	action := probe.HTTPGet
	// Match the Kubernetes kubelet: follow redirects and treat any 2xx or 3xx
	// (200 <= code < 400) response as a success.
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d%s", action.Port, action.Path), nil)
	if err != nil {
		return false, nil
	}
	req.Close = true

	for _, header := range action.HTTPHeaders {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return false, nil
	}
	for _, header := range action.ResponseHTTPHeaders {
		val, ok := resp.Header[header.Name]
		if !ok || header.Value != strings.Join(val[:], "") {
			return false, nil
		}
	}
	return true, nil
}
