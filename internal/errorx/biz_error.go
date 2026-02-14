package errorx

import "fmt"

// BizError 业务错误，携带错误码和用户友好的错误消息
type BizError struct {
	Code    string
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

// NewBizError 创建业务错误
func NewBizError(code, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// NewBizErrorf 创建业务错误，支持格式化消息
func NewBizErrorf(code, format string, args ...any) *BizError {
	return &BizError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
