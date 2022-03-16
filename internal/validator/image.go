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
	"os/exec"

	"github.com/spf13/cobra"
)

func CheckImage(cmd *cobra.Command, image string, runtime string) bool {
	command := exec.Command(runtime, "image", "inspect", image)
	err := command.Run()

	if err == nil {
		return true
	} else {
		cmd.Printf("Cannot find %s, pulling...", image)
		command = exec.Command(runtime, "pull", image)
		return command.Run() == nil
	}
}
