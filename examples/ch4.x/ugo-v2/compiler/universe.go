package compiler

var Universe *Scope = NewScope(nil)

var builtinObjects = []*Object{
	{Name: "println", LLName: "@ugo_builtin_println"},
	{Name: "exit", LLName: "@ugo_builtin_exit"},
}

func init() {
	for _, obj := range builtinObjects {
		Universe.Insert(obj)
	}
}
