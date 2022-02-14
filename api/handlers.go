package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var AssetDir string

const (
	RemoteURL  = "https://www.learningcontainer.com/wp-content/uploads/2020/04/sample-text-file.txt"
	SampleFile = "sample-data.txt"
)

func AssetFile() string {
	return fmt.Sprintf("%s/%s", AssetDir, SampleFile)
}

func assetExists() bool {
	if _, err := os.Stat(AssetFile()); os.IsNotExist(err) {
		return false
	}
	return true
}

func downloadFile() (err error) {
	log.Debugf("Downloading remote asset %s", RemoteURL)
	resp, err := http.Get(RemoteURL)
	if err != nil {
		log.Errorf("Error downloading remote file %s: %s", RemoteURL, err)
		return
	}
	defer resp.Body.Close()

	log.Debugf("Opening asset file %s", AssetFile())
	out, err := os.OpenFile(AssetFile(), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Errorf("Error opening file %s: %s", AssetFile(), err)
		return
	}
	defer out.Close()

	log.Debugf("Writing asset file %s", AssetFile())
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Errorf("Error writing file %s: %s", AssetFile(), err)
		return
	}

	return
}
