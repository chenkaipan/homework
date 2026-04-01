package utils

type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	return e.Msg
}

// 工厂函数
func NewBizError(code int, msg string) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
	}
}
