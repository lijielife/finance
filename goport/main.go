package main

import (
	"finance"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"os"
)

type Env struct {
	db finance.Datastore
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "schema",
			Usage: "Setup database schema",
			Action: func(c *cli.Context) error {
				db := finance.ConnectDatabase()
				defer db.Raw.Close()

				env := &Env{db}
				env.db.CreateTables()
				return nil
			},
		},
		{
			Name:    "import_stock_values",
			Aliases: []string{"s"},
			Usage:   "Import stock values by reading a CSV file",
			Action: func(c *cli.Context) error {
				args := c.Args()
				filePath, symbol := args.Get(0), args.Get(1)

				db := finance.ConnectDatabase()
				defer db.Raw.Close()

				finance.ImportStockValues(db, filePath, symbol)
				return nil
			},
		},
		{
			Name:  "webserver",
			Usage: "Run a web server",
			Action: func(c *cli.Context) error {
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "Hello World!")
				})
				http.ListenAndServe(":8080", nil)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
