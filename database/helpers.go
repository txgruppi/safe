package database

import (
	"encoding/binary"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/txgruppi/safe/errors"
	"github.com/txgruppi/safe/fs"
)

func itemValueCopy(tx *badger.Txn, key []byte) ([]byte, error) {
	item, err := tx.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func locationToID(tx *badger.Txn, location string) (string, error) {
	item, err := tx.Get([]byte(location))
	if err != nil {
		return "", err
	}
	idBytes, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}
	return string(idBytes), nil
}

func readFile(tx *badger.Txn, id string, loadFileData bool) (fs.File, error) {
	keys := keysForID(id)

	mimeTypeBytes, err := itemValueCopy(tx, keys.mimeType)
	if err != nil {
		return nil, err
	}
	sizeBytes, err := itemValueCopy(tx, keys.size)
	if err != nil {
		return nil, err
	}
	size, n := binary.Varint(sizeBytes)
	if n <= 0 {
		return nil, errors.ErrCantDecodeFileSize
	}

	file := fs.
		NewEmptyFile().
		SetMimeType(string(mimeTypeBytes)).
		SetSize(size)

	if loadFileData {
		dataBytes, err := itemValueCopy(tx, keys.data)
		if err != nil {
			return nil, err
		}
		file.SetData(dataBytes)
	}

	return file, nil
}

func writeFile(tx *badger.Txn, id string, file fs.File) error {
	if err := file.Validate(); err != nil {
		return err
	}
	keys := keysForID(id)

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, file.Size())

	if err := tx.Set(keys.mimeType, []byte(file.MimeType())); err != nil {
		return err
	}
	if err := tx.Set(keys.size, buf[:n]); err != nil {
		return err
	}
	if err := tx.Set(keys.data, file.Data()); err != nil {
		return err
	}

	return nil
}

func deleteFile(tx *badger.Txn, id string) error {
	keys := keysForID(id)
	if err := tx.Delete(keys.mimeType); err != nil {
		return err
	}
	if err := tx.Delete(keys.size); err != nil {
		return err
	}
	if err := tx.Delete(keys.data); err != nil {
		return err
	}
	return nil
}
