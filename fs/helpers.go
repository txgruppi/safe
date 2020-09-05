package fs

import (
	"mime"
	"path"
)

func SafeMimeType(filepath string) (mt string) {
	defer recover()
	mt = "application/octet-stream"
	mt = mime.TypeByExtension(path.Ext(filepath))
	return
}
