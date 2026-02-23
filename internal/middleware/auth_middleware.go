// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"forma/internal/config"
	"forma/internal/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	c config.Config
}

func NewAuthMiddleware(c config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		c: c,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			httpx.ErrorCtx(r.Context(), w, errorx.ErrUnauthorized)
			return
		}

		// TODO 先临时实现验证，后续替换为真正的认证逻辑
		if token != m.c.TempToken {
			httpx.ErrorCtx(r.Context(), w, errorx.ErrUnauthorized)
			return
		}

		next(w, r)
	}
}
