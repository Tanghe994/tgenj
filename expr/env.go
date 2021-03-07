package expr

import (
	"tgenj/database"
	"tgenj/document"
)

// Environment contains information about the context in which
// the expression is evaluated.
type Environment struct {
	//Params []Param
	//Vars   *document.FieldBuffer
	Doc    document.Document
	Tx     *database.Transaction

	Outer *Environment
}
