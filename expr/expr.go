package expr

import "tgenj/document"

// An Expr evaluates to a value.
type Expr interface {
	Eval(*Environment) (document.Value, error)
}