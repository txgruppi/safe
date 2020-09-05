package database

import (
	"sync"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/txgruppi/safe/fs"
)

type Iterator interface {
	Channel() <-chan fs.File
	Close()
	Error() error
}

type iterator struct {
	db           *badger.DB
	prefix       string
	started      bool
	ended        bool
	loadFileData bool
	err          error
	itemCh       chan fs.File
	closeCh      chan struct{}
	lock         sync.Mutex
}

func (t *iterator) start() {
	if t.started {
		return
	}
	t.started = true

	t.itemCh = make(chan fs.File)
	t.closeCh = make(chan struct{})

	go func() {
		prefix := []byte(t.prefix)
		t.err = t.db.View(func(tx *badger.Txn) error {
			it := tx.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			idBytes := make([]byte, 0)
			var file fs.File
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()
				idBytes, t.err = item.ValueCopy(idBytes)
				if t.err != nil {
					break
				}
				file, t.err = readFile(tx, string(idBytes), t.loadFileData)
				if t.err != nil {
					break
				}
				file.SetLocation(string(item.Key()))
				t.err = file.Validate()
				if t.err != nil {
					break
				}
				select {
				case t.itemCh <- file:
				case <-t.closeCh:
					break
				}
			}

			t.lock.Lock()
			defer t.lock.Unlock()

			if t.itemCh != nil {
				close(t.itemCh)
				t.itemCh = nil
			}
			if t.closeCh != nil {
				close(t.closeCh)
				t.closeCh = nil
			}
			t.ended = true

			return nil
		})
	}()
}

func (t *iterator) Close() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.closeCh != nil {
		close(t.closeCh)
		t.closeCh = nil
	}
}

func (t *iterator) Channel() <-chan fs.File {
	t.lock.Lock()
	defer t.lock.Unlock()
	if !t.started {
		t.start()
	}
	return t.itemCh
}

func (t *iterator) Error() error {
	return t.err
}
