package finance

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadCSV(filePath string, ch chan []string) {
	defer close(ch)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(f))
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		ch <- row
	}
}

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

func ImportAccounts(filePath string) error {
	db := ConnectDatabase()
	defer db.Raw.Close()

	ch := make(chan []string)
	go ReadCSV(filePath, ch)
	for row := range ch {
		accountID, _ := strconv.ParseUint(strings.TrimSpace(row[0]), 10, 64)
		name := strings.TrimSpace(row[1])

		fmt.Printf("Creating Account{%d, %s}\n", accountID, name)

		account := Account{
			ID:   accountID,
			Name: name,
		}

		res := db.Raw.Create(&account)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

// Imports trading records from a CSV file.
func ImportRecords(db *DB, filePath string) error {
	ch := make(chan []string)
	go ReadCSV(filePath, ch)
	for row := range ch {
		// "Test Account","AMZN",DEPOSIT,2017-02-24,1
		accountName := strings.TrimSpace(row[0])
		assetName := strings.TrimSpace(row[1])
		recordType := RecordType(strings.TrimSpace(row[2]))
		createdAt, _ := time.Parse("2006-01-02", strings.TrimSpace(row[3]))
		quantity, _ := strconv.ParseInt(strings.TrimSpace(row[4]), 10, 64)

		account := db.GetAccountByName(accountName)
		asset := db.GetAssetByName(assetName)

		record := Record{
			Account:   account,
			Asset:     asset,
			Type:      recordType,
			CreatedAt: createdAt,
			Quantity:  int(quantity),
		}
		res := db.Raw.Create(&record)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func ImportStockValues(filePath string, symbol string) error {
	db := ConnectDatabase()
	defer db.Raw.Close()

	var asset Asset
	db.Raw.First(&asset, "name = ?", symbol)
	if asset == (Asset{}) {
		db.InsertAsset(symbol, "")
	}

	ch := make(chan AssetValue)
	go ReadStockValues(filePath, ch)
	for v := range ch {
		res := db.Raw.Create(&v)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}
