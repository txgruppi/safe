package regexp

import "regexp"

var (
	FileLocation = regexp.MustCompile("^(/[a-z0-9_.-]+)+$")
)
