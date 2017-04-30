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
	db := finance.ConnectDatabase()
	env := &Env{db}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "schema",
			Usage: "Setup database schema",
			Action: func(c *cli.Context) error {
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

				finance.ImportStockValues(filePath, symbol)
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
