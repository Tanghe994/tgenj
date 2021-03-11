package parser

import (
	"bytes"
	"fmt"
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
		if tok,pos,lit := p.ScanIgnoreWhitespace();tok == scanner.EOF{
			return query.New(statements...),nil
		}else if tok == scanner.SEMICOLON{	// 分号
			semi =true
		}else {
			if !semi {
				return query.Query{},newParseError(scanner.Tokstr(tok,lit),[]string{";"},pos)
			}
			p.Unscan()
			s,err := p.ParseStatement()

		}
	}

}

// ParseStatement parses a Genji SQL string and returns a Statement AST object.
func (p *Parser) ParseStatement() (query.Statement, error) {
	tok, pos, lit := p.ScanIgnoreWhitespace()
	switch tok {
	case scanner.ALTER:
		return p.parseAlterStatement()
	case scanner.BEGIN:
		return p.parseBeginStatement()
	case scanner.COMMIT:
		return p.parseCommitStatement()
	case scanner.SELECT:
		return p.parseSelectStatement()
	case scanner.DELETE:
		return p.parseDeleteStatement()
	case scanner.UPDATE:
		return p.parseUpdateStatement()
	case scanner.INSERT:
		return p.parseInsertStatement()
	case scanner.CREATE:
		return p.parseCreateStatement()
	case scanner.DROP:
		return p.parseDropStatement()
	case scanner.EXPLAIN:
		return p.parseExplainStatement()
	case scanner.REINDEX:
		return p.parseReIndexStatement()
	case scanner.ROLLBACK:
		return p.parseRollbackStatement()
	}

	return nil, newParseError(scanner.Tokstr(tok, lit), []string{
		"ALTER", "BEGIN", "COMMIT", "SELECT", "DELETE", "UPDATE", "INSERT", "CREATE", "DROP", "EXPLAIN", "REINDEX", "ROLLBACK",
	}, pos)
}


/*ScanIgnoreWhitespace scans the next non-whitespace and non-comment token*/
/*TODO*/
func (p *Parser) ScanIgnoreWhitespace() (tok scanner.Token, pos scanner.Pos, lit string) {
	for {
		/*得到token的信息，位置信息，还有lit?*/
		tok,pos,lit=p.Scan()
		if tok == scanner.WS || tok == scanner.COMMIT {
			continue
		}
	}
}

/*TODO*/
/*Scan returns the next token from the underlying scanner*/
func (p *Parser) Scan()(tok scanner.Token,pos scanner.Pos,lit string)  {
	ti :=p.s.Scan()

	if p.buf != nil {
		// 将Raw（源）的信息append到buff中
		p.buf.WriteString(ti.Raw)
	}
	tok, pos,lit =ti.Tok,ti.Pos,ti.Lit
	return tok,pos,lit
}

// Unscan pushes the previously read token back onto the buffer.
func (p *Parser) Unscan() {
	if p.buf != nil {
		ti := p.s.Curr()
		p.buf.Truncate(p.buf.Len() - len(ti.Raw))
	}
	p.s.Unscan()
}


// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      scanner.Pos
}

func newParseError(found string, expected []string, pos scanner.Pos) *ParseError {
	return &ParseError{
		Found: found,
		Expected: expected,
		Pos: pos,
	}
}

func (e *ParseError) Error()string {
	if e.Message != "" {
		/*格式化输出字符串*/
		return fmt.Sprintf("%s at line %d, char %d", e.Message, e.Pos.Line+1, e.Pos.Char+1)
	}
	return fmt.Sprintf("found %s, expected %s at line %d, char %d", e.Found, strings.Join(e.Expected, ", "), e.Pos.Line+1, e.Pos.Char+1)
}