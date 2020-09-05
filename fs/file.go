package fs

import (
	"mime"

	"github.com/txgruppi/safe/errors"
	"github.com/txgruppi/safe/regexp"
)

type File interface {
	Location() string
	SetLocation(string) File
	MimeType() string
	SetMimeType(string) File
	Size() int64
	SetSize(int64) File
	Data() []byte
	SetData([]byte) File
	Validate() error
}

func NewEmptyFile() File {
	return &file{}
}

type file struct {
	location string
	mimeType string
	size     int64
	data     []byte
}

func (t *file) Location() string {
	return t.location
}

func (t *file) SetLocation(location string) File {
	t.location = location
	return t
}

func (t *file) MimeType() string {
	return t.mimeType
}

func (t *file) SetMimeType(mimeType string) File {
	t.mimeType = mimeType
	return t
}

func (t *file) Size() int64 {
	return t.size
}

func (t *file) SetSize(size int64) File {
	t.size = size
	return t
}

func (t *file) Data() []byte {
	if t.data == nil {
		return nil
	}
	d := make([]byte, len(t.data))
	for i := range t.data {
		d[i] = t.data[i]
	}
	return d
}

func (t *file) SetData(data []byte) File {
	if data == nil {
		return nil
	}
	d := make([]byte, len(data))
	for i := range data {
		d[i] = data[i]
	}
	t.data = data
	return t
}

func (t *file) Validate() error {
	if !regexp.FileLocation.MatchString(t.location) {
		return errors.ErrInvalidFileLocation
	}
	if _, _, err := mime.ParseMediaType(t.mimeType); err != nil {
		return err
	}
	if t.size < 0 {
		return errors.ErrInvalidSize
	}
	return nil
}
