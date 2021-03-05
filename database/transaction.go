package database

import "tgenj/transaction"

/*Transaction represents a database transaction*/
type Transaction struct {
	/*该事务所要操作的数据库*/
	db *Database

	/*该事务需要操作的底层引擎*/
	/*源码是engine.Transaction*/
	/*这是一个接口*/
	tx transaction.Transaction

	/*只分两种，只读或读写*/
	writable bool

	// if set to true, this transaction is attached to the database
	attached bool

	/*rollback and commit function*/
	onRollbackHooks []func()
	onCommitHooks   []func()
}

func (tx *Transaction) GetDB() *Database {
	return tx.db
}

func (tx *Transaction) Rollback() error {
	err := tx.tx.Rollback()
	if err != nil {
		return err
	}
	defer tx.Unlock()

	/*如果该事务已经附加到某个数据库上面*/
	if tx.attached {
		/*对连接事务加锁*/
		tx.db.attachedTxMu.Lock()
		/*结束后解锁*/
		defer tx.db.attachedTxMu.Unlock()
		/*清空连接的事务*/
		if tx.db.attachedTransaction != nil {
			tx.db.attachedTransaction = nil
		}
	}

	for i := len(tx.onRollbackHooks) - 1 ;i>=0;i--{
		tx.onRollbackHooks[i]()
	}
	return err
}

func (tx Transaction) Commit() error {
	err :=tx.tx.Commit()
	if err != nil {
		return err
	}
	defer tx.Unlock()

	if tx.attached {
		tx.db.attachedTxMu.Lock()
		defer tx.db.attachedTxMu.Unlock()

		if tx.db.attachedTransaction != nil {
			tx.db.attachedTransaction = nil
		}
	}

	for i := len(tx.onCommitHooks) - 1; i >= 0; i-- {
		tx.onCommitHooks[i]()
	}

	return nil
}


func (tx *Transaction) Unlock() {
	if tx.writable {
		tx.db.txmu.Unlock()
	} else {
		//Runlock方法解除rw的读取锁状态，如果m未加读取锁会导致运行时错误
		/*Rlock是加只读状态锁，DUnlock是解除只读状态*/
		tx.db.txmu.RUnlock()
	}
}

/*Writable返回该事务是否可写*/
func (tx *Transaction) Writable() bool {
	return tx.writable
}


/*事务的各项功能 TODO*/
func (tx *Transaction) CreateTable()  {

}


func (tx *Transaction) GetTable()  {

}

func (tx *Transaction) AddFieldConstraint()  {

}

func (tx *Transaction) RenameTable()  {

}

func (tx *Transaction) DropTable()  {

}

/*事务的索引功能 TODO*/
func (tx *Transaction) CreateIndex()  {

}

func (tx *Transaction) GetIndex()  {

}

func (tx *Transaction) DropIndex()  {

}

func (tx *Transaction) ListIndexes()  {

}

func (tx *Transaction) ReIndex()  {

}

func (tx *Transaction) ReIndexAll()  {

}

/*关于存储功能的实现 TODO*/
func (tx *Transaction) getTableStore()  {

}
func (tx *Transaction) getIndexStore()  {

}






