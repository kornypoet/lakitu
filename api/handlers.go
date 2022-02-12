package api

import (
	"errors"
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

func downloadAction() (err error) {
	if !assetExists() {
		err = downloadFile()
		return
	}
	return errors.New("file already downloaded")
}

func assetExists() bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", AssetDir, SampleFile)); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func downloadFile() (err error) {
	resp, err := http.Get(RemoteURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	os.MkdirAll(AssetDir, 0700)
	out, err := os.OpenFile(fmt.Sprintf("%s/%s", AssetDir, SampleFile), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return
}
