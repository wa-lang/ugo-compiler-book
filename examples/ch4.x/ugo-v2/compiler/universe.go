package compiler

var Universe *Scope

func init() {
	Universe = NewScope(nil)
}
