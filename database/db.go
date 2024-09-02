package database

import (
	"context"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/database/pg"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Db struct {
	DB *gorm.DB
}

type DbFunc func(db *gorm.DB, span trace.Span) *gorm.DB

func NewDb(c mysqltool.Config) (*Db, error) {
	db, err := mysqltool.NewMySQLFromConfig(&c)
	if err != nil {
		return nil, err
	}
	return &Db{DB: db}, nil
}

func NewDbFromPostgres(c pg.Config) (*Db, error) {
	db, err := pg.NewPostgres(&c)
	if err != nil {
		return nil, err
	}
	return &Db{DB: db}, nil
}

func (db *Db) Do(ctx context.Context, do DbFunc) {
	db.DoWithName(ctx, "DB.Do", do)
}

func (db *Db) DoWithName(ctx context.Context, traceName string, do DbFunc) {
	tracing.TraceFunc(ctx, traceName, func(c1 context.Context, span trace.Span) {
		result := do(db.DB, span)
		if result.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", result.Error.Error()))

			span.RecordError(result.Error)
		}
	})
}

func (db *Db) Create(ctx context.Context, value any, do DbFunc, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "DB.Create", func(c1 context.Context, span trace.Span) {
		var d = db.DB
		if do != nil {
			d = do(d, span)
		}
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}

		res = d.Create(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))

			span.RecordError(res.Error)
		}
	})
	return res
}

func (db *Db) CreateBatch(ctx context.Context, value any, batchSize int, do DbFunc, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "DB.Create", func(c1 context.Context, span trace.Span) {
		var d = db.DB
		if do != nil {
			d = do(d, span)
		}
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}

		res = d.CreateInBatches(value, batchSize)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))

			span.RecordError(res.Error)
		}
	})
	return res
}

func (db *Db) Update(ctx context.Context, column string, value any, do DbFunc) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "DB.Update", func(c1 context.Context, span trace.Span) {
		var d = db.DB
		if do != nil {
			d = do(d, span)
		}
		res = d.Update(column, value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))
		}
	})
	return res
}

func (db *Db) Updates(ctx context.Context, value any, do DbFunc, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "DB.Updates", func(c1 context.Context, span trace.Span) {
		var d = db.DB
		if do != nil {
			d = do(d, span)
		}
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}
		res = d.Updates(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))

			span.RecordError(res.Error)
		}
	})
	return res
}

// Save wrap gorm.Save
//
// @value any struct pointer to save data
//
// @do conds gorm.clause.Expression
//
// for example:
//
//	ctx := context.Background()
//
//	var req = struct {
//		UserId int64  `json:"userId"`
//		Name   string `json:"name"`
//		Phone  string `json:"phone"`
//		State  int    `json:"state"`
//	}{
//		UserId: 0,
//		Name:   "Bob",
//		Phone:  "13311112222",
//		State:  1,
//	}
//
//	var newUser = User{
//		UserId:   req.UserId,
//		Username: "",
//		Phone:    req.Phone,
//		Account:  req.Phone,
//	}
//
//	err = DB.Save(ctx,&newUser,clause.OnConflict{
//		Columns:      []clause.Column{{Name: "phone"}},
//		DoUpdates:    clause.Assignments(map[string]interface{}{"state": 1}),
//	}).Error
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log(newUser)
func (db *Db) Save(ctx context.Context, value any, conds ...clause.Expression) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "DB.Save", func(c1 context.Context, span trace.Span) {
		var d = db.DB
		if len(conds) > 0 {
			d = d.Clauses(conds...)
		}
		res = d.Save(value)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))

			span.RecordError(res.Error)
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
// var users []User
//
//	var count int64
//	err = DB.Find(context.Background(), &users, func(tx *gorm.DB, span trace.Span) *gorm.DB {
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
	tracing.TraceFunc(ctx, "DB.Find", func(c1 context.Context, span trace.Span) {
		d := db.DB
		if do != nil {
			d = do(d, span)
		}
		res = d.Find(dest)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))

			span.RecordError(res.Error)
		}
	})
	return res
}

// First wrap gorm.First
//
// @dest your result struct
//
// @do do func You can use do func to do sth.
//
// for example:
//
//	var user User
//	err = DB.First(context.Background(), &user, func(tx *gorm.DB, span trace.Span) *gorm.DB {
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
//	  return tx
//	}).Error
func (db *Db) First(ctx context.Context, dest any, do DbFunc) *gorm.DB {
	var res *gorm.DB
	tracing.TraceFunc(ctx, "DB.First", func(c1 context.Context, span trace.Span) {
		d := db.DB
		if do != nil {
			d = do(d, span)
		}
		res = d.First(dest)
		if res.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", res.Error.Error()))

			span.RecordError(res.Error)
		}
	})
	return res
}

// Transaction wrap gorm Transaction
func (db *Db) Transaction(ctx context.Context, do func(tx *gorm.DB, span trace.Span) error) error {
	var err error
	tracing.TraceFunc(ctx, "DB.Transaction", func(c1 context.Context, span trace.Span) {
		err = db.DB.Transaction(func(tx *gorm.DB) error {
			return do(tx, span)
		})
		if err != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", err.Error()))
			span.RecordError(err)
		}
	})
	return err
}

// FindInBatch wrap gorm.FindInBatch
//
// @dest your result struct
//
// @do do func You can use do func to do sth.
//
// @batchSize batchSize
func (db *Db) FindInBatch(ctx context.Context, dest any, batchSize int, do func(tx *gorm.DB, span trace.Span) error) *gorm.DB {
	var result *gorm.DB
	tracing.TraceFunc(ctx, "DB.FindInBatch", func(c1 context.Context, span trace.Span) {
		result = db.DB.FindInBatches(&dest, batchSize, func(tx *gorm.DB, batch int) error {
			return do(tx, span)
		})
		if result.Error != nil {
			span.SetAttributes(attribute.Bool("DB.error", true), attribute.String("DB.errorString", result.Error.Error()))
			span.RecordError(result.Error)
		}
	})
	return result
}
