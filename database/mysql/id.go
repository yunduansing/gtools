package mysqltool

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type IdType struct {
	Id        int64      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsDelType
}

type IsDelType struct {
	DeletedAt *time.Time            `json:"deleted_at"`
	IsDeleted soft_delete.DeletedAt `gorm:"type:tinyint(4);softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"is_deleted"`
}
