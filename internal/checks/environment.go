package checks

import (
	"fmt"
	"strings"

	"github.com/jacobtomlinson/containercanairy/internal/containertools/container"
)

func CheckUser(c container.Container, user string) (bool, error) {
	out, err := c.Exec("whoami")
	return out == user, err
}

func CheckUID(c container.Container, uid string) (bool, error) {
	out, err := c.Exec("id")
	return strings.Contains(out, fmt.Sprintf("uid=%s", uid)), err
}

func CheckPWD(c container.Container, home string) (bool, error) {
	out, err := c.Exec("pwd")
	return strings.Contains(out, home), err
}
