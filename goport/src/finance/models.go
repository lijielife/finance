package finance

import (
	"database/sql/driver"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"time"
)

type Datastore interface {
	CreateTables()
	GetAssetByName(name string) Asset
	InsertAsset(name string, description string) (Asset, []error)
}

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
	ID   uint64 `gorm:"primary_key"`
	Name string `sql:"type:varchar(255);"`
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

// Record represents a single financial trade
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

// ConnectDatabase connects to a database and returns a wrapper object
// containing an instance of `gorm.DB`.
func ConnectDatabase() *DB {
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

	return &DB{db}
}

// CreateTables creates all necessary tables.
func (db *DB) CreateTables() {
	// Any better way to handle this?
	db.Raw.Exec("DROP TYPE IF EXISTS granularity CASCADE")
	db.Raw.Exec("CREATE TYPE granularity AS ENUM('1sec', '1min', '5min', '1hour', '1day', '1week', '1month', '1year')")

	db.Raw.Exec("DROP TYPE IF EXISTS record_type CASCADE")
	db.Raw.Exec("CREATE TYPE record_type AS ENUM('deposit', 'withdraw', 'balance_adjustment')")

	// Migrate the schema
	db.Raw.AutoMigrate(&Account{})
	db.Raw.AutoMigrate(&Asset{})
	db.Raw.AutoMigrate(&AssetValue{})
	db.Raw.AutoMigrate(&Record{})

	// // Update - update product's price to 2000
	// db.Raw.Model(&product).Update("Price", 2000)

	// // Delete - delete product
	// db.Raw.Delete(&product)
}

///////////////////////////////////////////////////////////////////////////////

// A wrapper for `gorm.DB` object.
type DB struct {
	Raw *gorm.DB
}

// Returns an `Asset` instance matching the given name.
func (db *DB) GetAssetByName(name string) Asset {
	var asset Asset
	// FIXME: Deal with cases where no matching record found
	db.Raw.First(&asset, "name = ?", name)
	return asset
}

func (db *DB) GetAccountByName(name string) Account {
	var account Account
	// FIXME: Deal with cases where no matching record found
	db.Raw.First(&account, "name = ?", name)
	return account
}

func (db *DB) InsertAsset(name string, description string) (Asset, []error) {
	asset := Asset{
		Name:        name,
		Description: description,
	}
	res := db.Raw.Create(&asset)
	return asset, res.GetErrors()
}

func (db *DB) InsertRecord(account Account, asset Asset,
	recordType RecordType, createdAt time.Time, quantity int) (Record, []error) {

	record := Record{
		Account:   account,
		Asset:     asset,
		Type:      recordType,
		CreatedAt: createdAt,
		Quantity:  quantity,
	}
	res := db.Raw.Create(&record)
	return record, res.GetErrors()
}

func (db *DB) InsertAccount(name string) (Account, error) {
	account := Account{
		Name: name,
	}
	res := db.Raw.Create(&account)
	return account, res.Error
}
