package finance

import (
	"database/sql/driver"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"time"
)

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

type Asset struct {
	ID   uint64 `gorm:"primary_key"`
	Name string `sql:"type:varchar(255);"`
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

func ConnectDatabase() *gorm.DB {
	dbUrl, _ := os.LookupEnv("DB_URL")

	fmt.Printf("Connecting to %s...\n", dbUrl)
	db, err := gorm.Open("postgres", dbUrl+"?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	return db
}

func CreateTables() {
	db := ConnectDatabase()

	// Any better way to handle this?
	db.Exec("DROP TYPE IF EXISTS granularity CASCADE")
	db.Exec("CREATE TYPE granularity AS ENUM('1sec', '1min', '5min', '1hour', '1day', '1week', '1month', '1year')")

	// Migrate the schema
	db.AutoMigrate(&Asset{})
	db.AutoMigrate(&AssetValue{})

	// Create
	//db.Create(&Asset{Name: "Test"})

	// Read
	// var product Product
	// db.First(&product, 1) // find product with id 1
	// fmt.Printf("product = %+v\n", product)
	// db.First(&product, "code = ?", "L1212") // find product with code l1212
	// fmt.Printf("product = %+v\n", product)

	// // Update - update product's price to 2000
	// db.Model(&product).Update("Price", 2000)

	// // Delete - delete product
	// db.Delete(&product)
}
