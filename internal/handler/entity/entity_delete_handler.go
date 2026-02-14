// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"net/http"

	"forma/internal/logic/entity"
	"forma/internal/svc"
	"forma/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EntityDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EntityDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := entity.NewEntityDeleteLogic(r.Context(), svcCtx)
		err := l.EntityDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
