package main

import (
	"github.com/lingmiaotech/stir/cmd/db"
	"github.com/urfave/cli"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "Stir"
	app.HelpName = "stir"
	app.Version = "1.0.1"
	app.Usage = "Stir is a command line tool for tonic framework"
	app.Commands = []cli.Command{
		db.Command,
	}

	app.Run(os.Args)
}
