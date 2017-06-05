package main

import (
	"errors"
	"testing"
)

func TestMakeExitError(t *testing.T) {

	err := MakeExitError(errors.New("some error"), "doing something")

	expect(t, err.Error(), "Got error (some error) while doing something.")

}
