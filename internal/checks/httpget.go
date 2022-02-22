package checks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jacobtomlinson/containercanairy/internal/apis/config"
	"github.com/jacobtomlinson/containercanairy/internal/containertools/container"
)

func HTTPGetCheck(c container.Container, action *config.HTTPGetAction) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d%s", action.Port, action.Path), nil)
	if err != nil {
		fmt.Println(err.Error())
		return false, nil
	}
	req.Close = true

	for _, header := range action.HTTPHeaders {
		req.Header.Set(header.Name, header.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false, nil
	}
	for _, header := range action.ResponseHTTPHeaders {
		if val, ok := resp.Header[header.Name]; ok {
			if header.Value != strings.Join(val[:], "") {
				return false, nil
			}
		}
	}
	defer resp.Body.Close()
	return true, nil
}
