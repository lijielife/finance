package finance

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateTables(t *testing.T) {
	db := ConnectDatabase()
	db.Raw.DropTable(&Account{})
	db.Raw.DropTable(&Asset{})
	db.Raw.DropTable(&AssetValue{})
	db.Raw.DropTable(&Record{})
	db.CreateTables()
}

func TestInsertAsset(t *testing.T) {
	db := ConnectDatabase()
	asset, errors := db.InsertAsset("KRW", "Korean Won")

	if len(errors) > 0 {
		t.Error(errors)
	}
	if asset.Name != "KRW" {
		t.Error("Asset was not properly inserted")
	}
}

func TestGetAssetByName(t *testing.T) {
	db := ConnectDatabase()
	asset, err := db.GetAssetByName("KRW")
	if asset == (Asset{}) {
		t.Error("No such asset found")
	}
	if err != nil {
		t.Error(err)
	}
	fmt.Println(asset)
}

func TestInsertRecord(t *testing.T) {
	db := ConnectDatabase()
	account, err := db.InsertAccount("Default Account")
	if err != nil {
		t.Error(err)
	}
	asset, _ := db.InsertAsset("NVDA", "")
	db.InsertRecord(account, asset, DEPOSIT, time.Now(), 10)
}
