package finance

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func ReadStockValues() {
	f, err := os.Open("./test.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(records)
}
