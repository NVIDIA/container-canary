// Copyright (c) 2022, NVIDIA CORPORATION.

package terminal

import (
	"fmt"
)

func PrintCheckItem(emoji string, description string, status string) {
	fmt.Printf("%s %-50s [%s]\n", emoji, description, status)
}
