package errorx

// 业务错误码常量
const (
	CodeSuccess      = "0"     // 成功
	CodeInvalidParam = "10001" // 参数校验错误
	CodeNotFound     = "10002" // 资源不存在
	CodeUnauthorized = "10003" // 未授权
	CodeForbidden    = "10004" // 无权限
	CodeInternal     = "99999" // 系统内部错误
)

// 预定义常用业务错误
var (
	ErrInvalidParam = NewBizError(CodeInvalidParam, "参数错误")
	ErrNotFound     = NewBizError(CodeNotFound, "资源不存在")
	ErrUnauthorized = NewBizError(CodeUnauthorized, "未授权")
	ErrForbidden    = NewBizError(CodeForbidden, "无权限")
	ErrInternal     = NewBizError(CodeInternal, "系统内部错误")
)
