package database

import (
	"context"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Db struct {
	Mysql *gorm.DB
}

var Context Db

type DbFunc func(db *gorm.DB, span trace.Span) *gorm.DB

func NewDb(c mysqltool.Config) (*Db, error) {
	db, err := mysqltool.NewMySQLFromConfig(&c)
	if err != nil {
		return nil, err
	}
	return &Db{Mysql: db}, nil
}

func (db *Db) Do(ctx context.Context, do DbFunc) {
	db.DoWithName(ctx, "db.Do", do)
}

func (db *Db) DoWithName(ctx context.Context, traceName string, do DbFunc) {
	tracing.TraceFunc(ctx, traceName, func(span trace.Span) {
		result := do(db.Mysql, span)
		if result.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", result.Error.Error()))
		}
	})
}

func (db *Db) Create(ctx context.Context, value any, do DbFunc, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Create", func(span trace.Span) {
		var d = db.Mysql
		if do != nil {
			d = do(d, span)
		}
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}

		res = d.Create(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

func (db *Db) CreateBatch(ctx context.Context, value any, batchSize int, do DbFunc, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Create", func(span trace.Span) {
		var d = db.Mysql
		if do != nil {
			d = do(d, span)
		}
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}

		res = d.CreateInBatches(value, batchSize)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

func (db *Db) Update(ctx context.Context, column string, value any, do DbFunc) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Update", func(span trace.Span) {
		var d = db.Mysql
		if do != nil {
			d = do(d, span)
		}
		res = d.Update(column, value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

func (db *Db) Updates(ctx context.Context, value any, do DbFunc, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Updates", func(span trace.Span) {
		var d = db.Mysql
		if do != nil {
			d = do(d, span)
		}
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}
		res = d.Updates(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

func (db *Db) Save(ctx context.Context, value any, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Save", func(span trace.Span) {
		var d = db.Mysql
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}
		res = d.Save(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

// Find wrap gorm.Find and opentelemetry trace,finds records using giving conditions by do func
//
// @dest your result struct
//
// @do do func You can use do func to do sth.
//
// for example:
//
// var users []UserPageItem
//
//	var count int64
//	err = db.Find(context.Background(), &users, func(tx *gorm.DB, span trace.Span) *gorm.DB {
//	  tx = tx.Table("t_app_user a").Joins("left join t_user_vip b on a.user_id=b.user_id")
//	  if req.UserId > 0 {
//		tx = tx.Where("a.user_id=?", req.UserId)
//	  }
//	  if req.Name != "" {
//		tx = tx.Where("a.username like ?", fmt.Sprintf("%%%s%%", req.Name))
//	  }
//	  if req.IsVip > 0 {
//		tx = tx.Where("a.is_vip=?", req.IsVip)
//	  }
//	  return tx.Count(&count).Order("user_id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
//	  }).Error
//	  if err != nil {
//		t.Error(err)
//		return
//	  }
//	  t.Log(count, users)
func (db *Db) Find(ctx context.Context, dest any, do DbFunc) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Find", func(span trace.Span) {
		d := db.Mysql
		if do != nil {
			d = do(d, span)
		}
		res = d.Find(dest)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

func (db *Db) First(ctx context.Context, dest any, do DbFunc) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.First", func(span trace.Span) {
		d := db.Mysql
		if do != nil {
			d = do(d, span)
		}
		res = d.First(dest)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", res.Error.Error()))
		}
	})
	return res
}

// Transaction wrap gorm Transaction
func (db *Db) Transaction(ctx context.Context, do func(tx *gorm.DB, span trace.Span) error) error {
	var err error
	tracing.TraceFunc(ctx, "gorm.Transaction", func(span trace.Span) {
		err = db.Mysql.Transaction(func(tx *gorm.DB) error {
			return do(tx, span)
		})
		if err != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", err.Error()))
		}
	})
	return err
}

func (db *Db) FindInBatch(ctx context.Context, dest any, batchSize int, do func(tx *gorm.DB, span trace.Span) error) *gorm.DB {
	var result *gorm.DB
	tracing.TraceFunc(ctx, "gorm.FindInBatch", func(span trace.Span) {
		result = db.Mysql.FindInBatches(&dest, batchSize, func(tx *gorm.DB, batch int) error {
			return do(tx, span)
		})
		if result.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true), attribute.String("db.errorString", result.Error.Error()))
		}
	})
	return result
}
