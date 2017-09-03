//
//magnet to torrent
//
//03621694F0E8B2CE87216C99CB5CA3AF23029E37
//        ||
//        ||  base32(hashinfo+`xxd -r -p`) or
//        ||  base32(StrToBinaryEncode(hashinfo))
//        \/
//ANRBNFHQ5CZM5BZBNSM4WXFDV4RQFHRX
package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	//_ "github.com/go-sql-driver/mysql"

	bencode "github.com/jackpal/bencode-go"
)

type FileDict struct {
	Length int64    "length"
	Path   []string "path"
	Md5sum string   "md5sum"
}

type InfoDict struct {
	FileDuration []int64 "file-duration"
	FileMedia    []int64 "file-media"

	// Single file
	Name   string "name"
	Length int64  "length"
	Md5sum string "md5sum"

	// Multiple files
	Files       []FileDict "files"
	PieceLength int64      "piece length"
	Pieces      string     "pieces"
	Private     int64      "private"
}

type MetaInfo struct {
	Info         InfoDict   "info"
	InfoHash     string     "info hash"
	Announce     string     "announce"
	AnnounceList [][]string "announce-list"
	CreationDate int64      "creation date"
	Comment      string     "comment"
	CreatedBy    string     "created by"
	Encoding     string     "encoding"
}

func (metaInfo *MetaInfo) ReadTorrentMetaInfoFile(r io.Reader) bool {

	fileMetaData, er := bencode.Decode(r)
	if er != nil {
		return false
	}

	metaInfoMap, ok := fileMetaData.(map[string]interface{})
	if !ok {
		return false
	}

	var bytesBuf bytes.Buffer
	for mapKey, mapVal := range metaInfoMap {
		switch mapKey {
		case "info":
			if er = bencode.Marshal(&bytesBuf, mapVal); er != nil {
				return false
			}

			infoHash := sha1.New()
			infoHash.Write(bytesBuf.Bytes())
			metaInfo.InfoHash = string(infoHash.Sum(nil))

			if er = bencode.Unmarshal(&bytesBuf, &metaInfo.Info); er != nil {
				return false
			}

		case "announce-list":
			if er = bencode.Marshal(&bytesBuf, mapVal); er != nil {
				return false
			}
			if er = bencode.Unmarshal(&bytesBuf, &metaInfo.AnnounceList); er != nil {
				return false
			}

		case "announce":
			if aa, ok := mapVal.(string); ok {
				metaInfo.Announce = aa
			}

		case "creation date":

			if tt, ok := mapVal.(int64); ok {
				metaInfo.CreationDate = tt
			}

		case "comment":
			if cc, ok := mapVal.(string); ok {
				metaInfo.Comment = cc
			}

		case "created by":
			if cb, ok := mapVal.(string); ok {
				metaInfo.CreatedBy = cb
			}

		case "encoding":
			if ed, ok := mapVal.(string); ok {
				metaInfo.Encoding = ed
			}
		}
	}

	return true
}

func NewTorrentUrl() {

	url := makeUrl("03621694F0E8B2CE87216C99CB5CA3AF23029E37")

	fmt.Println(url)
	//makeUrl("f8181597b51c157fb470e5ee236e364c6fbc2af2")

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

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func pullTorrent(url string) (int, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 1, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Host", "bt.box.n0808.com")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "Keep-Alive")

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := &http.Client{
		Transport: &transport,
	}

	resp, err := client.Do(req)

	if err != nil {
		return 2, err
	}
	defer resp.Body.Close()

	var metaTorrent MetaInfo
	ok := metaTorrent.ReadTorrentMetaInfoFile(resp.Body)
	if !ok {
		return 3, nil
	}

	name := metaTorrent.Info.Name
	hashInfo := fmt.Sprintf("%X", metaTorrent.InfoHash)
	created := metaTorrent.CreationDate

	var fileLength int64
	var fileDownLoadList bytes.Buffer
	var fileList string

	for _, fileDict := range metaTorrent.Info.Files {
		fileLength += fileDict.Length
		for _, path := range fileDict.Path {
			fileDownLoadList.WriteString(path)
			fileDownLoadList.WriteString("\r\n")
		}
	}
	fileList = fileDownLoadList.String()

	fmt.Println(name)
	fmt.Println(hashInfo)
	fmt.Println(created)
	fmt.Println(fileList)

	return 0, nil
}
