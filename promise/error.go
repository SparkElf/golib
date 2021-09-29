package promise

//1-500 is used for internal errors.
const (
	SUCCESS = iota
	TIMEOUT
	INTERRUPT
)

var (
	codemap = map[int]string{
		//0-500
		SUCCESS:   "success",
		TIMEOUT:   "timeout",
		INTERRUPT: "interrupt",
	}
)

type Err struct {
	code int
	Msg  string
}

func (e *Err) Error() string {
	return codemap[e.code]
}
