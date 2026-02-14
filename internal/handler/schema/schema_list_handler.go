// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schema

import (
	"net/http"

	"forma/internal/logic/schema"
	"forma/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SchemaListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := schema.NewSchemaListLogic(r.Context(), svcCtx)
		resp, err := l.SchemaList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
