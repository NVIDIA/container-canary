package validator

import "fmt"

func getStatus(check bool, err error) string {
	var status string
	if err != nil {
		status = fmt.Sprintf("error - %s", err.Error())
	} else {
		status = fmt.Sprintf("%t", check)
	}
	return status
}
