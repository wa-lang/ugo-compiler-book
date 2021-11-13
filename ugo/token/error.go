package token

type Error struct {
	Filename string
	Source   []byte
	Pos      Pos
	Msg      string
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Filename != "" || e.Pos.IsValid() {
		return PosString(e.Filename, e.Source, e.Pos) + ": " + e.Msg
	}
	return e.Msg
}
