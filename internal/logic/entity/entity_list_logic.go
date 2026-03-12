// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package entity

import (
	"context"

	"forma/internal/ent"
	entApp "forma/internal/ent/app"
	"forma/internal/ent/entityfieldvalue"
	"forma/internal/ent/entityrecord"
	"forma/internal/ent/fielddef"
	"forma/internal/ent/schemadef"
	"forma/internal/service"
	"forma/internal/svc"
	"forma/internal/types"
	"forma/internal/util"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type EntityListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEntityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EntityListLogic {
	return &EntityListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EntityListLogic) EntityList(req *types.EntityListReq) (resp *types.EntityListResp, err error) {

	baseFilter := entityrecord.HasSchemaDefWith(
		schemadef.NameEQ(req.SchemaName),
		schemadef.HasAppWith(entApp.CodeEQ(req.AppCode)),
	)

	err = util.WithTkClient(l.ctx, l.svcCtx.Ent, func(client *ent.Client) error {
		// 查询总数
		total, countErr := client.EntityRecord.Query().
			Where(baseFilter).
			Count(l.ctx)
		if countErr != nil {
			return countErr
		}

		query := client.EntityRecord.Query().
			WithFieldValues(func(q *ent.EntityFieldValueQuery) {
				q.WithFieldDef()
			}).
			Where(baseFilter)

		if req.SortParam != nil && req.SortParam.Field != "" {
			sortField := req.SortParam.Field
			query.Order(func(s *sql.Selector) {
				// 关联子查询：获取指定排序字段的值
				efv := sql.Table(entityfieldvalue.Table).As("sort_efv")
				fd := sql.Table(fielddef.Table).As("sort_fd")
				sub := sql.Select(efv.C(entityfieldvalue.FieldValue)).
					From(efv).
					Join(fd).On(efv.C(entityfieldvalue.FieldDefColumn), fd.C(fielddef.FieldID)).
					Where(sql.And(
						sql.ColumnsEQ(efv.C(entityfieldvalue.EntityRecordColumn), s.C(entityrecord.FieldID)),
						sql.EQ(fd.C(fielddef.FieldName), sortField),
					)).
					Limit(1)
				dir := " ASC"
				if req.SortParam.Direction == "desc" {
					dir = " DESC"
				}
				s.OrderExpr(sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString("(")
					b.Join(sub)
					b.WriteString(")")
					b.WriteString(dir)
				}))
			})
		}

		// 分页查询
		records, queryErr := query.
			Offset((req.Page - 1) * req.PageSize).
			Limit(req.PageSize).
			All(l.ctx)
		if queryErr != nil {
			return queryErr
		}

		list := lo.Map(records, func(r *ent.EntityRecord, _ int) *types.EntityDetailResp {
			return service.ToEntityDetailResp(r, req.SchemaName)
		})

		resp = &types.EntityListResp{
			Total: int64(total),
			List:  list,
		}
		return nil
	})

	return resp, err
}
