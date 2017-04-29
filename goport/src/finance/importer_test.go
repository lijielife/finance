package finance

import (
	"fmt"
	"testing"
)

func TestImportStockValues(t *testing.T) {
	ImportStockValues("test-data/MSFT.csv", "MSFT")
}

func TestReadCSV(t *testing.T) {
	files := []string{"accounts", "assets"}

	for _, filename := range files {
		ch := make(chan []string)
		go ReadCSV(fmt.Sprintf("test-data/%s.csv", filename), ch)
		for v := range ch {
			fmt.Println("Processing", v)
		}
	}
}
