package query

import (
	"errors"
	"tgenj/database"
	"tgenj/expr"
)

/**
 *  @ClassName:alter.go
 *  @Description:TODO
 *  @Author:jackey
 *  @Create:2021/3/11 下午5:20
 */

/*AlterStmt is a DSL that allows creating a full ALTER TABLE query.*/
type AlterStmt struct {
	TableName string
	NewTableName string
}

/* IsReadOnly always returns false. It implements the Statement interface.*/
func (stmt AlterStmt) isReadOnly() bool {
	return false
}

func (stmt AlterStmt) Run(tx *database.Transaction, _ []expr.Param) (Result,error) {
	var res Result

	if stmt.TableName=="" {
		return res,errors.New("missing table name")
	}

	if stmt.NewTableName == "" {
		return res,errors.New("missing new table name")
	}

	if stmt.NewTableName == stmt.TableName {
		return res,database.ErrIndexAlreadyExists
	}

	err := tx.
}
