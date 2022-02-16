package terminal

import (
	"fmt"
)

func PrintCheckItem(emoji string, description string, status string) {
	fmt.Printf("%s %s [%s]\n", emoji, description, status)
}
