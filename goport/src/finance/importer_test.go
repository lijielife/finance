package finance

import (
	"fmt"
	"testing"
)

func TestReadCSV(t *testing.T) {
	// TODO: Look for something like @pytest.mark.parametrize
	files := []string{"accounts", "assets"}

	for _, filename := range files {
		ch := make(chan []string)
		go ReadCSV(fmt.Sprintf("test-data/%s.csv", filename), ch)
		for v := range ch {
			fmt.Println("Processing", v)
		}
	}
}

func DeleteExistingAccounts(db *DB, filePath string) {
	ch := make(chan []string)
	go ReadCSV(filePath, ch)
	for row := range ch {
		db.Raw.Delete(Account{}, "id = ?", row[0])
	}
}

func TestImportAccounts(t *testing.T) {
	db := ConnectDatabase()
	defer db.Raw.Close()

	filePath := "test-data/accounts.csv"
	DeleteExistingAccounts(db, filePath)

	err := ImportAccounts(filePath)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan []string)
	go ReadCSV(filePath, ch)
	for row := range ch {
		name := row[1]
		var account Account
		db.Raw.First(&account, "name = ?", name)
		if account == (Account{}) {
			t.Errorf("Account '%s' not found", name)
		}
	}
}

func TestImportRecords(t *testing.T) {
	db := ConnectDatabase()
	defer db.Raw.Close()

	filePath := "test-data/records.csv"

	err := ImportRecords(db, filePath)
	if err != nil {
		t.Error(err)
	}
}
