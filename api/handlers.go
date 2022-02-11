package api

import (
	"fmt"
	"io"
	"os"
	"net/http"
)

const (
	AssetDir   = "assets"
	RemoteURL  = "https://www.learningcontainer.com/wp-content/uploads/2020/04/sample-text-file.txt"
	SampleFile = "sample-data.txt"
)

func downloadFile() (err error) {
	resp, err := http.Get(RemoteURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// this relies on the dir already existing
	out, err := os.Create(fmt.Sprintf("%s/%s", AssetDir, SampleFile))
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return
}

