package errors

import (
	"fmt"

	"github.com/txgruppi/safe/regexp"
)

var (
	ErrInvalidFileLocation = fmt.Errorf("file location is invalid, it must match %q", regexp.FileLocation)
	ErrInvalidSize         = fmt.Errorf("size must be greater or equal to 0")
)
