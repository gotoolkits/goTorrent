package torrent

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gotoolkits/gorequest"
)

func DownLoadTorrentFile(url, filePath string) int {

	request := gorequest.New()
	resp, body, errs := request.Get(url).
		Retry(3, 5*time.Second, http.StatusGatewayTimeout, http.StatusInternalServerError, http.StatusServiceUnavailable).
		End()

	defer resp.Body.Close()
	if errs != nil {
		fmt.Println("get torrent file failed:", errs)
		return 1
	}

	if resp.StatusCode != 200 {
		fmt.Println("get torrent status:", resp.Status)
		return 2
	}

	ioutil.WriteFile(filePath, []byte(body), 0666)
	return 0
}

func logFile(msg string) {
	f, err := os.OpenFile("logfile_torrent.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(msg)
}
