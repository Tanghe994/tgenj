package document


/*TODO*/
type Document interface {
	Iterate(fn func(field string,value Value)error)error

	GetByField(field string)(Value,error)
}


// A Path represents the path to a particular value within a document.
type Path []PathFragment


// PathFragment is a fragment of a path representing either a field name or
// the index of an array.
type PathFragment struct {
	FieldName  string
	ArrayIndex int
}