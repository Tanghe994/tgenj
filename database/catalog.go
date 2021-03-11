package database

import (
	"errors"
	"math/rand"
	"sync"
)

// Catalog holds all table and index informations.
type Catalog struct {
	cache *catalogCache
}

// catalogCache 需要互斥访问
type catalogCache struct {
	tables           map[string]*TableInfo
	indexes          map[string]*IndexInfo
	indexesPerTables map[string][]*IndexInfo

	mu sync.RWMutex
}

// RenameTable renames a table.
// If it doesn't exist, it returns ErrTableNotFound.
// TODO
func (c *Catalog) RenameTable(tx *Transaction, oldName, newName string) error {
	newTi, newIdxs, err := c.cache.updateTable(tx, oldName, func(clone *TableInfo) error {
		clone.tableName = newName
		return nil
	})

	if err != nil {
		return err
	}

	tableStore := tx.getTableStore()

	err = tableStore.Insert(tx,newName,newTi)
	if err != nil {
		return err
	}




	return err
}

/*updateTable 更新表info和索引info，但是并没有更新底层的TableStore*/
func (c *catalogCache) updateTable(tx *Transaction, tableName string, fn func(clone *TableInfo) error) (*TableInfo, []*IndexInfo, error) {
	/*加锁*/
	c.mu.Lock()
	defer c.mu.Unlock()

	/*判断要修改的表是否在缓存中,ti=tableInfo*/
	ti, ok := c.tables[tableName]
	if ok {
		return nil, nil, ErrTableNotFound
	}

	if ti.readOnly {
		return nil, nil, errors.New("connot write to read-only table")
	}

	clone := ti.Clone()

	/*更改名字*/
	err := fn(clone)
	/*旧表的名字已经更改成新表名字*/

	if err != nil {
		return nil, nil, err
	}

	/*创建索引*/
	var oldIndexes, newIndexes []*IndexInfo

	if clone.tableName != tableName { // 不等于才是更改成功了
		/*删除目录缓存中的索引*/
		delete(c.tables, tableName)

		/*还要删除旧表的索引，增加更名后表的索引*/
		for _, idx := range c.indexes { // 这里可能不止一个索引
			if idx.TableName == tableName {
				idxClone := idx.Clone() /*创建新的索引*/
				idxClone.TableName = clone.tableName

				newIndexes = append(newIndexes, idxClone)
				/*为了事务回滚*/
				oldIndexes = append(oldIndexes, idx)

				c.indexes[idxClone.IndexName] = idxClone
			}
		}
	}

		c.tables[clone.tableName] = clone

		tx.onRollbackHooks = append(tx.onRollbackHooks, func() {
			c.mu.Lock()
			defer c.mu.Unlock()

			delete(c.tables,clone.tableName)
			c.tables[tableName]=ti

			for _,idx := range oldIndexes{
				c.indexes[idx.IndexName]= idx
			}

		})
	return clone,newIndexes,nil
}
