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
	"testing"
	"time"

	canaryv1 "github.com/nvidia/container-canary/internal/apis/v1"
	"github.com/nvidia/container-canary/internal/container"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestExecuteCheckSuccess(t *testing.T) {
	require := require.New(t)
	logger, buf := logger()

	start := time.Now()
	attempts := 0
	checkFunc := func(c container.ContainerInterface, probe *canaryv1.Probe, e *zerolog.Event) (bool, error) {
		require.GreaterOrEqual(time.Since(start), (2+time.Duration(attempts))*time.Second)
		require.Less(time.Since(start), (2+time.Duration(attempts))*time.Second+100*time.Millisecond)
		attempts++
		e.Str("data", "test")
		if attempts < 5 {
			return false, nil
		}
		return true, nil
	}

	check := canaryv1.Check{
		Name: "test-check",
		Probe: canaryv1.Probe{
			InitialDelaySeconds: 2,
			PeriodSeconds:       1,
			TimeoutSeconds:      60,
			SuccessThreshold:    3,
			FailureThreshold:    10,
		},
	}

	expectedLog := "{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":1,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":2,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":3,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":4,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":5,\"data\":\"test\",\"pass\":true}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":6,\"data\":\"test\",\"pass\":true}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":7,\"data\":\"test\",\"pass\":true}\n"

	result, err := executeCheck(checkFunc, logger, &dummyContainer{}, check)
	require.NoError(err)
	require.True(result)
	require.Equal(7, attempts)
	require.Equal(expectedLog, buf.String())
}

func TestExecuteCheckFailureThreshold(t *testing.T) {
	require := require.New(t)
	logger, buf := logger()

	start := time.Now()
	attempts := 0
	checkFunc := func(c container.ContainerInterface, probe *canaryv1.Probe, e *zerolog.Event) (bool, error) {
		require.GreaterOrEqual(time.Since(start), (2+time.Duration(attempts))*time.Second)
		require.Less(time.Since(start), (2+time.Duration(attempts))*time.Second+100*time.Millisecond)
		attempts++
		e.Str("data", "test")
		if attempts < 5 {
			return false, nil
		}
		return true, nil
	}

	check := canaryv1.Check{
		Name: "test-check",
		Probe: canaryv1.Probe{
			InitialDelaySeconds: 2,
			PeriodSeconds:       1,
			TimeoutSeconds:      60,
			SuccessThreshold:    3,
			FailureThreshold:    4,
		},
	}

	expectedLog := "{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":1,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":2,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":3,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":4,\"data\":\"test\",\"pass\":false}\n"

	result, err := executeCheck(checkFunc, logger, &dummyContainer{}, check)
	require.NoError(err)
	require.False(result)
	require.Equal(4, attempts)
	require.Equal(expectedLog, buf.String())
}

func TestExecuteCheckTimeout(t *testing.T) {
	require := require.New(t)
	logger, buf := logger()

	start := time.Now()
	attempts := 0
	checkFunc := func(c container.ContainerInterface, probe *canaryv1.Probe, e *zerolog.Event) (bool, error) {
		require.GreaterOrEqual(time.Since(start), (2+time.Duration(attempts))*time.Second)
		require.Less(time.Since(start), (2+time.Duration(attempts))*time.Second+100*time.Millisecond)
		attempts++
		e.Str("data", "test")
		if attempts < 5 {
			return false, nil
		}
		return true, nil
	}

	check := canaryv1.Check{
		Name: "test-check",
		Probe: canaryv1.Probe{
			InitialDelaySeconds: 2,
			PeriodSeconds:       1,
			TimeoutSeconds:      4,
			SuccessThreshold:    3,
			FailureThreshold:    10,
		},
	}

	expectedLog := "{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":1,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":2,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":3,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":4,\"data\":\"test\",\"pass\":false}\n" +
		"{\"level\":\"info\",\"check\":\"test-check\",\"attempt\":5,\"data\":\"test\",\"pass\":true}\n"

	result, err := executeCheck(checkFunc, logger, &dummyContainer{}, check)
	require.Error(err, "check timed out after 4 seconds")
	require.False(result)
	require.Equal(5, attempts)
	require.Equal(expectedLog, buf.String())
}
