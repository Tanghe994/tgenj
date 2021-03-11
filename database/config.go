package database

import (
	"bytes"
	"encoding/binary"
	"tgenj/document"
	"tgenj/engine"
)

/**
 *  @ClassName:config
 *  @Description:TODO
 *  @Author:jackey
 *  @Create:2021/3/11 下午5:35
 */
const storePrefix = 't'

// FieldConstraint describes constraints on a particular field.
// 字段约束
type FieldConstraint struct {
	Path         document.Path
	Type         document.ValueType
	IsPrimaryKey bool
	IsNotNull    bool
	DefaultValue document.Value
	IsInferred   bool
	InferredBy   []document.Path
}

// FieldConstraints is a list of field constraints.
type FieldConstraints []*FieldConstraint

// TableInfo contains information about a table.
// TODO
type TableInfo struct {
	// name of the table.
	tableName string
	// name of the store associated with the table.
	storeName []byte
	readOnly  bool

	FieldConstraints FieldConstraints
}

func (ti *TableInfo) Clone() *TableInfo {
	cp := *ti
	cp.FieldConstraints = nil

	/*clone 字段制约*/
	for _, fc := range ti.FieldConstraints {
		cp.FieldConstraints = append(cp.FieldConstraints, fc)
	}
	return &cp
}

/*TODO ToDocument turns ti a document*/
/*将表的一个信息转换为一个文件*/
func (ti *TableInfo) ToDocument() document.Document {

}

// IndexInfo holds the configuration of an index.
type IndexInfo struct {
	TableName string
	IndexName string
	Path      document.Path

	// If set to true, values will be associated with at most one key. False by default.
	Unique bool

	// If set, the index is typed and only accepts that type
	Type document.ValueType
}

func (i IndexInfo) Clone() *IndexInfo {
	return &i
}

// tableStore manages table information.
// It loads table information during database startup
// and holds it in memory.
type tableStore struct {
	db *Database
	st engine.Store
}

/* Insert a new tableInfo for the given table name.*/
func (t *tableStore) Insert(tx *Transaction,tableName string,info *TableInfo) error{
	tblName := []byte(tableName)
	
	_,err := t.st.Get(tblName)

	/*你要修改名字。如果在里面找到已经有的名字，那就说明已经有同名的表存在engine中了*/
	if err == nil {	// nil的话代表已经找到
		return ErrTableAlreadyExists
	}

	/*如果没有找到，会饭返回 ErrTableNotFound错误*/
	if err != ErrTableNotFound {
		return nil
	}


	// TODO 这里没太明白
	if info.storeName == nil {
		seq,err := t.st.NextSequence()
		if err != nil {
			return err
		}

		/*c创建一个huff*/
		buf := make([]byte,binary.MaxVarintLen64+1)
		buf[0]	= storePrefix

		// 向buf中写数据
		n := binary.PutUvarint(buf[1:],seq)
		info.storeName=buf[:n+1]
		/*此时buf中是存有数据的*/
	}

	var buf bytes.Buffer	/*这是一个结构体*/
	encoder := t.db.Codec.NewEncoder(&buf)
		defer encoder.Close()

	/*TODO*/
	err = encoder.EncodeDocument(info.ToDocument())
	if err != nil {
		return err
	}

	err = t.st.Put([]byte(tableName),buf.Bytes())
	if err != nil {
		return err
	}

	return nil





	if  {
		
	}
}
