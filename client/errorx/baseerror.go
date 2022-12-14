package errorx

const defaultCode = 10001
type CodeError struct {
	Code int   	`json:"code"`
	Msg  string 	`json:"msg"`
}

type CodeErrorResponse struct {
	Code int 	`json:"code"`
	Msg string 	`json:"msg"`
}

func NewCodeDefaultError(msg string) error {
	return NewCodeError(defaultCode,msg)
}
func NewCodeError(code int,msg string) error {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}
func(c *CodeError) Error() string {
	return c.Msg
}

func(c *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: c.Code,
		Msg:  c.Msg,
	}
}
