package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Get HTTP GET request
func Get(url string) (int, []byte) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		_ = fmt.Errorf("error sending GET request, url: %s, %q", url, err)
		return http.StatusInternalServerError, nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("error sending GET request, url: %s, %q", url, err)
		}
	}(resp.Body)
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			_ = fmt.Errorf("error sending GET request, url: %s, %q", url, err)
		}
	}
	return resp.StatusCode, result.Bytes()
}
