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

type DoFunc func(db *gorm.DB)

func NewDb(c mysqltool.Config) (*Db, error) {
	db, err := mysqltool.NewMySQLFromConfig(&c)
	if err != nil {
		return nil, err
	}
	return &Db{Mysql: db}, nil
}

func (db *Db) Do(ctx context.Context, do DoFunc) {
	tracing.TraceFunc(ctx, "db.Do", func(span trace.Span) {
		do(db.Mysql)
	})
}

func (db *Db) DoWithName(ctx context.Context, traceName string, do DoFunc) {
	tracing.TraceFunc(ctx, traceName, func(span trace.Span) {
		do(db.Mysql)
	})
}

func (db *Db) Create(ctx context.Context, value any, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Create", func(span trace.Span) {
		var d = db.Mysql
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}

		res = d.Create(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true))
		}
	})
	return res
}

func (db *Db) Update(ctx context.Context, column string, value any) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Update", func(span trace.Span) {
		res = db.Mysql.Update(column, value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true))
		}
	})
	return res
}

func (db *Db) Updates(ctx context.Context, value any, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Updates", func(span trace.Span) {
		var d = db.Mysql
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}
		res = d.Updates(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true))
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
			span.SetAttributes(attribute.Bool("db.error", true))
		}
	})
	return res
}

func (db *Db) Find(ctx context.Context, dest any) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.Find", func(span trace.Span) {
		res = db.Mysql.Find(dest)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true))
		}
	})
	return res
}

func (db *Db) First(ctx context.Context, dest any) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "gorm.First", func(span trace.Span) {
		res = db.Mysql.First(dest)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("db.error", true))
		}
	})
	return res
}

func (db *Db) Transaction(ctx context.Context, do func(tx *gorm.DB) error) error {
	var err error
	tracing.TraceFunc(ctx, "gorm.Transaction", func(span trace.Span) {
		err = db.Mysql.Transaction(func(tx *gorm.DB) error {
			return do(tx)
		})
		if err != nil {
			span.SetAttributes(attribute.Bool("db.error", true))
		}
	})
	return err
}
