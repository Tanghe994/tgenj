package engine

import (
	"context"
	"tgenj/options"
	"tgenj/transaction"
)

/*负责存储数据*/
type Engine interface {
	/*returns a read-only or read/write transaction depending on whether writable is set to false or true, respectively.*/
	Begin (ctx context.Context, opts options.TxOptions)(transaction.Transaction,error)

	/*close engine*/
	Close() error
}

// A Store manages key value pairs. It is an abstraction on top of any data structure that can provide
// random read, random write, and ordered sequential read.
type Store interface {
	// Get returns a value associated with the given key. If no key is not found, it returns ErrKeyNotFound.
	Get(k []byte) ([]byte, error)
	// Put stores a key value pair. If it already exists, it overrides it.
	// Both k and v must be not nil.
	Put(k, v []byte) error
	// Delete a key value pair. If the key is not found, returns ErrKeyNotFound.
	Delete(k []byte) error
	// Truncate deletes all the key value pairs from the store.
	Truncate() error
	// Iterator creates an iterator with the given options.
	// The initial position depends on the implementation.
	// TODO Iterator(opts IteratorOptions) Iterator
	// NextSequence returns a monotonically increasing integer.
	NextSequence() (uint64, error)
}
