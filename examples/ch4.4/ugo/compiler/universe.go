package compiler

var Universe *Scope = NewScope(nil)

var builtinObjects = []*Object{
	{Name: "println", MangledName: "@ugo_builtin_println"},
	{Name: "exit", MangledName: "@ugo_builtin_exit"},
}

func init() {
	for _, obj := range builtinObjects {
		Universe.Insert(obj)
	}
}
