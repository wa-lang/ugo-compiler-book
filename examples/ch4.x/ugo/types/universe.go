package types

var Universe *Scope

func init() {
	Universe = NewScope(nil)
}
