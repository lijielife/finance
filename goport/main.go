package main

import (
	"finance"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "schema",
			Usage: "Setup database schema",
			Action: func(c *cli.Context) error {
				finance.CreateTables()
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
	}
	app.Run(os.Args)
}
