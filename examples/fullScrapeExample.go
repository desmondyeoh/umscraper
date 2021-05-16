package main

import (
	"time"

	"github.com/desmondyeoh/umscraper"
)

// This example shows how to scrape the data fully.

func main() {
	// Constants
	// Sending too many requests without sufficient sleep time might cause IP to be blacklisted
	// or some requests not able to get response successfullly.
	sleepBetweenPtj := 100 * time.Millisecond
	sleepBetweenPtjJab := 1000 * time.Millisecond

	// Stage 0 - Create a data directory to store scraped data
	umscraper.CreateDirectoryIfNotExist("data", 0700)

	// Stage 1 - Scrape full list of (ptjCode, ptjText, jabCode, jabText)
	ptjJabTable := umscraper.ScrapePtjJabTable(sleepBetweenPtj)
	umscraper.WriteCsv(ptjJabTable, "data/ptjJabTable.csv")

	// Stage 2 - Scrape full list of (ptjCode, ptjText, jabCode, jabText, staffName, staffDetails...)
	ptjJabTable = umscraper.ReadCsv("data/ptjJabTable.csv")
	staffTable := umscraper.ScrapeStaffTable(ptjJabTable, sleepBetweenPtjJab)
	umscraper.WriteCsv(staffTable, "data/staffTable.csv")

}
