package torrent

import (
	"fmt"
	"strings"
)

const (
	SAMPLE = "magnet:?xt=urn:btih:03621694F0E8B2CE87216C99CB5CA3AF23029E37"
)

type magnet struct {
	btFilePath string
	btPrefix   string
	btHashInfo string
	mURI       string
}

var mg magnet

func InitMagnet() {
	mg.btPrefix = "magnet:?xt=urn:btih:"
}

func (m *magnet) ExtractHashInfoFromMagnet() string {
	strslic := strings.Split(m.mURI, ":")
	return strslic[len(strslic)-1]
}

func (m *magnet) GenerateMagnetHashInfo() string {

	hashkey, err := parseTorrentHash(m.btFilePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hashkey

}

func (m *magnet) CreateMangnetURI() string {
	return m.btPrefix + m.GenerateMagnetHashInfo()
}
