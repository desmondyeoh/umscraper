package umscraper

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"strings"
)

func ExtractImgBase64(imgSrc string) string {
	i := strings.Index(imgSrc, ",")
	if i < 0 {
		log.Fatal("no comma in imgSrc")
	}
	return imgSrc[i+1:]
}

func DecodeImgBase64(imgBase64 string) []byte {
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgBase64))
	buf, err := ioutil.ReadAll(dec)
	checkError("ioutil ReadAll decoded image into buffer error", err)
	return buf
}
