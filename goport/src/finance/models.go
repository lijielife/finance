package finance

import (
	"database/sql/driver"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"time"
)

//
// AccountType
//

type AccountType string

const (
	CHECKING    AccountType = "checking"
	SAVINGS     AccountType = "savings"
	INVESTMENT  AccountType = "investment"
	CREDIT_CARD AccountType = "credit_card"
	VIRTUAL     AccountType = "virtual"
)

func (u *AccountType) Scan(value interface{}) error {
	*u = AccountType(value.(string))
	return nil
}
func (u AccountType) Value() (driver.Value, error) { return string(u), nil }

//
// Granularity
//

type Granularity string

const (
	SEC      Granularity = "1sec"
	MIN      Granularity = "1min"
	FIVE_MIN Granularity = "5min"
	HOUR     Granularity = "1hour"
	DAY      Granularity = "1day"
	WEEK     Granularity = "1week"
	MONTH    Granularity = "1month"
	YEAR     Granularity = "1year"
)

func (u *Granularity) Scan(value interface{}) error {
	*u = Granularity(value.(string))
	return nil
}
func (u Granularity) Value() (driver.Value, error) { return string(u), nil }

//
// RecordType
//

type RecordType string

const (
	DEPOSIT            RecordType = "deposit"
	WITHDRAW           RecordType = "withdraw"
	BALANCE_ADJUSTMENT RecordType = "balance_adjustment"
)

func (u *RecordType) Scan(value interface{}) error {
	*u = RecordType(value.(string))
	return nil
}
func (u RecordType) Value() (driver.Value, error) { return string(u), nil }

///////////////////////////////////////////////////////////////////////////////

type Account struct {
	Name string
}

type Asset struct {
	ID          uint64 `gorm:"primary_key"`
	Name        string `sql:"type:varchar(255);"`
	Description string
}

type AssetValue struct {
	ID          uint64 `gorm:"primary_key"`
	Asset       Asset
	AssetID     uint64
	BaseAsset   Asset
	BaseAssetID uint64
	EvaluatedAt time.Time   `sql:"DEFAULT:current_timestamp"`
	Granularity Granularity `sql:"not null;type:GRANULARITY"`
	Open        float64     `sql:"type:decimal(10,4);"`
	High        float64     `sql:"type:decimal(10,4);"`
	Low         float64     `sql:"type:decimal(10,4);"`
	Close       float64     `sql:"type:decimal(10,4);"`
	Volume      int64
}

type Record struct {
	ID        uint64 `gorm:"primary_key"`
	Account   Account
	AccountID uint64
	Asset     Asset
	AssetID   uint64
	Type      RecordType `sql:"not null;type:record_type"`
	CreatedAt time.Time  `sql:"DEFAULT:current_timestamp"`
	Quantity  int
}

func ConnectDatabase() *gorm.DB {
	dbUrl, found := os.LookupEnv("DB_URL")
	if !found {
		panic("Could not find an environment variable DB_URL")
	}

	fmt.Printf("Connecting to %s...\n", dbUrl)
	db, err := gorm.Open("postgres", dbUrl+"?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	return db
}

func CreateTables(db *gorm.DB) {
	// Any better way to handle this?
	db.Exec("DROP TYPE IF EXISTS granularity CASCADE")
	db.Exec("CREATE TYPE granularity AS ENUM('1sec', '1min', '5min', '1hour', '1day', '1week', '1month', '1year')")

	// Migrate the schema
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Asset{})
	db.AutoMigrate(&AssetValue{})
	db.AutoMigrate(&Record{})

	// // Update - update product's price to 2000
	// db.Model(&product).Update("Price", 2000)

	// // Delete - delete product
	// db.Delete(&product)
}

///////////////////////////////////////////////////////////////////////////////

func GetAssetByName(db *gorm.DB, name string) Asset {
	var asset Asset
	db.First(&asset, "name = ?", name)
	return asset
}

func InsertAsset(db *gorm.DB, name string, description string) (Asset, []error) {
	asset := Asset{
		Name:        name,
		Description: description,
	}
	res := db.Create(&asset)
	return asset, res.GetErrors()
}

func InsertRecord(db *gorm.DB, account Account, asset Asset,
	recordType RecordType, createdAt time.Time, quantity int) (Record, []error) {

	record := Record{
		Account:   account,
		Asset:     asset,
		Type:      recordType,
		CreatedAt: createdAt,
		Quantity:  quantity,
	}
	res := db.Create(&record)
	return record, res.GetErrors()
}
