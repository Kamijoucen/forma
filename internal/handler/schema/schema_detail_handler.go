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

func SchemaDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SchemaDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := schema.NewSchemaDetailLogic(r.Context(), svcCtx)
		resp, err := l.SchemaDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
