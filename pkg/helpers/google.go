package helpers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func DownloadImage(url string) ([]byte, string, int64, string, error) {
	// Make a GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return nil, "", 0, "", err
	}
	defer response.Body.Close()

	// Read the response body into a byte slice
	imageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", 0, "", err
	}

	// Get the file name from the URL
	fileName := filepath.Base(url)

	// Get the file size
	fileSize := int64(len(imageBytes))

	// Get the file type
	fileType := http.DetectContentType(imageBytes)

	return imageBytes, fileName, fileSize, fileType, nil
}
