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
	"net"
	"sync"
	"testing"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestTCPSocketSuccess(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		TCPSocket: &canaryv1.TCPSocketAction{
			Port: 54321,
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	listener, err := net.Listen("tcp", "localhost:54321")
	require.NoError(err)
	success := false
	go func() {
		conn, err := listener.Accept()
		require.NoError(err)
		conn.Close()
		success = true
		wg.Done()
	}()

	logger, buf := logger()
	e := logger.Info()
	result, err := TCPSocketCheck(c, probe, e)
	e.Send()

	require.True(result)
	require.NoError(err)
	require.Equal("{\"level\":\"info\"}\n", buf.String())
	wg.Wait()
	require.True(success)
}

func TestTCPSocketFailure(t *testing.T) {
	require := require.New(t)

	c := &dummyContainer{}
	probe := &canaryv1.Probe{
		TCPSocket: &canaryv1.TCPSocketAction{
			Port: 54321,
		},
	}

	logger, buf := logger()
	e := logger.Info()
	result, err := TCPSocketCheck(c, probe, e)
	e.Send()

	require.False(result)
	require.NoError(err)
	require.Equal("{\"level\":\"info\"}\n", buf.String())
}
