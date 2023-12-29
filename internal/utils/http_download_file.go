package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(request *http.Request, output string) error {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("unable to perform download request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected response from pdf download: %d != 200", response.StatusCode)
	}

	fmt.Printf("response headers: %+v\n", response.Header)

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("unable to create dest file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}
