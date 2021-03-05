package tgenj

import (
	"context"
	"errors"
	"tgenj/database"
	"tgenj/document"
	"tgenj/query"
	"tgenj/result"
	"tgenj/sql/parser"
)

/*DB represents a collection of tables stored in the underlying engine.*/
type DB struct {
	DB *database.Database

	ctx context.Context
}

/*添加上下文信息*/
func (db *DB) WithContext(ctx context.Context) (*DB, error) {
	if ctx == nil {
		return nil, errors.New("创建失败")
	}
	return &DB{
		DB:  db.DB,
		ctx: ctx,
	}, nil
}

/*close database*/
func (db *DB) Close() error {
	return db.DB.Close()
}

/*Begin starts a new transaction*/
func (db *DB) Begin(writable bool) (*Tx, error) {

	tx, err := db.DB.BeginTx(db.ctx, &database.TxOptions{
		ReadOnly: !writable,
	})

	if err != nil {
		return nil, err
	}

	return &Tx{
		Transaction: tx,
	}, nil
}

/*View starts a read only transaction, runs fn and automatically rolls it back.*/
func (db *DB) View(fn func(tx *Tx) error) error {
	tx, err := db.Begin(false)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	return fn(tx)
}

/*Update starts a read-write transaction, runs fn and automatically commits it*/
func (db *DB) Update(fn func(tx *Tx) error) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

/* Exec a query against the database without returning the result.*/
/*TODO*/
func (db *DB) Exec(q string, args ...interface{}) error {
	res, err := db.Query(q, args...)
	if err != nil {
		return err
	}

	/*TODO defer res.close*/

	return res.Iterate(func(d document.Document) error {
		return nil
	})

}

/* Query the database and return the result.*/
/*TODO*/
func (db *DB) Query(q string, args ...interface{}) (*query.Result,error) {
	pd,err := parser.ParseQuery(q)
	if err != nil {
		return nil, err
	}
	return pd.Run(db.ctx,db.DB)
}

func (db *DB) QueryDocument() {

}

/*Tx represents a database transaction.*/
/*Read-only or read/write*/
type Tx struct {
	*database.Transaction
}

/*Tx 包含几个操作*/
/*TODO*/
func (tx *Tx) Query(str string, args ...interface{}) (*result.Result, error) {
	return nil, nil
}
