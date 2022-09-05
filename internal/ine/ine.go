package ine

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/nleeper/goment"
	"github.com/xuri/excelize/v2"
)

var shortMonths = [12]string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}
var longMonths = [12]string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Setiembre", "Octubre", "Noviembre", "Diciembre"}

const xlsFileUrl = "http://www.ine.gub.uy/c/document_library/get_file?uuid=1dcbe20a-153b-4caf-84a7-7a030d109471"

type DolarPrice struct {
	Date           time.Time
	BuyPrice       float64
	SellPrice      float64
	EbrouBuyPrice  float64
	EbrouSellPrice float64
}

func GetIneDollarPrices() ([]DolarPrice, error) {
	file, err := downloadXlsFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	prices, err := parseXlsFile(file)
	if err != nil {
		return nil, err
	}

	return prices, err
}

func parsePrice(price string) float64 {
	value, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return 0
	}
	return value
}

func parseXlsFile(file io.Reader) ([]DolarPrice, error) {
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)

	if err != nil {
		return nil, err
	}

	// skip first 8 rows
	entries := make([]DolarPrice, 0)

	// Initialize the date to 30 Dic 1999
	currentDate, err := goment.New("1999-12-30")
	currentDate.SetLocale("es")

	if err != nil {
		return nil, err
	}

	for _, row := range rows[8:] {
		// Stop when we reach the end of the table
		if len(row) == 0 {
			break
		}

		// Get the values to create a Date
		day := row[0]
		month := row[1]
		year := row[2]

		// Some rows have a year, if so, update the current date with the new year
		if year != "" {
			yearAsNumber, _ := strconv.ParseInt(year, 10, 32)
			currentDate = currentDate.SetYear(int(yearAsNumber))
		}

		// Some rows have a month, if so, update the current date with the new month
		if month != "" {
			for i, shortMonth := range shortMonths {
				if shortMonth == month {
					currentDate.SetMonth(i + 1)
					break
				}
			}
			for i, longMonth := range longMonths {
				if longMonth == month {
					currentDate.SetMonth(i + 1)
					break
				}
			}
		}

		// Some rows have a day, if so, update the current date with the new day
		if day != "" {
			dayAsNumber, _ := strconv.ParseInt(day, 10, 32)
			currentDate.SetDate(int(dayAsNumber))
		}

		dollarEntry := DolarPrice{
			Date:           currentDate.ToTime(),
			BuyPrice:       parsePrice(row[3]),
			SellPrice:      parsePrice(row[4]),
			EbrouBuyPrice:  parsePrice(row[6]),
			EbrouSellPrice: parsePrice(row[7]),
		}

		entries = append(entries, dollarEntry)
	}

	return entries, nil
}

func downloadXlsFile() (io.ReadCloser, error) {
	resp, err := http.Get(xlsFileUrl)
	if err != nil {
		return nil, err
	}
	return resp.Body, err

}
