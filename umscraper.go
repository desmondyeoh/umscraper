package umscraper

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

const rootURL = "https://www.um.edu.my/"

// ScrapePtjJabTable concurrently scrapes full list of (Ptj, Jab) pairs.
// Returns the table with headers [PtjCode, PtjText, JabCode, JabText].
func ScrapePtjJabTable(sleepBetweenPtj time.Duration) [][]string {
	// Scrape Ptj list.
	ptjMap := ScrapePtjMap()

	// Scrape Ptj pages for jabatan list.
	c := make(chan [][]string)
	data := [][]string{}
	for ptjCode, ptjText := range ptjMap {
		go ScrapeJabTable(ptjCode, ptjText, data, c)
		time.Sleep(sleepBetweenPtj)
	}

	// Create table with header
	table := [][]string{}
	header := []string{"ptjCode", "ptjText", "jabCode", "jabText"}
	table = append(table, header)

	// Append table body
	for range ptjMap {
		table = append(table, <-c...)
	}

	return table
}

// ScrapeStaffTable concurrently scrapes all staff's profiles.
// Return a table with headers [PtjCode, PtjText, JabCode, JabText, Name, NameQueryEsc, Details...]
func ScrapeStaffTable(ptjJabTable [][]string, sleepBetweenPtjJab time.Duration) [][]string {
	outTable := [][]string{}

	c := make(chan [][]string)
	for _, row := range ptjJabTable[1:] {
		ptjCode := row[0]
		ptjText := row[1]
		jabCode := row[2]
		jabText := row[3]
		go ScrapePtjJabSearchRes(ptjCode, ptjText, jabCode, jabText, c)
		time.Sleep(sleepBetweenPtjJab)
	}

	for range ptjJabTable[1:] {
		outTable = append(outTable, <-c...)
	}

	return outTable
}

// ScrapePtjMap scrapes full list of Ptj.
// Returns a map of { PtjCode: PtjText }
func ScrapePtjMap() map[string]string {
	fmt.Println("Scraping PtjMap...")
	resp, err := soup.Get(rootURL + "list-staff.php")
	checkError("Unable to ScrapePtj", err)

	// Parse the HTML to get all departments
	doc := soup.HTMLParse(resp)
	options := doc.Find("select", "id", "kodPTJ").FindAll("option")

	// Soup -> ptjMap
	ptjMap := map[string]string{}
	for _, jab := range options {
		ptjCode := jab.Attrs()["value"]
		ptjText := jab.Text()
		if ptjCode == "" && ptjText == "Select PTJ / Faculty" {
			continue
		}
		ptjMap[ptjCode] = ptjText
	}

	return ptjMap
}

// ScrapeJabTable scrapes full list of Ptj and Jab for a Ptj.
// Returns a table with headers [PtjCode, PtjText, JabCode, JabText].
func ScrapeJabTable(ptjCode string, ptjText string, data [][]string, c chan [][]string) {
	fmt.Println("Scraping JabTable...", ptjCode, ptjText)

	// Scrape ptjCode
	resp, err := soup.Get(rootURL + "list-staff.php?kodPTJ=" + ptjCode)
	checkError("Unable to ScrapeJabatan", err)

	// Parse the HTML to get all departments
	doc := soup.HTMLParse(resp)
	options := doc.Find("select", "id", "kodJAB").FindAll("option")

	// Soup -> subTable
	subTable := [][]string{}
	for _, jab := range options {
		jabCode := jab.Attrs()["value"]
		jabText := jab.Text()
		if jabCode == "" && jabText == "Select Department" {
			continue
		}
		subTable = append(subTable, []string{ptjCode, ptjText, jabCode, jabText})
	}

	// Return departments with Ptj code in a table
	c <- subTable
}

// ScrapePtjJabSearchRes scrapes full list of staffs for a (ptjCode, jabCode) pair
// Return a subtable with headers [PtjCode, PtjText, JabCode, JabText, Name, NameQueryEsc, Details...]
func ScrapePtjJabSearchRes(ptjCode string, ptjText string, jabCode string, jabText string, c chan [][]string) {
	fmt.Println("Scraping PtjJabSearchRes...", ptjCode, jabCode)
	resp, err := soup.Get(rootURL + "list-staff.php?kodPTJ=" + ptjCode + "&kodJAB=" + jabCode)
	checkError("Unable to ScrapeProfiles "+ptjCode+" "+jabCode, err)

	// Parse the HTML to get all departments
	subTables := ParsePtjJabSearchResHTML(ptjCode, ptjText, jabCode, jabText, resp)
	c <- subTables
}

// ParsePtjJabSearchResHTML parses HTML string and extracts staff's Name & Details
// Returns subTable with headers [PtjCode, PtjText, JabCode, JabText, Name, NameQueryEsc, Details...]
func ParsePtjJabSearchResHTML(ptjCode string, ptjText string, jabCode string, jabText string, text string) [][]string {
	doc := soup.HTMLParse(text)
	staffs := doc.Find("table", "id", "sample_3").Find("tbody").FindAll("tr")

	subTable := [][]string{}

	for _, staff := range staffs {
		name := staff.Find("div", "class", "nama-staf")
		img := staff.Find("img", "class", "avatar")
		details := staff.FindAll("div", "class", "nama-style")

		nameText := name.Text()
		nameTextQueryEsc := url.QueryEscape(nameText)

		// add row
		row := []string{ptjCode, ptjText, jabCode, jabText, nameText, nameTextQueryEsc}
		for _, detail := range details {
			row = append(row, strings.TrimSpace(detail.FullText()))
		}
		subTable = append(subTable, row)

		// save image
		imgDir := "data/img/" + ptjCode + "/"
		CreateDirectoryIfNotExist(imgDir, 0700)
		imgSrc := img.Attrs()["src"]
		imgBase64 := ExtractImgBase64(imgSrc)
		buf := DecodeImgBase64(imgBase64)
		WriteFile(buf, filepath.Join(imgDir, nameTextQueryEsc+".jpg"))
	}

	return subTable
}
