package db

import (
	"database/sql"
	"errors"
	"github.com/lingmiaotech/stir/utils"
	"github.com/lingmiaotech/tonic"
	"github.com/pressly/goose"
	"github.com/urfave/cli"
	"os"
)

var Command = cli.Command{
	Name:  "db",
	Usage: "database related operations",
	Subcommands: []cli.Command{
		UpCommand,
	},
}

var UpCommand = cli.Command{
	Name:  "up",
	Usage: "up to latest version",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "env, e", Value: "./configs/development.yaml", Usage: "configs path", EnvVar: "APP_ENV"},
	},
	Action: func(c *cli.Context) (err error) {
		env := c.String("env")
		os.Setenv("APP_ENV", env)

		err = database("up", "./migrations")
		if err != nil {
			return err
		}

		return nil
	},
}

func database(command string, dir string, args ...string) error {
	var err error

	err = tonic.InitConfigs()
	if err != nil {
		return utils.MakeExitError(err, "initializing configs file")
	}

	driver := tonic.Configs.GetString("database.driver")
	if driver == "" {
		return utils.MakeExitError(errors.New("empty_driver_configs"), "initializing configs file")
	}

	err = goose.SetDialect(driver)
	if err != nil {
		return utils.MakeExitError(err, "applying dialect driver")
	}

	dbstring := tonic.Configs.GetString("database.dbstring")
	if dbstring == "" {
		return utils.MakeExitError(errors.New("empty_dbstring_config"), "initializing configs file")
	}

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		return utils.MakeExitError(err, "connecting database")
	}

	err = goose.Run(command, db, dir, args...)
	if err != nil {
		return utils.MakeExitError(err, "executing commands")
	}

	return nil
}
