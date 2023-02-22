package mysqltool

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type IdKey struct {
	Id        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time            `json:"deleted_at"`
	IsDeleted soft_delete.DeletedAt `gorm:"type:tinyint(4);softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"is_deleted"`
}
