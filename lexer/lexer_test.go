package lexer

import (
	"myRVCC/logger"
	"myRVCC/token"
	"testing"
	"text/scanner"
)

func TestScanNext(t *testing.T) {
	l := New("1 + 2 - 3")
	for {
		logger.Info("%+v", l.Pos())
		next := l.Next()
		if next == scanner.EOF {
			break
		}
		logger.Info("%c", next)
	}
}
func TestNextToken(t *testing.T) {
	l := New("1 + 2 - 33+55")
	for {
		t := l.NextToken()
		logger.Info("%s %s", t.Kind, t.Literal)
		if t.Kind == token.EOF {
			break
		}
	}
}
func TestError(t *testing.T) {
	l := New("1>2")
	for {
		t := l.NextToken()
		logger.Info("%s %s", t.Kind, t.Literal)
		if t.Kind == token.EOF {
			break
		}
	}
}
