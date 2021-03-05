package document

/*TODO*/

type ValueType uint8

type Value struct {
	Type ValueType
	V    interface{}
}