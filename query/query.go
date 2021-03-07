package query

import (
	"context"
	"tgenj/database"
	"tgenj/document"
	"tgenj/expr"
)

/**
	A Query can execute statements against the database.
	It can read or write data from any table, or even alter the structure of the database.
	Results are returned as streams.
 */
/*TODO*/
type Query struct {
	//Statements []Statement
	tx         *database.Transaction
	autoCommit bool
}

/*TODO*/
func (q Query) Run(ctx context.Context,db *database.Database) (*Result,error) {
	return nil,nil
}

/*Result of a query*/
type Result struct {
	Iterator document.Iterator
	Tx       *database.Transaction
	closed   bool
}


func (r *Result) Iterate(fn func(d document.Document) error) error {
	if r.Iterator == nil {
		return nil
	}
	return r.Iterator.Iterate(fn)
}

/* A Statement represents a unique action that can be executed against the database*/
type Statement interface {
	Run(*database.Transaction, []expr.Param) (Result, error)
	IsReadOnly() bool
}
