package parser

import "tgenj/expr"

/*Options of the Sql parser*/
type Options struct {
	/*A map of builtin SQL functions*/
	Functions expr.Functions
}

func defaultOptions() *Options {
	return &Options{
		Functions: expr.NewFunctions(),
	}
}