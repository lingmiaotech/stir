package main

import (
	"database/sql"
	"errors"
	"github.com/lingmiaotech/tonic"
	"github.com/pressly/goose"
	"github.com/urfave/cli"
	"os"
)

var DBCommand = cli.Command{
	Name:  "db",
	Usage: "database related operations",
	Subcommands: []cli.Command{
		DBCreateCommand,
		DBUpCommand,
		DBUpToCommand,
		DBDownCommand,
		DBDownToCommand,
	},
}

var DBCreateCommand = cli.Command{
	Name:  "create",
	Usage: "create migration",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "env, e", Value: "./configs/development.yaml", Usage: "configs path", EnvVar: "APP_ENV"},
		cli.StringFlag{Name: "message, m", Value: "", Usage: "migration message"},
	},
	Action: func(c *cli.Context) (err error) {
		message := c.String("message")
		if message == "" {
			return MakeExitError(errors.New("missing_argument_message"), "preparing operation")
		}

		env := c.String("env")
		os.Setenv("APP_ENV", env)

		err = database("create", "./migrations", "message", "sql")
		if err != nil {
			return err
		}

		return nil
	},
}

var DBUpCommand = cli.Command{
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

var DBUpToCommand = cli.Command{
	Name:  "up-to",
	Usage: "up to specific version",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "env, e", Value: "./configs/development.yaml", Usage: "configs path", EnvVar: "APP_ENV"},
		cli.StringFlag{Name: "version, v", Value: "", Usage: "version of target migration"},
	},
	Action: func(c *cli.Context) (err error) {
		version := c.String("version")
		if version == "" {
			return MakeExitError(errors.New("missing_argument_version"), "preparing operation")
		}

		env := c.String("env")
		os.Setenv("APP_ENV", env)

		err = database("up-to", "./migrations", version)
		if err != nil {
			return err
		}

		return nil
	},
}

var DBDownCommand = cli.Command{
	Name:  "down",
	Usage: "down to earliest version",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "env, e", Value: "./configs/development.yaml", Usage: "configs path", EnvVar: "APP_ENV"},
	},
	Action: func(c *cli.Context) (err error) {
		env := c.String("env")
		os.Setenv("APP_ENV", env)

		err = database("down", "./migrations")
		if err != nil {
			return err
		}

		return nil
	},
}

var DBDownToCommand = cli.Command{
	Name:  "down-to",
	Usage: "down to specific version",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "env, e", Value: "./configs/development.yaml", Usage: "configs path", EnvVar: "APP_ENV"},
		cli.StringFlag{Name: "version, v", Value: "", Usage: "version of target migration"},
	},
	Action: func(c *cli.Context) (err error) {
		version := c.String("version")
		if version == "" {
			return MakeExitError(errors.New("missing_argument_version"), "preparing operation")
		}

		env := c.String("env")
		os.Setenv("APP_ENV", env)

		err = database("down-to", "./migrations", version)
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
		return MakeExitError(err, "initializing configs file")
	}

	driver := tonic.Configs.GetString("database.driver")
	if driver == "" {
		return MakeExitError(errors.New("empty_driver_configs"), "initializing configs file")
	}

	err = goose.SetDialect(driver)
	if err != nil {
		return MakeExitError(err, "applying dialect driver")
	}

	dbstring := tonic.Configs.GetString("database.dbstring")
	if dbstring == "" {
		return MakeExitError(errors.New("empty_dbstring_config"), "initializing configs file")
	}

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		return MakeExitError(err, "connecting database")
	}

	err = goose.Run(command, db, dir, args...)
	if err != nil {
		return MakeExitError(err, "executing commands")
	}

	return nil
}
