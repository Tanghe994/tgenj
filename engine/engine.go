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
