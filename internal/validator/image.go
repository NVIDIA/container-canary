package validator

import "os/exec"

func CheckImage(image string, runtime string) bool {
	cmd := exec.Command(runtime, "image", "inspect", image)
	return cmd.Run() == nil
}
