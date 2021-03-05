package document


/*TODO*/
type Document interface {
	Iterate(fn func(field string,value Value)error)error

	GetByField(field string)(Value,error)
}
