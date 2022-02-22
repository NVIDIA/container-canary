package checks

import (
	"github.com/jacobtomlinson/containercanary/internal/containertools/container"
	v1 "k8s.io/api/core/v1"
)

func ExecCheck(c container.Container, action *v1.ExecAction) (bool, error) {
	_, err := c.Exec(action.Command...)
	if err != nil {
		return false, nil
	}
	return true, nil
}
