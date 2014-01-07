//scan func from the source

package lib

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

type Func struct {
	Name     string
	Receiver Value
	Values   Values
	Retruns  Values
}

type Funcs []Func

func (f *Funcs) Scan(tok token.Token, lit string, s *scanner.Scanner) {
	funcation := Func{}

	if funcation.scanFunc(tok, lit, s) {
		f = append(f, funcation)
	}
}

func (f *Func) String() string {
	str := fmt.Sprintf("\nrececiver:%v\nname:%s\nvalues:%v\nreturn:%v", f.Receiver, f.Name, f.Values, f.Retruns)
	return str
}

func scan(s *scanner.Scanner) (token.Token, string) {
	_, tok, lit := s.Scan()
	return tok, lit
}

func (f *Func) scanFunc(tok token.Token, lit string, s *scanner.Scanner) bool {
	if tok != token.FUNC {
		return false
	}

	tok, lit = scan(s)
	if tok != token.LPAREN {
		return false
	}
	f.Receiver = scanReceiver(tok, lit, s)

	tok, lit = scan(s)
	if tok != token.IDENT {
		return false
	}
	f.Name = lit

	tok, lit = scan(s)
	if tok != token.LPAREN {
		return false
	}
	f.Values = scanValues(tok, lit, s)

	tok, lit = scan(s)
	f.Retruns = scanReturn(tok, lit, s)

	return true
}

const (
	VALUE_START = iota
	VALUE_NAME
	VALUE_TYPE
	VALUE_END
)

func scanReceiver(tok token.Token, lit string, s *scanner.Scanner) (value Value) {
	if tok != token.LPAREN {
		return
	}

	tok, lit = scan(s)
	if tok != token.IDENT {
		return
	}
	value.Name = lit

	tok, lit = scan(s)
	if tok != token.MUL {
		value.Type = lit
		return
	}
	value.Type = tok.String()

	tok, lit = scan(s)
	value.Type += lit
	scan(s)
	return
}

func scanValues(tok token.Token, lit string, s *scanner.Scanner) (values []Value) {
	// fmt.Println("--->scan values")
	state := VALUE_START
	var value Value
	for {
		// fmt.Printf("\n%v \t%s \t%q\n", state, tok, lit)

		switch state {
		case VALUE_START:
			if tok == token.LPAREN {
				state = VALUE_NAME
			} else {
				state = VALUE_END
			}
		case VALUE_NAME:
			if tok == token.RPAREN {
				state = VALUE_END
				break
			}
			value.Name = lit
			// fmt.Println(">>>get name:", value.Name)
			state = VALUE_TYPE
		case VALUE_TYPE:
			if tok == token.RPAREN {
				state = VALUE_END
				values = append(values, value)
				break
			}
			if tok == token.COMMA {
				state = VALUE_NAME
				values = append(values, value)
				value.Name = ""
				value.Type = ""
				break
			}
			if tok == token.MUL {
				value.Type = tok.String()
			} else {
				value.Type += lit
				// fmt.Println(">>>get type:", value.Type)

				for index, v := range values {
					if v.Type == "" {
						v.Type = value.Type
						values[index] = v
					}
				}
			}
		}

		if state == VALUE_END {
			break
		}
		if state != VALUE_START {
			_, tok, lit = s.Scan()
		}
		if tok == token.EOF {
			break
		}
	}
	// fmt.Println("-------------", tok)
	return
}

func scanReturn(tok token.Token, lit string, s *scanner.Scanner) (values []Value) {
	// fmt.Println("--->scan return")
	if tok == token.LPAREN {
		values = scanValues(tok, lit, s)
	}

	if tok == token.IDENT {
		value := Value{Type: lit}
		values = append(values, value)
		scan(s)
	}

	return
}