package constant

type ErrCode struct {
	Code int32
	Msg  string
}

var (
	Success = ErrCode{Code: 00000, Msg: "successfully"}
)
