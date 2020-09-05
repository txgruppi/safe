package errors

import "fmt"

var (
	ErrDBAlreadyUnlocked  = fmt.Errorf("db already unlocked")
	ErrDBNotUnlocked      = fmt.Errorf("db not unlocked")
	ErrCantDecodeFileSize = fmt.Errorf("can't decode file size")
)
