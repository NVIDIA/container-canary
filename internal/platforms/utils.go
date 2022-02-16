package platforms

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

func all(bools []bool) bool {
	for _, b := range bools {
		if !b {
			return false
		}
	}
	return true
}
