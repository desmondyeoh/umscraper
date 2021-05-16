# umscraper

**UM Scraper scrapes all staffs names, ptj (faculty), jab (department), details and images**

*umscraper* is a small GoLang web scraper package to scrape all UM's staffs info.

Exported variables and functions implemented till now :
```go
func ScrapePtjJabTable(time.Duration) [][]string {} // Concurrently scrapes full list of (Ptj, Jab) pairs. Returns the table with headers [PtjCode, PtjText, JabCode, JabText].
func ScrapeStaffTable([][]string time.Duration) [][]string {} // Concurrently scrapes all staff's profiles. Return a table with headers [PtjCode, PtjText, JabCode, JabText, Name, NameQueryEsc, Details...]
func WriteCsv([][]string, string) {} // Writes a table of [][]string into specified filename.
func ReadCsv(string) [][]string {} // Reads a table of [][]string from specified filename.
```

## Dataset
The scraped data are in this repository [desmondyeoh/umscraper-data](https://github.com/desmondyeoh/umscraper-data)
- `ptjJabTable.csv` has the full list of (Ptj, Jab) pairs with columns {ptjCode, ptjText, jabCode, jabText}
- `staffTable.csv` has the full staff list with columns {ptjCode, ptjText, jabCode, jabText, name, nameQueryEsc, details...}
  - `img/<ptjCode>/<nameQueryEsc>.jpg` can be used to find the staff's respective image file.
- `img/` has all the staff images files, organised into ptjCode folders.

## Installation
Install the package using the command
```bash
go get github.com/desmondyeoh/umscraper
```

## Usage
Check out the `examples/` folder for a quick getting started example!

## Contributions
This package was developed in my free time. However, contributions from everybody in the community are welcome, to make it a better web scraper. If you think there should be a particular feature or function included in the package, feel free to open up a new issue or pull request.
