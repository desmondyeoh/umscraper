package umscraper

import "log"

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
