package svc

import (
	"context"
	"errors"

	"forma/internal/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ResponseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func initHandler() {
	httpx.SetOkHandler(func(ctx context.Context, a any) any {
		return &ResponseBody{
			Code:    "200",
			Message: "success",
			Data:    a,
		}
	})

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		var bizErr *errorx.BizError
		if errors.As(err, &bizErr) {
			return 200, &ResponseBody{
				Code:    bizErr.Code,
				Message: bizErr.Message,
			}
		}

		// 非业务错误，记录日志，不暴露内部细节
		logx.WithContext(ctx).Errorf("internal error: %v", err)
		return 200, &ResponseBody{
			Code:    errorx.CodeInternal,
			Message: "系统内部错误",
		}
	})
}
