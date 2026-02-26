// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package app

import (
	"net/http"

	"forma/internal/logic/app"
	"forma/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AppListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := app.NewAppListLogic(r.Context(), svcCtx)
		resp, err := l.AppList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
