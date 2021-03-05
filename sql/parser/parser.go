package parser

import (
	"bytes"
	"tgenj/query"
	"tgenj/sql/scanner"
)

type Parser struct {
	s             *scanner.BufScanner
	orderedParams int
	namedParams   int
	buf           *bytes.Buffer
	//functions     expr.Functions
}

/*ParseQuery parses a query string and returns its AST representation.*/
/*TODO*/
func ParseQuery(s string)(query.Query,error)  {
	return nil,nil
}