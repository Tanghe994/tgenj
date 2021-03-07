package parser

import (
	"bytes"
	"io"
	"strings"
	"tgenj/expr"
	"tgenj/query"
	"tgenj/sql/scanner"
)

/*Parser represents an Tgenj SQL Parser */
/*TODO*/
type Parser struct {
	s             *scanner.BufScanner
	orderedParams int
	namedParams   int
	buf           *bytes.Buffer
	functions     expr.Functions
}

func NewParser(r io.Reader) *Parser {
	return NewParserWithOptions(r, nil)
}

/*NewParserWithOptions returns a new instance of Parser using given Options*/
func NewParserWithOptions(r io.Reader, opts *Options) *Parser {
	if opts == nil {
		opts = defaultOptions()
	}
	return &Parser{
		s:         scanner.NewBufScanner(r),
		functions: opts.Functions,
	}
}

/*ParseQuery parses a query string and returns its AST representation.*/
/*TODO*/
func ParseQuery(s string) (query.Query, error) {
	/**
	根据查询字符串创建一个新的阅读器
	通过NewParser创建一个Parser
	然后执行ParserQuery创建一个Query
	*/
	return NewParser(strings.NewReader(s)).ParserQuery()
}

/*ParserQuery parser a Tgenj Sql string and returns a Query*/
func (p *Parser) ParserQuery() (query.Query, error) {
	/*操作 接口切片*/
	var statements []query.Statement
	semi := true

	for {
		if {

		}
	}

}


/*ScanIgnoreWhitespace scans the next non-whitespace and non-comment token*/
/*TODO*/
func (p *Parser) ScanIgnoreWhitespace() (tok scanner.Token, pos scanner.Pos, lit string) {
	for true {
		tok,pos,lit=p.
	}
}

/*TODO*/
func (p *Parser) Scan()(tok scanner.Token,pos scanner.Pos,lit string)  {
	ti :=p.s.scan
}