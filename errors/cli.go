package errors

import "fmt"

var (
	ErrMissingFileArgument    = fmt.Errorf("missing file argument")
	ErrWrongNumberOfArguments = fmt.Errorf("wrong number of arguments, expected list of location/file pair")
)
