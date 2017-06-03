package utils

import (
	"fmt"
	"github.com/urfave/cli"
)

func MakeExitError(err error, doing string) *cli.ExitError {
	return cli.NewExitError(fmt.Sprintf("Got error (%s) while %s.", err.Error(), doing), 1)
}
