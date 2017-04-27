package finance

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadStockValues(filePath string, ch chan AssetValue) {
	defer close(ch)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// FIXME: Could we make the following code a bit cleaner?
		evaluatedAt, _ := time.Parse("2006-1-2", record[0])
		open, _ := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		high, _ := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		low, _ := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		close, _ := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		volume, _ := strconv.ParseInt(strings.TrimSpace(record[5]), 10, 64)

		ch <- AssetValue{
			EvaluatedAt: evaluatedAt,
			Granularity: DAY,
			Open:        open,
			High:        high,
			Low:         low,
			Close:       close,
			Volume:      volume,
		}
	}
}

// FIXME: Need to make an abstraction layer for grom.DB to avoid strong binds
// with the gorm library
func InsertStockSymbol(db *gorm.DB, symbol string) {
	// FIXME: Separate the following part into a different function (checking
	// for an existing Asset record)
	var asset Asset
	db.First(&asset, "name = ?", symbol)

	if asset == (Asset{}) {
		asset := Asset{
			Name: symbol,
		}
		db.Create(&asset)
	}
}

func ImportStockValues(filePath string, symbol string) {
	db := ConnectDatabase()
	defer db.Raw.Close()

	InsertStockSymbol(db.Raw, symbol)

	ch := make(chan AssetValue)
	go ReadStockValues(filePath, ch)
	for v := range ch {
		fmt.Println("Processing", v)
		db.Raw.Create(&v)
	}
}
