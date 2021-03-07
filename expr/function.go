package expr


// Functions represents a map of builtin SQL functions.
type Functions struct {
	m map[string]func(args ...Expr) (Expr, error)
}

func NewFunctions() Functions {
	return Functions{
		m: BuiltinFunctions(),
	}
}



/*BuiltinFunctions returns default map of builtin functions.*/
func BuiltinFunctions() map[string] func(args ... Expr)(Expr,error){
	return nil,nil
}