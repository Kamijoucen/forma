package service

import (
	"context"
	"time"

	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/errorx"
	"forma/internal/types"
)

// QueryAppByCode 按 code 查询 App，不存在时返回 ErrAppNotFound
func QueryAppByCode(ctx context.Context, client *ent.Client, code string) (*ent.App, error) {
	a, err := client.App.Query().
		Where(entApp.CodeEQ(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorx.ErrAppNotFound
		}
		return nil, err
	}
	return a, nil
}

// ToAppDetailResp 将 Ent App 实体转为 API 响应
func ToAppDetailResp(a *ent.App) *types.AppDetailResp {
	return &types.AppDetailResp{
		Code:        a.Code,
		Name:        a.Name,
		Description: a.Description,
		CreatedAt:   a.CreateTime.Format(time.DateTime),
		UpdatedAt:   a.UpdateTime.Format(time.DateTime),
	}
}
