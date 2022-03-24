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
	"github.com/nvidia/container-canary/internal/container"
	v1 "k8s.io/api/core/v1"
)

func ExecCheck(c container.ContainerInterface, action *v1.ExecAction) (bool, error) {
	_, err := c.Exec(action.Command...)
	if err != nil {
		return false, nil
	}
	return true, nil
}
