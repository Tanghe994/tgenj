package database

import (
	"context"
	"errors"
	"sync"
	"tgenj/engine"
	"tgenj/options"
)

type Database struct {
	/*明确底层的存储引擎*/
	engine              engine.Engine

	/*主要用于检测当前是否有其它的事务在进行*/
	attachedTransaction *Transaction
	attachedTxMu        sync.Mutex
	/*This controls concurrency on read-only and read/write transactions*/
	txmu 				sync.RWMutex

}

/*Close the underlying engine*/
func (db *Database) Close() error {
	return db.engine.Close()
}

/* Begin starts a new transaction with default options.*/
func (db Database) Begin(writable bool) (*Transaction, error) {
	txo := &TxOptions{
		ReadOnly: !writable,
	}
	return db.BeginTx(context.Background(), txo)
}

func (db *Database) BeginTx(ctx context.Context, opts *TxOptions) (*Transaction, error) {
	if opts == nil {
		opts = new(TxOptions)
	}
	/*ReadOnly 是可以写入的*/
	/*这里需要添加互斥锁*/
	if !opts.ReadOnly {
		/*可以读取则上锁进行写入*/
		db.txmu.Lock()
	} else {
		/*RLock方法将rw锁定为读取状态，禁止其他线程写入，但不禁止读取。*/
		db.txmu.RLock()
	}

	db.attachedTxMu.Lock()
	defer db.attachedTxMu.Unlock()

	if db.attachedTransaction != nil {
		return nil,errors.New("cannot open a transaction within a transaction")
	}

	ntx, err := db.engine.Begin(ctx, options.TxOptions{
		Writable: !opts.ReadOnly,
	})

	if err != nil {
		return nil,err
	}

	tx := Transaction{
		db:       db,
		tx:       ntx,
		writable: !opts.ReadOnly,
		attached: opts.Attached,
	}
	if opts.Attached {
		db.attachedTransaction=&tx
	}

	return &tx, nil
}


// TxOptions are passed to Begin to configure transactions.
type TxOptions struct {
	// Open a read-only transaction.
	ReadOnly bool
	// Set the transaction as global at the database level.
	// Any queries run by the database will use that transaction until it is
	// rolled back or commited.
	Attached bool
}

func (db Database) GetAttachedTx() *Transaction {

	/*要上锁进行读取*/
	db.attachedTxMu.Lock()
	defer db.attachedTxMu.Unlock()

	return db.attachedTransaction
}

/*TODO*/
func (db Database) GetCatalog() error {
	return nil
}


