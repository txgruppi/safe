package database

import (
	"sync"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/txgruppi/safe/errors"
	"github.com/txgruppi/safe/fs"
)

type IteratorCloseFunc func() error

type Database interface {
	Unlock(password []byte) error
	Lock() error
	Tidy() error
	Set(fs.File) error
	Get(location string) (fs.File, error)
	Del(location string) error
	Iterator(prefix string, loadFileData bool) (Iterator, error)
}

func New(databasePath string, verbose bool) (Database, error) {
	return &database{
		path:    databasePath,
		verbose: verbose,
	}, nil
}

type database struct {
	path    string
	db      *badger.DB
	verbose bool
	lock    sync.Mutex
}

func (t *database) Unlock(password []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.db != nil {
		return errors.ErrDBAlreadyUnlocked
	}
	opts := badger.DefaultOptions(t.path).WithEncryptionKey(password)
	if !t.verbose {
		opts.Logger = nil
	}
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	t.db = db
	return nil
}

func (t *database) Lock() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.db != nil {
		if err := t.db.Close(); err != nil {
			return err
		}
		t.db = nil
	}
	return nil
}

func (t *database) ensureDBIsUnlocked() error {
	if t.db == nil {
		return errors.ErrDBNotUnlocked
	}
	return nil
}

func (t *database) Tidy() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err := t.ensureDBIsUnlocked(); err != nil {
		return err
	}
	for {
		if err := t.db.RunValueLogGC(0.5); err != nil {
			if err == badger.ErrNoRewrite {
				return nil
			}
			return err
		}
	}
}

func (t *database) Del(location string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err := t.ensureDBIsUnlocked(); err != nil {
		return err
	}

	return t.db.Update(func(tx *badger.Txn) error {
		id, err := locationToID(tx, location)
		if err != nil {
			return err
		}
		if err := deleteFile(tx, id); err != nil {
			return err
		}
		return nil
	})
}

func (t *database) Set(file fs.File) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err := t.ensureDBIsUnlocked(); err != nil {
		return err
	}

	if err := file.Validate(); err != nil {
		return err
	}
	id, err := generateID()
	if err != nil {
		return err
	}
	return t.db.Update(func(tx *badger.Txn) error {
		locationKey := []byte(file.Location())
		if err := writeFile(tx, id, file); err != nil {
			return err
		}
		if err := tx.Set(locationKey, []byte(id)); err != nil {
			return err
		}
		return nil
	})
}

func (t *database) Get(location string) (file fs.File, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err := t.ensureDBIsUnlocked(); err != nil {
		return nil, err
	}

	err = t.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte(location))
		if err != nil {
			return err
		}
		idBytes, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		f, err := readFile(tx, string(idBytes), true)
		if err != nil {
			return err
		}
		f.SetLocation(location)
		if err := f.Validate(); err != nil {
			return err
		}
		file = f
		return nil
	})
	return
}

func (t *database) Iterator(prefix string, loadFileData bool) (Iterator, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err := t.ensureDBIsUnlocked(); err != nil {
		return nil, err
	}

	return &iterator{
		db:           t.db,
		prefix:       prefix,
		loadFileData: loadFileData,
	}, nil
}
