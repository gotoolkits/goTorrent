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
	"errors"
	"fmt"
	"io"
	"io/ioutil"

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

func NewTorrentUrl(hashinfo string) {

	//hashinfo := "03621694F0E8B2CE87216C99CB5CA3AF23029E37"
	filePath := "/tmp/" + hashinfo + ".torrent"
	retcode := DownLoadTorrentFile(makeUrl(hashinfo), filePath)
	fmt.Println(retcode)

	if retcode == 0 {
		parseTorrentInfo(filePath)
	}
	//makeUrl("f8181597b51c157fb470e5ee236e364c6fbc2af2")
}

// var timeout = time.Duration(2 * time.Second)

// func dialTimeout(network, addr string) (net.Conn, error) {
// 	return net.DialTimeout(network, addr, timeout)
// }

func parseTorrentInfo(filepath string) (int, error) {

	bt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return 1, err
	}

	var metaTorrent MetaInfo
	ok := metaTorrent.ReadTorrentMetaInfoFile(bytes.NewReader(bt))
	if !ok {
		return 2, nil
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

	fmt.Println("Torrent Info Name:", name)
	fmt.Println("Torrent Hash key:", hashInfo)
	fmt.Println("Torrent Created Date:", created)
	fmt.Println("Torrent File list:", "\n", fileList)

	return 0, nil
}

func parseTorrentHash(filepath string) (string, error) {

	bt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	var metaTorrent MetaInfo
	ok := metaTorrent.ReadTorrentMetaInfoFile(bytes.NewReader(bt))
	if !ok {
		return "", errors.New("read torrent mate info file error")
	}

	return fmt.Sprintf("%X", metaTorrent.InfoHash), nil
}
