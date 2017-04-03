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

func ReadStockValues(filePath string) {
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

		av := AssetValue{
			EvaluatedAt: evaluatedAt,
			Open:        open,
			High:        high,
			Low:         low,
			Close:       close,
			Volume:      volume,
		}

		fmt.Println(av)
	}
}
