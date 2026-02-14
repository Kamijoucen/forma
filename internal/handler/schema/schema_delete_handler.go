// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"net/http"

	"forma/internal/logic/schema"
	"forma/internal/svc"
	"forma/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SchemaDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SchemaDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := schema.NewSchemaDeleteLogic(r.Context(), svcCtx)
		err := l.SchemaDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
